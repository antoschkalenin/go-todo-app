package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/todo/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// InitRoutes - иницилизируем маршруты приложения
func (h *Handler) InitRoutes() *gin.Engine {
	// иницилизация роутера gin
	router := gin.New()

	// сгруппируем по маршрутам наши методы
	// группа end points для регистрации и авторизации
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	// группа end points для работы со списками и их задачами,
	// так же установили прослойку middleware для группы /api указав метод userIdentity.
	// В userIdentity делаем проверку токена от пользователя его парсинг и созранение в контекст
	api := router.Group("/api", h.userIdentity)
	{
		// группа lists
		lists := api.Group("/lists")
		{
			// :id - bind на id
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateListById)
			lists.DELETE("/:id", h.deleteListById)

			// группа для элементов списка
			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
			}
		}
		// группа для элементов списка для изменения (нет необходимости принимать list id каждый раз)
		items := api.Group("/items")
		{
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItemById)
			items.DELETE("/:id", h.deleteItemById)
		}
	}
	return router
}
