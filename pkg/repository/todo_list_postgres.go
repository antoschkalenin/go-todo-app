package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/zhashkevych/todo"
	"strings"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

// NewTodoListPostgres конструктор
func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

// Create метод создания списка, метод транзакционный
func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	// для использования транзакции используем метод Begin
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var listsTableId int
	createListQuery := fmt.Sprintf("insert into %s (title, description) values ($1, $2) returning id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	// используя метод Scan записываем в переменную значение из ответа
	if err := row.Scan(&listsTableId); err != nil {
		tx.Rollback()
		return 0, err
	}
	createUsersListQuery := fmt.Sprintf("insert into %s (user_id, list_id) values ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, listsTableId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return listsTableId, tx.Commit()
}

// GetAll возвращает все списки пользователя
func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	query := fmt.Sprintf("select tl.id, tl.title, tl.description from %s tl "+
		"inner join %s ul on tl.id = ul.list_id where ul.user_id = $1",
		todoListsTable, usersListsTable)
	// r.db.Select - работает аналогично r.db.Get только используется
	//для выборки больше одного элемента и записи в slice
	err := r.db.Select(&lists, query, userId)
	return lists, err
}

// GetById возвращает список пользователя по listId
func (r *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList
	query := fmt.Sprintf("select tl.id, tl.title, tl.description from %s tl "+
		"inner join %s ul on tl.id = ul.list_id where ul.user_id = $1 and tl.id = $2",
		todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)
	return list, err
}

// Delete удаляет список пользователя по listId
func (r *TodoListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("delete from %s tl using %s ul "+
		"where tl.id = ul.list_id and ul.list_id = $1 and ul.user_id = $2",
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, listId, userId)
	return err
}

// Update обновляет список пользователя по listId
func (r *TodoListPostgres) Update(userId, listId int, input todo.UpdateListInput) error {
	// иницилизуерм slice строк, интрфейсов
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	// иницилизуерм аргумент = 1
	argId := 1
	// проверка полей, если не nil то будем добавлять
	// в список элемкнты для формирования запросов в БД
	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		// в slice аргументов добавим само значение title
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}
	// соединим наши элементы slice в одному строку
	// пример результата: title=$1, description=$2
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("update %s tl set %s from %s ul where tl.id = ul.list_id and ul.list_id = $%d and ul.user_id = $%d",
		todoListsTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, listId, userId)
	logrus.Debugf("update query: %s", query)
	logrus.Debugf("args query: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
