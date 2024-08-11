package task

import (
	"errors"
	"time"

	"github.com/AramLab/todo-list/internal/entities"
	"github.com/AramLab/todo-list/internal/storage"
)

type TaskService struct {
	TaskRepository storage.TaskRepository
}

// Create task.
func (uc *TaskService) Create(task *entities.Task) (int64, error) {
	// Проверка заголовка.
	if err := task.ValidateTitle(); err != nil {
		return 0, errors.New(`{"error":"Не указан заголовок задачи"}`)
	}

	// Проверка даты.
	if err := task.ValidateAndFormatDate(); err != nil {
		return 0, errors.New(`{"error":"Дата представлена в формате, отличном от 20060102"}`)
	}

	// Проверка шаблонов правил.
	if err := task.ValidatePepeat(); err != nil {
		return 0, errors.New(`{"error":"Правило повторения указано в неправильном формате"}`)
	}

	id, err := uc.TaskRepository.Save(task)
	if err != nil {
		return 0, errors.New(`{"error":"Ошибка создания задачи"}`)
	}
	return id, nil
}

// Get task.
func (uc *TaskService) Get(id string) (entities.Task, error) {
	task, err := uc.TaskRepository.FindById(id)
	if err != nil {
		return entities.Task{}, errors.New(`{"error":"Ошибка получения задачи"}`)
	}
	return task, nil
}

// Get all tasks.
func (uc *TaskService) GetAll() ([]entities.Task, error) {
	var tasks []entities.Task
	tasks, err := uc.TaskRepository.FindAll()
	if err != nil {
		return []entities.Task{}, errors.New(`{"error":"Ошибка получения задач"}`)
	}
	return tasks, nil
}

// Update by id.
func (uc *TaskService) UpdateDate(id string, task *entities.Task) error {
	// Проверка валидности даты.
	if err := task.ValidateAndFormatDate(); err != nil {
		return errors.New(`{"error":"дата представлена в формате, отличном от 20060102"}`)
	}

	err := uc.TaskRepository.UpdateDate(id, task)
	if err != nil {
		return errors.New(`{"error":"Ошибка обновления даты задачи"}`)
	}
	return nil
}

// Update task.
func (uc *TaskService) Update(task *entities.Task) error {
	// Проверка наличия `id`.
	if task.ID == "" {
		return errors.New(`{"error":"Не указан ID задачи"}`)
	}

	// Проверка существования задачи по данному ID.
	if _, err := uc.Get(task.ID); err != nil {
		return errors.New(`{"error":"Задача не найдена"}`)
	}

	// Проверка валидности заголовка.
	if err := task.ValidateTitle(); err != nil {
		return errors.New(`{"error":"Не указан заголовок задачи"}`)
	}

	// Проверка валидности правила повторения.
	if err := task.ValidatePepeat(); err != nil {
		return errors.New(`{"error":"Правило повторения указано в неправильном формате"}`)
	}

	// Проверка валидности даты.
	if err := task.ValidateAndFormatDate(); err != nil {
		return errors.New(`{"error":"дата представлена в формате, отличном от 20060102"}`)
	}

	// Обновляем задачу в базе данных.
	err := uc.TaskRepository.Update(task)
	if err != nil {
		return errors.New(`{"error":"Ошибка обновления задачи"}`)
	}
	return nil
}

// Delete task by id.
func (uc *TaskService) Delete(id string) error {
	// Проверка наличия `id`.
	if id == "" {
		return errors.New(`{"error":"Не указан ID задачи"}`)
	}

	// Получаем задачу по `id` и удостоверяемся в её наличии.
	_, err := uc.Get(id)
	if err != nil {
		return errors.New(`{"error":"Задача не найдена"}`)
	}

	// Удаляем задачу из базы данных.
	err = uc.TaskRepository.DeleteById(id)
	if err != nil {
		return errors.New(`{"error":"Ошибка удаления задачи"}`)
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------

// func Complete()
// Complete task.
// Завершает задачу и исходя из ее свойств либо удаляет её, либо переносит на следующую дату.
func (uc *TaskService) Complete(id string) error {
	// Проверка наличия `id`.
	if id == "" {
		return errors.New(`{"error":"Не указан ID задачи"}`)
	}

	// Получаем задачу по `id`.
	task, err := uc.Get(id)
	if err != nil {
		return errors.New(`{"error":"Задача не найдена"}`)
	}

	// Рассчитываем следующее время для периодической задачи.
	if task.Repeat != "" {
		task.Date, err = task.NextDate(time.Now())
		if err != nil {
			return errors.New(`{"error":"Ошибка в работе функции NextDate()"}`)
		}

		// Обновляем дату задачи.
		err := uc.UpdateDate(id, &task)
		if err != nil {
			return errors.New(`{"error":"Ошибка обновления даты задачи"}`)
		}
		// Удаляем одноразовую задачу.
	} else {
		err := uc.Delete(id)
		if err != nil {
			return errors.New(`{"error":"Ошибка удаления задачи"}`)
		}
	}
	return nil
}

func NewTaskService(TaskRepository storage.TaskRepository) *TaskService {
	return &TaskService{TaskRepository: TaskRepository}
}
