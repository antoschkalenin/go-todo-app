package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zhashkevych/todo/model"
)

type AuthRepository struct {
	db *sqlx.DB
}

// NewAuthRepository конструктор для работы со слоем БД
func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser создание пользователя в БД
func (r *AuthRepository) CreateUser(user model.User) (int, error) {
	var id int
	query := fmt.Sprintf("insert into %s (name, username, password_hash) values ($1,$2,$3) returning id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	// используя метод Scan записываем в переменную значение из ответа
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// GetUser получаем пользователя по логину и паролю
func (r *AuthRepository) GetUser(username, password string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("select id from %s where username = $1 and password_hash = $2", usersTable)
	// передаем указатель &user на структуру куда записываем ответ
	err := r.db.Get(&user, query, username, password)
	return user, err
}
