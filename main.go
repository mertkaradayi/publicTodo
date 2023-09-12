package main

import (
	"examples.com/jwt-auth/controllers"
	"examples.com/jwt-auth/initializers"
	"examples.com/jwt-auth/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/todos", middleware.RequireAuth, controllers.CreateTodo)
	r.GET("/todos", middleware.RequireAuth, controllers.GetTodos)
	r.POST("logout", middleware.RequireAuth, controllers.Logout)
	r.GET("/user", middleware.RequireAuth, controllers.GetUserProfile)
	r.DELETE("/user", middleware.RequireAuth, controllers.DeleteUser)
	r.PUT("/user", middleware.RequireAuth, controllers.ChangePassword)
	r.PUT("/todos/:id", middleware.RequireAuth, controllers.UpdateTodo)
	r.DELETE("/todos/:id", middleware.RequireAuth, controllers.DeleteTodo)
	return r
}

func main() {

	r := SetupRouter()
	r.Run()
}
