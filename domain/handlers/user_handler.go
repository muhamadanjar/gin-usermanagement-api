package handlers

import (
	"net/http"
	"usermanagement-api/domain/models"
	"usermanagement-api/domain/services"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/", CreateUser)
		userGroup.GET("/", GetUsers)
		userGroup.GET("/:id", GetUserByID)
		userGroup.PUT("/:id", UpdateUser)
		userGroup.DELETE("/:id", DeleteUser)
	}
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func GetUsers(c *gin.Context) {

}

func UpdateUser(c *gin.Context) {

}

func DeleteUser(c *gin.Context) {

}

func GetUserByID(c *gin.Context) {

}
