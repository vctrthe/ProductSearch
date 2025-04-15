package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vctrthe/ProductSearch/config"
	"github.com/vctrthe/ProductSearch/controller"
	"github.com/vctrthe/ProductSearch/util"
)

func main() {
	config.InitElastic()
	util.LoadAndIndexData(config.ES)

	router := gin.Default()
	router.GET("/search", controller.SearchProduct)
	router.Run(":8080")
}
