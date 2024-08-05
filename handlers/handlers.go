package handlers

import (
	"database/sql"
	"encoding/json"

	//"log"
	"net/http"
	"time"

	"github.com/AramLab/todo-list/types"
	"github.com/AramLab/todo-list/utils"
)

type Handlers struct {
	DB *sql.DB
}

func NewHandlers(db *sql.DB) *Handlers {
	return &Handlers{DB: db}
}

func (h *Handlers) NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		now := r.URL.Query().Get("now")
		date := r.URL.Query().Get("date")
		repeat := r.URL.Query().Get("repeat")

		timeNow, err := time.Parse("20060102", now)
		if err != nil {
			return
		}

		nextDate, err := utils.NextDate(timeNow, date, repeat)
		if err != nil {
			return
		}

		response := nextDate
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}
}

func (h *Handlers) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task types.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.ValidateTitle(task.Title); err != nil {
		http.Error(w, `{"error":"title is required"}`, http.StatusBadRequest)
		return
	}

	// Проверка даты.
	// Устанавливаем текущую дату, если поле даты пусто.
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	// Проверка формата даты.
	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, `{"error":"дата представлена в формате, отличном от 20060102"}`, http.StatusBadRequest)
		return
	}

	// Проеврка актуальности даты.
	today := time.Now().Truncate(24 * time.Hour)
	if date.Before(today) {
		if task.Repeat == "" {
			task.Date = today.Format("20060102")
		} else {
			nextDate, err := utils.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"правило повторения указано в неправильном формате"}`, http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	}

	// Проверка шаблонов правил.
	if err = utils.ValidatePepeat(task.Repeat); err != nil {
		http.Error(w, `{"error":"invalid repeat pattern"}`, http.StatusBadRequest)
		return
	}

	res, err := h.DB.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"id": id,
	}

	json.NewEncoder(w).Encode(response)
}

/*func (h *Handlers) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var task types.Task
		var buf bytes.Buffer

		// Чтение тела запроса.
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, "ошибка чтения запроса", http.StatusBadRequest)
			return
		}

		// Десериализация запроса.
		if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
			http.Error(w, "ошибка десериализации JSON", http.StatusBadRequest)
			return
		}

		// Проверка наличия заголовка.
		if err := utils.ValidateTitle(task.Title); err != nil {
			http.Error(w, "не указан заголовок задачи", http.StatusBadRequest)
			return
		}

		// Проверка корректности формата даты.
		date, err := utils.ValidateDate(task.Date)
		if err != nil {
			http.Error(w, "дата представлена в формате, отличном от 20060102", http.StatusBadRequest)
			return
		} else if date.Before(time.Now()) {
			if task.Repeat != "" {
				newDate, err := utils.NextDate(time.Now(), date.Format("20060102"), task.Repeat)
				if err != nil {
					http.Error(w, "ошибка при вызове функции NextDate", http.StatusBadRequest)
					return
				}
				date, err = time.Parse("20060102", newDate)
				if err != nil {
					http.Error(w, "ошибка преобразования строки в время", http.StatusBadRequest)
					return
				}
			} else {
				date = time.Now()
			}
		}

		dateStr := date.Format("20060102")

		if err := utils.ValidatePepeat(task.Repeat); err != nil {
			http.Error(w, "правило повторения указано в неправильном формате", http.StatusBadRequest)
			return
		}

		// Делаем SQL-запрос на добавление задачи в таблицу.
		res, err := h.DB.Exec("insert into scheduler (date, title, comment, repeat) values (:date, :title, :comment, :repeat)",
			sql.Named("date", dateStr),
			sql.Named("title", task.Title),
			sql.Named("comment", task.Comment),
			sql.Named("repeat", task.Repeat))
		if err != nil {
			http.Error(w, "ошибка добалвения задачи в БД", http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "ошибка получения индекса заадчи", http.StatusInternalServerError)
			return
		}

		index := strconv.FormatInt(id, 10)

		response := types.Response{ID: index}

		responseJson, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "ошибка сериализации ответа", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJson)
	}
}*/

func (h *Handlers) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := utils.GetTasks(h.DB)
	if err != nil {
		//log.Printf("Error retrieving tasks: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//log.Printf("Retrieved %d tasks", len(tasks))

	// Если tasks == nil, инициализируйте его как пустой слайс
	if tasks == nil {
		tasks = []types.Task{}
	}

	response := types.TasksResponse{Tasks: tasks}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		//log.Printf("Error marshalling JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *Handlers) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}
	var task types.Task
	task, err := utils.GetTask(id, h.DB)
	if err != nil {
		http.Error(w, `{"error":"Ошибка чтения данных из базы данных"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handlers) PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask types.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, `{"error":"Ошибка десериализации"}`, http.StatusBadRequest)
		return
	}

	// Проверка наличия `id`.
	if newTask.ID == "" {
		http.Error(w, `{"error":"Не указан ID задачи"}`, http.StatusBadRequest)
		return
	}

	// Проверка существования задачи по данному ID.
	if _, err := utils.GetTask(newTask.ID, h.DB); err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusBadRequest)
		return
	}

	// Проверка валидности заголовка.
	if err := utils.ValidateTitle(newTask.Title); err != nil {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Проверка валидности правила повторения.
	if err := utils.ValidatePepeat(newTask.Repeat); err != nil {
		http.Error(w, `{"error":"Правило повторения указано в неправильном формате"}`, http.StatusBadRequest)
		return
	}

	// Проверка валидности даты.
	// Устанавливаем текущую дату, если поле даты пусто.
	if newTask.Date == "" {
		newTask.Date = time.Now().Format("20060102")
	}

	// Проверка формата даты.
	date, err := time.Parse("20060102", newTask.Date)
	if err != nil {
		http.Error(w, `{"error":"Дата представлена в формате, отличном от 20060102"}`, http.StatusBadRequest)
		return
	}

	// Проеврка актуальности даты.
	today := time.Now().Truncate(24 * time.Hour)
	if date.Before(today) {
		if newTask.Repeat == "" {
			newTask.Date = today.Format("20060102")
		} else {
			nextDate, err := utils.NextDate(time.Now(), newTask.Date, newTask.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Правило повторения указано в неправильном формате"}`, http.StatusBadRequest)
				return
			}
			newTask.Date = nextDate
		}
	}

	// Обновление задачи в базе данных.
	_, err = h.DB.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		newTask.Date, newTask.Title, newTask.Comment, newTask.Repeat, newTask.ID)
	if err != nil {
		http.Error(w, `{"error":"Ошибка обновления данных в базе данных"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}
