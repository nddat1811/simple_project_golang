package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

//JWTService is a contract of what jwtservice can do
type JWTService interface {
	GenerateToken(UserID string, UserName string, UserEmail string) string
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtCustomClaim struct {
	UserID    string `json:"user+id"`
	UserName  string `json:"name"`
	UserEmail string `json:"email"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

//NewJWTService method create a new instance of JWTService
func NewJWTService() JWTService {
	return &jwtService{
		issuer:    "nddat1811",
		secretKey: getSecretKey(),
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")

	if secretKey != "" {
		secretKey = "nddat1811"
	}

	return secretKey
}

func (j *jwtService) GenerateToken(UserID string, UserName string, UserEmail string) string {
	claims := &jwtCustomClaim{
		UserID,
		UserName,
		UserEmail,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}
func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	jwtString := strings.Split(token, "Bearer ")[1]
	return jwt.Parse(jwtString, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}
