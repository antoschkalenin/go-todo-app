package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/todo/model"
	"net/http"
	"strconv"
)

// структура для ответа пользвоателю
type getAllListsResponse struct {
	Data []model.TodoList `json:"data"`
}

// метод создания списка
func (r *Routes) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	// получение от пользователя данных и валидация
	var input model.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// создаем список
	id, err := r.services.TodoList.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// ответ клиенту
	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

// метод возвращает все списки пользователя
func (r *Routes) getAllLists(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	// получение списка пользователя
	lists, err := r.services.TodoList.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ клиенту
	c.JSON(http.StatusOK, getAllListsResponse{Data: lists})
}

// метод возвращает список по listId
func (r *Routes) getListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	// получение id списка из пути запроса
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "не валидный параметр id списка")
		return
	}
	// получение списка по userId и listId
	list, err := r.services.TodoList.GetById(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ клиенту
	c.JSON(http.StatusOK, list)
}

// метод обновляет список пользователя по lisId
func (r *Routes) updateListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	// получение id списка из пути запроса
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "не валидный параметр id списка")
		return
	}
	// получение от пользователя данных и валидация
	var input model.UpdateListInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err = r.services.TodoList.Update(userId, listId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ клиенту
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

// метод удаляет список пользователя по listId
func (r *Routes) deleteListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	// получение id списка из пути запроса
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "не валидный параметр id списка")
		return
	}
	// удаление по userId и listId
	err = r.services.TodoList.Delete(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ клиенту
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
