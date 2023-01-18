package controller

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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

// Authorize авторизация пользователя с помощью casbin
func Authorize(obj string, act string, adapter persist.Adapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, existed := c.Get(userCtx)
		if !existed {
			c.AbortWithStatusJSON(401, errorResponse{"пользователь еще не авторизовался"})
			return
		}
		idInt, ok := userId.(int)
		if !ok {
			c.AbortWithStatusJSON(500, errorResponse{"ошибка при конвертации user id"})
		}
		ok, err := enforce(strconv.Itoa(idInt), obj, act, adapter)
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatusJSON(500, errorResponse{"ошибка при авторизации пользователя"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(403, errorResponse{"доступ закрыт"})
			return
		}
		c.Next()
	}
}

// метод создает enforce, загружает политики и роли
func enforce(sub, obj, act string, adapter persist.Adapter) (bool, error) {
	enforcer, err := casbin.NewEnforcer("configs/rbac_model.conf", adapter)
	if err != nil {
		return false, fmt.Errorf("не удалось создать casbin Enforcer: %w", err)
	}
	// динамическая загрузка политик из БД
	err = enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("не удалось загрузить политику из БД: %w", err)
	}
	ok, err := enforcer.Enforce(sub, obj, act)
	return ok, err
}
