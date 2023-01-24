package model

// User описание сущности пользователя
// поля совпадают по стркутре с БД
// json - это json тег для корректного вывода/ввода http ответа/запроса
// binding - валидация полей (часть фреймворка gin)
// db - для описания полей в БД (необходимо для работы метода r.db.Get в репозитории)
type User struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	IsAd     bool   `json:"isId" db:"is_ad"`
}
