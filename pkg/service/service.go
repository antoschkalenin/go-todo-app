package service

import (
	"github.com/zhashkevych/todo/model"
	"github.com/zhashkevych/todo/pkg/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list model.TodoList) (int, error)
	GetAll(userId int) ([]model.TodoList, error)
	GetById(userId, listId int) (model.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input model.UpdateListInput) error
}

type TodoItem interface {
	Create(userId, listId int, item model.TodoItem) (int, error)
	GetAll(userId, listId int) ([]model.TodoItem, error)
	GetById(itemId, userId int) (model.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input model.UpdateItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

// NewService конструктор
// иницилизируем сервис в конструкторе
func NewService(r *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(r.Authorization),
		TodoList:      NewTodoListService(r.TodoList),
		TodoItem:      NewTodoItemService(r.TodoItem, r.TodoList),
	}
}
