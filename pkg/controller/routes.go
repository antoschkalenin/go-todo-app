package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/todo/pkg/service"
)

type Routes struct {
	services *service.Service
}

func NewRouters(services *service.Service) *Routes {
	return &Routes{services: services}
}

// InitRoutes - иницилизируем маршруты приложения
func (r *Routes) InitRoutes() *gin.Engine {
	// иницилизация роутера gin
	router := gin.New()

	// сгруппируем по маршрутам наши методы
	// группа end points для регистрации и авторизации
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", r.signUp)
		auth.POST("/sign-in", r.signIn)
	}

	// группа end points для работы со списками и их задачами,
	// так же установили прослойку middleware для группы /api указав метод userIdentity.
	// В userIdentity делаем проверку токена от пользователя его парсинг и созранение в контекст
	api := router.Group("/api", r.userIdentity)
	{
		// группа lists
		lists := api.Group("/lists")
		{
			// :id - bind на id
			lists.POST("/", r.createList)
			lists.GET("/", r.getAllLists)
			lists.GET("/:id", r.getListById)
			lists.PUT("/:id", r.updateListById)
			lists.DELETE("/:id", r.deleteListById)

			// группа для элементов списка
			items := lists.Group(":id/items")
			{
				items.POST("/", r.createItem)
				items.GET("/", r.getAllItems)
			}
		}
		// группа для элементов списка для изменения (нет необходимости принимать list id каждый раз)
		items := api.Group("/items")
		{
			items.GET("/:id", r.getItemById)
			items.PUT("/:id", r.updateItemById)
			items.DELETE("/:id", r.deleteItemById)
		}
	}
	return router
}
