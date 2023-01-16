package service

import (
	"errors"
	"github.com/zhashkevych/todo"
	"github.com/zhashkevych/todo/pkg/repository"
)

// TodoItemService структура для сервиса в котором будем использовать repo
type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

// NewTodoItemService конструктор
func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

// Create метод создания элемента в списке пользователя
func (s *TodoItemService) Create(userId, listId int, item todo.TodoItem) (int, error) {
	// проверка по user id
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, errors.New("нет списка или он не пренадлежит данному пользователю")
	}
	return s.repo.Create(listId, item)
}

// GetAll метод возвращает элементы пользователя
func (s *TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

// GetById метод возвращает элемент пользователя
func (s *TodoItemService) GetById(itemId, userId int) (todo.TodoItem, error) {
	return s.repo.GetById(itemId, userId)
}

// Delete метод удаляет элемент пользователя
func (s *TodoItemService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

// Update метод обновляет элемент пользователя
func (s *TodoItemService) Update(userId, itemId int, input todo.UpdateItemInput) error {
	// валидация структуры
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, itemId, input)
}
