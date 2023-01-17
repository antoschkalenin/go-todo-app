package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/zhashkevych/todo/model"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, password string) (model.User, error)
}

type TodoList interface {
	Create(userId int, list model.TodoList) (int, error)
	GetAll(userId int) ([]model.TodoList, error)
	GetById(userId, listId int) (model.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input model.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, item model.TodoItem) (int, error)
	GetAll(userId, listId int) ([]model.TodoItem, error)
	GetById(itemId, userId int) (model.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input model.UpdateItemInput) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

// NewRepository конструктор
// иницилизируем репо в конструкторе
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
		TodoList:      NewTodoListRepository(db),
		TodoItem:      NewTodoItemRepository(db),
	}
}
