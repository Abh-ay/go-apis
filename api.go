package main

import (
	loginHandler "go-apis/Handler"
	middleware "go-apis/Middleware"
	util "go-apis/Util"
	dbConnection "go-apis/connection"

	"github.com/gin-gonic/gin"
)

func main() {
	dbConnection.ConnectDB()
	r := gin.Default()
	publicRoutes := r.Group("/public")
	{
		publicRoutes.POST("/login", loginHandler.Login)
		publicRoutes.POST("/register", loginHandler.Register)
		publicRoutes.GET("/getkeys", util.KeyWrite)
	}
	protectedRoutes := r.Group("/protected")
	protectedRoutes.Use(middleware.AuthenticationMiddleware())
	{
		// Protected routes here
	}

	r.Run(":8080")
}
