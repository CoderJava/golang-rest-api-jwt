package v1

import (
	"golang-rest-api-jwt/common/obj"
	"golang-rest-api-jwt/common/response"
	"golang-rest-api-jwt/dto"
	"golang-rest-api-jwt/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Register(ctx *gin.Context)

	Login(ctx *gin.Context)
}

type authHandler struct {
	userService service.UserService
	authService service.AuthService
	jwtService  service.JWTService
}

func NewAuthHandler(
	userService service.UserService,
	authService service.AuthService,
	jwtService service.JWTService,
) AuthHandler {
	return &authHandler{
		userService: userService,
		authService: authService,
		jwtService:  jwtService,
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

func (h *authHandler) Login(ctx *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := ctx.ShouldBind(&loginRequest); err != nil {
		response := response.BuildErrorResponse("Failed to process request", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if err := h.authService.VerifyCredential(loginRequest.Email, loginRequest.Pasword); err != nil {
		response := response.BuildErrorResponse("Failed to login", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	user, _ := h.userService.FindUserByEmail(loginRequest.Email)

	token := h.jwtService.GenerateToken(strconv.Itoa(int(user.ID)))
	user.AccessToken = token
	response := response.BuildSuccessResponse(true, "Success", user)
	ctx.JSON(http.StatusOK, response)
}
