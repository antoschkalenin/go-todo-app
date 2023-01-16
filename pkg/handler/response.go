package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// errorResponse структура для ответа с ошибокой
type errorResponse struct {
	Message string `json:"message"`
}

// statusResponse структура для овтета со статусом
type statusResponse struct {
	Status string `json:"status"`
}

// функция для создания ответов с ошибками
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	// блокирует последующее выполнение обработчиков
	//и записывает ответ код и сообщение в ответ
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
