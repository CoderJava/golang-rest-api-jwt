package v1

import (
	"golang-rest-api-jwt/common/obj"
	"golang-rest-api-jwt/common/response"
	"golang-rest-api-jwt/dto"
	"golang-rest-api-jwt/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Register(ctx *gin.Context)
}

type authHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) AuthHandler {
	return &authHandler{
		userService: userService,
	}
}

func (h *authHandler) Register(ctx *gin.Context) {
	var registerRequest dto.RegisterRequest

	if err := ctx.ShouldBind(&registerRequest); err != nil {
		response := response.BuildErrorResponse("Failed to process request", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	user, err := h.userService.CreateUser(registerRequest)
	if err != nil {
		response := response.BuildErrorResponse(err.Error(), err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := response.BuildSuccessResponse(true, "Success", user)
	ctx.JSON(http.StatusCreated, response)
}