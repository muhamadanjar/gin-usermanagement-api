package handlers

import (
	"net/http"
	"usermanagement-api/internal/user/dto"
	"usermanagement-api/internal/user/usecase"
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
		userUsecase usecase.UserUsecase
	}
)

func NewUserHandler(us usecase.UserUsecase) UserHandler {
	return &userHandler{
		userUsecase: us,
	}
}

func (c *userHandler) CreateUser(ctx *gin.Context) {
	var user dto.CreateUserRequest
	if err := ctx.ShouldBind(&user); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
}
