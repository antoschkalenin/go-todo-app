package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/todo/model"
	"net/http"
)

// структура для логина и пароля от пользователя при аутентификации
type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// отвечат за регистрацию
func (r *Routes) signUp(c *gin.Context) {
	var input model.User

	// делаем бинд JSON и валидируем по правилам в структуре todo.User
	if err := c.BindJSON(&input); err != nil {
		// возвращаем ответ 400 с ошибкой
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := r.services.Authorization.CreateUser(input)
	// если метод создания пользователя вернет ошибку то вернем ответ
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	// ответ JSON со статусом 200 и id нового пользователя
	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

// отвечает за аутентификацию с помощью JWT
func (r *Routes) signIn(c *gin.Context) {
	var input signInInput

	// делаем бинд JSON и валидируем по правилам в структуре signInInput
	if err := c.BindJSON(&input); err != nil {
		// возвращаем ответ 400 с ошибкой
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := r.services.Authorization.GenerateToken(input.Username, input.Password)
	// если метод генерации токена пользователя вернет ошибку то вернем ответ
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ JSON со статусом 200 и JWT token пользователя
	c.JSON(http.StatusOK, map[string]interface{}{"token": token})
}
