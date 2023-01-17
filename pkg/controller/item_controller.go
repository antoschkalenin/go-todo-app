package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/todo/model"
	"net/http"
	"strconv"
)

// метод создает элемент
func (r *Routes) createItem(c *gin.Context) {
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
	// бинд входящих значений с валидацией
	var item model.TodoItem
	if err := c.BindJSON(&item); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "не валидный item")
		return
	}
	// создаем элемент
	id, err := r.services.TodoItem.Create(userId, listId, item)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ клиенту
	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

// метод возвращает список элементов пользователя
func (r *Routes) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "не валидный параметр id списка")
		return
	}
	items, err := r.services.TodoItem.GetAll(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ клиенту
	c.JSON(http.StatusOK, items)
}

// метод возващает элемент пользователя
func (r *Routes) getItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "не валидный параметр id элемента")
		return
	}
	item, err := r.services.TodoItem.GetById(itemId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ клиенту
	c.JSON(http.StatusOK, item)
}

// метод обновляет элемент пользователя
func (r *Routes) updateItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	// получение id элемента из пути запроса
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "не валидный параметр id элемента")
		return
	}
	// получение от пользователя данных и валидация
	var input model.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err = r.services.TodoItem.Update(userId, itemId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ клиенту
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

// метод удаляет элемент пользователя
func (r *Routes) deleteItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	// получение id элемента из пути запроса
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "не валидный параметр id элемента")
		return
	}
	// удаление по userId и itemId
	err = r.services.TodoItem.Delete(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ клиенту
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
