package model

import "errors"

// TodoList - сущность задания
// json - это json тег для корректного вывода/ввода http ответа/запроса
// binding - валидация полей (часть фреймворка gin)
// db - дает возможность выборки из БД
type TodoList struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

// UsersList сущность описывает список заданий у пользователя
type UsersList struct {
	Id     int
	UserId int
	ListId int
}

// TodoItem сущность задания
type TodoItem struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}

// ListsItem сущность со списком заданий
type ListsItem struct {
	Id     int
	ListId int
	ItemId int
}

// UpdateListInput для обновления списка
// *string - добавили тип указателя на строку и если не будет
// иметь значения то будет равно не пустой строке "" по умолчанию а nil
type UpdateListInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (i UpdateListInput) Validate() error {
	if i.Title == nil && i.Description == nil {
		return errors.New("нет значений для обновления списка")
	}
	return nil
}

type UpdateItemInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (i UpdateItemInput) Validate() error {
	if i.Title == nil && i.Description == nil && i.Done == nil {
		return errors.New("нет значений для обновления элемента")
	}
	return nil
}
