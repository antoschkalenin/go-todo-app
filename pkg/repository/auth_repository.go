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
	query := fmt.Sprintf("insert into %s (name, username, password_hash, is_ad) values ($1,$2,$3,$4) returning id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password, user.IsAd)
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

// GetUserByUsernameAndByAd получаем пользователя по логину и признаку isAd = true/false
func (r *AuthRepository) GetUserByUsernameAndByAd(username string, isAd bool) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("select id from %s where username = $1 and is_ad = $2", usersTable)
	// передаем указатель &user на структуру куда записываем ответ
	err := r.db.Get(&user, query, username, isAd)
	return user, err
}

// UpdatePassword обновляет hash пароля
func (r *AuthRepository) UpdatePassword(userId int, password string) error {
	query := fmt.Sprintf("update %s set password_hash = $1 where id = $2", usersTable)
	_, err := r.db.Exec(query, password, userId)
	return err
}
