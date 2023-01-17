package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

// userIdentity метод из заголовка получает токен пользователя и парсит его,
// далее добавляет в конекст id пользвоателя для использования в приложении
func (r *Routes) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	// проверяем что заголовок из конекста запроса не пустой
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "пустой Authorization заголовок")
		return
	}
	// разбиваем заголовок на массив, разделитель пробел
	headerPart := strings.Split(header, " ")
	if len(headerPart) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "не валидный Authorization заголовок")
		return
	}
	// парсим токен и получаем id пользователя
	userId, err := r.services.Authorization.ParseToken(headerPart[1])

	// проверяем что не получили ошибок от парсинга токена
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	// запишем id пользователя в конеткст, это мы делаем для того что бы иметь
	// доступ id пользователя который делает запрос в последующих обработчиках которые
	// вызываются после прослойки middleware
	c.Set(userCtx, userId)
}

// метод проверяет id польщователя из контекста и приводит к типу int
func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "отсутствует id пользователя")
		return 0, errors.New("отсутствует id пользователя")
	}
	// приводим к типу int
	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "id пользователя не валидный")
		return 0, errors.New("id пользователя не валидный")
	}
	return idInt, nil
}
