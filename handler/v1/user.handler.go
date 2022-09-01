package v1

import (
	"fmt"
	"golang-rest-api-jwt/common/obj"
	"golang-rest-api-jwt/common/response"
	"golang-rest-api-jwt/dto"
	"golang-rest-api-jwt/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type UserHandler interface {
	Profile(ctx *gin.Context)

	Update(ctx *gin.Context)
}

type userHandler struct {
	userService service.UserService
	jwtService  service.JWTService
}

func NewUserHandler(userService service.UserService, jwtService service.JWTService) UserHandler {
	return &userHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (h *userHandler) Profile(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	bearerToken := strings.Split(header, " ")
	if len(bearerToken) < 2 {
		response := response.BuildErrorResponse("Error", "Token not provided", obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusForbidden, response)
	}

	token := h.jwtService.ValidateToken(bearerToken[1], ctx)

	if token == nil {
		response := response.BuildErrorResponse("Error", "Failed to validate token", obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
	}

	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	user, err := h.userService.FindUserByID(id)

	if err != nil {
		response := response.BuildErrorResponse("Error", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	response := response.BuildSuccessResponse(true, "Success", user)
	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) Update(ctx *gin.Context) {
	var updateUserRequest dto.UpdateUserRequest

	if err := ctx.ShouldBind(&updateUserRequest); err != nil {
		response := response.BuildErrorResponse("Failed to process request", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	strUserID := h.getUserIDByHeader(ctx)
	if strUserID == "" {
		response := response.BuildErrorResponse("Error", "Failed to validate token", obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
	}

	userID, err := strconv.ParseInt(strUserID, 0, 64)
	if err != nil {
		response := response.BuildErrorResponse("Error", "Failed to convert user ID to integer", obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
	}

	updateUserRequest.ID = userID
	resultUpdateUser, err := h.userService.UpdateUser(updateUserRequest)
	if err != nil {
		response := response.BuildErrorResponse("Error", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
	}

	response := response.BuildSuccessResponse(true, "Success", resultUpdateUser)
	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) getUserIDByHeader(ctx *gin.Context) string {
	header := ctx.GetHeader("Authorization")
	bearerToken := strings.Split(header, " ")
	if len(bearerToken) < 2 {
		response := response.BuildErrorResponse("Error", "Token not provided", obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusForbidden, response)
	}

	token := h.jwtService.ValidateToken(bearerToken[1], ctx)
	if token == nil {
		response := response.BuildErrorResponse("Error", "Failed to validate token", obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
	}

	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
