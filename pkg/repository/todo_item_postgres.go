package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/zhashkevych/todo"
	"strings"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

// NewTodoItemPostgres конструктор
func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

// Create метод создания элемента в списке пользователя
func (r *TodoItemPostgres) Create(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("insert into %s (title, description) values ($1, $2) returning id", todoItemsTable)
	row := r.db.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("insert into %s (list_id, item_id) values ($1, $2)", listsItemTable)
	_, err = r.db.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return itemId, tx.Commit()
}

// GetAll метод возвращает элементы пользователя
func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := fmt.Sprintf("select ti.id, ti.title, ti.description, ti.done "+
		"from %s ti "+
		"inner join %s li on ti.id = li.item_id "+
		"inner join %s ul on ul.list_id = li.list_id "+
		"where ul.user_id = $1 and li.list_id = $2", todoItemsTable, listsItemTable, usersListsTable)
	if err := r.db.Select(&items, query, userId, listId); err != nil {
		return nil, err
	}
	return items, nil
}

// GetById метод возвращает элемент пользователя
func (r *TodoItemPostgres) GetById(itemId, userId int) (todo.TodoItem, error) {
	var item todo.TodoItem
	query := fmt.Sprintf("select ti.id, ti.title, ti.description, ti.done "+
		"from %s ti "+
		"inner join %s li on ti.id = li.item_id "+
		"inner join %s ul on ul.list_id = li.list_id "+
		"where ul.user_id = $1 and ti.id = $2", todoItemsTable, listsItemTable, usersListsTable)
	if err := r.db.Get(&item, query, userId, itemId); err != nil {
		return item, err
	}
	return item, nil
}

// Delete метод удаляет элемент пользователя
func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf("delete from %s ti using %s li, %s ul "+
		"where ti.id = li.item_id and li.list_id = ul.list_id and "+
		"ti.id = $1 and ul.user_id = $2",
		todoItemsTable, listsItemTable, usersListsTable)
	_, err := r.db.Exec(query, itemId, userId)
	return err
}

// Update метод обновляет элемент пользователя
func (r *TodoItemPostgres) Update(userId, itemId int, input todo.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}
	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("update %s ti set %s from %s li, %s ul "+
		"where ti.id = li.item_id and ul.list_id = li.list_id and ti.id = $%d and ul.user_id = $%d",
		todoItemsTable, setQuery, listsItemTable, usersListsTable, argId, argId+1)
	args = append(args, itemId, userId)
	logrus.Debugf("update query: %s", query)
	logrus.Debugf("args query: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
