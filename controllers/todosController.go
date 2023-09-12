package controllers

import (
	"net/http"

	"examples.com/jwt-auth/initializers"
	"examples.com/jwt-auth/models"
	"github.com/gin-gonic/gin"
)

func CreateTodo(c *gin.Context) {
	// Get the user from the context (added by the authentication middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get the Status and Description from the request body
	var body struct {
		Status      string
		Description string
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Create the todo associated with the user's ID
	todo := models.Todo{
		UserID:      user.(models.User).ID, // Use the user's ID as the foreign key
		Status:      body.Status,
		Description: body.Description,
	}

	result := initializers.DB.Create(&todo)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create todo",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo created successfully",
	})
}

func GetTodos(c *gin.Context) {
	// Get the user from the context (added by the authentication middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Find todos associated with the user's ID (excluding User relationship)
	var todos []models.Todo
	initializers.DB.Where("user_id = ?", user.(models.User).ID).Find(&todos)

	c.JSON(http.StatusOK, todos)
}

func UpdateTodo(c *gin.Context) {
	// Get the user from the context (added by the authentication middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get the todo ID from the URL parameter
	todoID := c.Param("id")

	// Find the todo by ID
	var todo models.Todo
	result := initializers.DB.First(&todo, todoID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Todo not found",
		})
		return
	}

	// Check if the todo belongs to the authenticated user
	if todo.UserID != user.(models.User).ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to update this todo",
		})
		return
	}

	// Get the updated Status and Description from the request body
	var body struct {
		Status      string
		Description string
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Update the todo with the new values
	todo.Status = body.Status
	todo.Description = body.Description
	result = initializers.DB.Save(&todo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update todo",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo updated successfully",
	})
}

func DeleteTodo(c *gin.Context) {
	// Get the user from the context (added by the authentication middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get the todo ID from the URL parameter
	todoID := c.Param("id")

	// Find the todo by ID
	var todo models.Todo
	result := initializers.DB.First(&todo, todoID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Todo not found",
		})
		return
	}

	// Check if the todo belongs to the authenticated user
	if todo.UserID != user.(models.User).ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to delete this todo",
		})
		return
	}

	// Delete the todo
	result = initializers.DB.Delete(&todo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete todo",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo deleted successfully",
	})
}
