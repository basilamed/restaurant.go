package router

import (
	"github.com/gin-gonic/gin"
	"items/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	authorized := router.Group("/show")
	{
		authorized.GET("/items/getAll", controllers.GetItems)
		authorized.POST("/items/addItem", controllers.AddItem)
		authorized.GET("/items/getItemById/:id", controllers.GetItemById)
		authorized.PUT("/items/updateItem/:id", controllers.UpdateItemById)
		authorized.DELETE("/items/deleteItem/:id", controllers.DeleteItemById)
	}

	return router

}
