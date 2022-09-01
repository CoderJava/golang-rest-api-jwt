package service

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateToken(userID string) string

	ValidateToken(token string, ctx *gin.Context) *jwt.Token
}

type jwtCustomClaim struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	return &jwtService{
		issuer:    "admin",
		secretKey: getSecretKey(),
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey != "" {
		secretKey = "admin"
	}
	return secretKey
}

func (s *jwtService) GenerateToken(userID string) string {
	expirationTime := time.Now().Add(2 * time.Minute)

	claims := &jwtCustomClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: s.issuer,
			ExpiresAt: &jwt.NumericDate{
				Time: expirationTime,
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (s *jwtService) ValidateToken(token string, ctx *gin.Context) *jwt.Token {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, isOk := t.Method.(*jwt.SigningMethodHMAC); !isOk {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil
	}

	return t
}
