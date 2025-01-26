package handlers

import (
	"net/http"
	"usermanagement-api/domain/models"
	"usermanagement-api/services"
	"usermanagement-api/utils"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, userHandler UserHandler) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/", userHandler.CreateUser)
	}
}

type (
	UserHandler interface {
		CreateUser(ctx *gin.Context)
	}

	userHandler struct {
		userService services.UserService
	}
)

func NewUserHandler(us services.UserService) UserHandler {
	return &userHandler{
		userService: us,
	}
}

func (c *userHandler) CreateUser(ctx *gin.Context) {
	var user models.UserCreateRequest
	if err := ctx.ShouldBind(&user); err != nil {
		res := utils.BuildResponseFailed(models.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
}
