package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/todo"
	"net/http"
	"strconv"
)

// метод создает элемент
func (h *Handler) createItem(c *gin.Context) {
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
	var item todo.TodoItem
	if err := c.BindJSON(&item); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "не валидный item")
		return
	}
	// создаем элемент
	id, err := h.services.TodoItem.Create(userId, listId, item)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ клиенту
	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

// метод возвращает список элементов пользователя
func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "не валидный параметр id списка")
		return
	}
	items, err := h.services.TodoItem.GetAll(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ клиенту
	c.JSON(http.StatusOK, items)
}

// метод возващает элемент пользователя
func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "не валидный параметр id элемента")
		return
	}
	item, err := h.services.TodoItem.GetById(itemId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ клиенту
	c.JSON(http.StatusOK, item)
}

// метод обновляет элемент пользователя
func (h *Handler) updateItemById(c *gin.Context) {
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
	var input todo.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err = h.services.TodoItem.Update(userId, itemId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ клиенту
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

// метод удаляет элемент пользователя
func (h *Handler) deleteItemById(c *gin.Context) {
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
	err = h.services.TodoItem.Delete(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// ответ клиенту
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
