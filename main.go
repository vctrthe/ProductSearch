package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vctrthe/ProductSearch/controller"
)

func main() {
	router := gin.Default()
	router.GET("/search", controller.SearchProduct)
	router.Run(":8080")
}
