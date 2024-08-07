package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AramLab/todo-list/models"
	"github.com/AramLab/todo-list/storage"
)

type Handlers struct {
	DB *sql.DB
}

func NewHandlers(db *sql.DB) *Handlers {
	return &Handlers{DB: db}
}

func (h *Handlers) NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var task models.Task
		now := r.URL.Query().Get("now")
		task.Date = r.URL.Query().Get("date")
		task.Repeat = r.URL.Query().Get("repeat")

		timeNow, err := time.Parse(models.DatePattern, now)
		if err != nil {
			return
		}

		nextDate, err := task.NextDate(timeNow)
		if err != nil {
			return
		}

		response := nextDate

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(response))
		if err != nil {
			log.Printf("Ошибка записи в ответ: %v", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handlers) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := storage.AddTask(&task, h.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"id": id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Ошибка записи в ответ: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

func (h *Handlers) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := storage.GetTasks(h.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Если tasks == nil, инициализируйте его как пустой слайс
	if tasks == nil {
		tasks = []models.Task{}
	}

	response := models.TasksResponse{Tasks: tasks}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Printf("Ошибка запипси в ответ: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}
	var task models.Task
	task, err := storage.GetTask(id, h.DB)
	if err != nil {
		http.Error(w, `{"error":"Ошибка чтения данных из базы данных"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Printf("Ошибка запипси в ответ: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task
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
	if _, err := storage.GetTask(newTask.ID, h.DB); err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusBadRequest)
		return
	}

	// Проверка валидности заголовка.
	if err := newTask.ValidateTitle(); err != nil {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Проверка валидности правила повторения.
	if err := newTask.ValidatePepeat(); err != nil {
		http.Error(w, `{"error":"Правило повторения указано в неправильном формате"}`, http.StatusBadRequest)
		return
	}

	// Проверка валидности даты.
	if err := newTask.ValidateAndFormatDate(&newTask); err != nil {
		http.Error(w, `{"error":"дата представлена в формате, отличном от 20060102"}`, http.StatusBadRequest)
		return
	}

	// Обновляем задачу в базе данных.
	err = storage.UpdateTask(&newTask, h.DB)
	if err != nil {
		http.Error(w, `{"error":"Ошибка обновления данных в базе данных"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{}`))
	if err != nil {
		log.Printf("Ошибка запипси в ответ: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	// Проверка наличия `id`.
	if id == "" {
		http.Error(w, `{"error":"Не указан ID задачи"}`, http.StatusBadRequest)
		return
	}

	// Получаем задачу по `id`.
	task, err := storage.GetTask(id, h.DB)
	if err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusBadRequest)
		return
	}

	// Рассчитываем следующее время для периодической задачи.
	if task.Repeat != "" {
		task.Date, err = task.NextDate(time.Now())
		if err != nil {
			http.Error(w, `{"error":"Ошибка в работе функции NextDate()"}`, http.StatusInternalServerError)
			return
		}

		// Обновляем дату задачи.
		err = storage.UpdateTaskDate(id, &task, h.DB)
		if err != nil {
			http.Error(w, `{"error":"Ошибка обновления даты задачи"}`, http.StatusInternalServerError)
			return
		}
		// Удаляем одноразовую задачу.
	} else {
		err := storage.DeleteTaskById(id, h.DB)
		if err != nil {
			http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{}`))
	if err != nil {
		log.Printf("Ошибка запипси в ответ: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// Проверка наличия `id`.
	if id == "" {
		http.Error(w, `{"error":"Не указан ID задачи"}`, http.StatusBadRequest)
		return
	}

	// Получаем задачу по `id` и удостоверяемся в её наличии.
	_, err := storage.GetTask(id, h.DB)
	if err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusBadRequest)
		return
	}

	// Удаляем задачу из базы данных.
	err = storage.DeleteTaskById(id, h.DB)
	if err != nil {
		http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{}`))
	if err != nil {
		log.Printf("Ошибка запипси в ответ: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
}
