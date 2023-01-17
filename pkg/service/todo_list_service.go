package service

import (
	"github.com/zhashkevych/todo/model"
	"github.com/zhashkevych/todo/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

// NewTodoListService конструктор
func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

// Create метод создания списка
// userId - id пользователя
// list - данные списка от пользователя
func (s *TodoListService) Create(userId int, list model.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

// GetAll возвращает все списки пользователя
func (s *TodoListService) GetAll(userId int) ([]model.TodoList, error) {
	return s.repo.GetAll(userId)
}

// GetById возвращает список пользователя по listId
func (s *TodoListService) GetById(userId, listId int) (model.TodoList, error) {
	return s.repo.GetById(userId, listId)
}

// Delete удаляет список пользователя по listId
func (s *TodoListService) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}

// Update обновляет список пользователя по listId
func (s *TodoListService) Update(userId, listId int, input model.UpdateListInput) error {
	// валидация структуры
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, input)
}
