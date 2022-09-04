package v1

import (
	"fmt"
	"golang-rest-api-jwt/common/obj"
	"golang-rest-api-jwt/common/response"
	"golang-rest-api-jwt/dto"
	"golang-rest-api-jwt/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type ProductHandler interface {
	CreateProduct(ctx *gin.Context)

	All(ctx *gin.Context)
}

type productHandler struct {
	productService service.ProductService
	jwtService     service.JWTService
}

func NewProductHandler(productService service.ProductService, jwtService service.JWTService) ProductHandler {
	return &productHandler{
		productService: productService,
		jwtService:     jwtService,
	}
}

func (h *productHandler) CreateProduct(ctx *gin.Context) {
	var createProductRequest dto.CreateProductRequest
	if err := ctx.ShouldBind(&createProductRequest); err != nil {
		response := response.BuildErrorResponse("Failed to process request", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	bearer := strings.Split(authHeader, " ")
	if len(bearer) < 2 {
		// sudah dihandle didalam auth.middleware.go
		return
	}

	bearerToken := bearer[1]
	token := h.jwtService.ValidateToken(bearerToken, ctx)
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])

	productResponse, err := h.productService.CreateProduct(createProductRequest, userID)
	if err != nil {
		response := response.BuildErrorResponse("Failed to process request", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := response.BuildSuccessResponse(true, "Success", productResponse)
	ctx.JSON(http.StatusCreated, response)
}

func (h *productHandler) All(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	bearer := strings.Split(authHeader, " ")
	if len(bearer) < 2 {
		// sudah dihandle didalam auth.middleware.go
		return
	}

	bearerToken := bearer[1]
	token := h.jwtService.ValidateToken(bearerToken, ctx)
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])

	products, err := h.productService.All(userID)
	if err != nil {
		response := response.BuildErrorResponse("Failed to process request", err.Error(), obj.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := response.BuildSuccessResponse(true, "Success", products)
	ctx.JSON(http.StatusOK, response)
}
