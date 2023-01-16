package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// описали названия таблиц в БД
const (
	usersTable      = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemTable  = "lists_item"
)

// Config для подключения к БД нам нужны параметры
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB - возвращает указательно на структуру sqlx.DB и ошибку
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	// проверяем подключение к БД с помощью функции Ping
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
