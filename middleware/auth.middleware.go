package middleware

import (
	"fmt"
	"golang-rest-api-jwt/common/response"
	"golang-rest-api-jwt/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response := response.BuildErrorResponse("Failed to process request", "No token provided", nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		bearer := strings.Split(authHeader, " ")
		if len(bearer) < 2 {
			response := response.BuildErrorResponse("Failed to process request", "No token provided", nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		bearerToken := bearer[1]
		token := jwtService.ValidateToken(bearerToken, ctx)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			fmt.Println("Claim[user_id]: ", claims["user_id"])
			fmt.Println("Claim[iss]: ", claims["iss"])
		} else {
			response := response.BuildErrorResponse("Error", "Your token is not valid", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
	}
}
