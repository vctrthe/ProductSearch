package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vctrthe/ProductSearch/config"
	"github.com/vctrthe/ProductSearch/controller"
	_ "github.com/vctrthe/ProductSearch/docs"
	"github.com/vctrthe/ProductSearch/util"
)

// @title Product Search API
// @version 1.0
// @description This is a product search API using Elasticsearch
// @host localhost:8080
// @BasePath /
func main() {
	config.InitElastic()
	util.LoadAndIndexData(config.ES)

	router := gin.Default()

	// Swagger documentation endpoint
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/search", controller.SearchProduct)
	router.Run(":8080")
}
