package service

import (
	"github.com/zhashkevych/todo"
	"github.com/zhashkevych/todo/pkg/repository"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId, listId int) (todo.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input todo.UpdateListInput) error
}

type TodoItem interface {
	Create(userId, listId int, item todo.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo.TodoItem, error)
	GetById(itemId, userId int) (todo.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input todo.UpdateItemInput) error
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
