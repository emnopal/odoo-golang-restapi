package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	jwt_schema "github.com/emnopal/odoo-golang-restapi/app/schemas/auth"
	utils "github.com/emnopal/odoo-golang-restapi/app/utils/env"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	DEFAULT_JWT_ACCESS_SECRET   = "nbEKnc3E1p5a8wK74PkkmJjAsHVxfJIi"
	DEFAULT_TOKEN_HOUR_LIFESPAN = "1"
)

func GetBearerToken(authorizationHeader string) string {
	if len(strings.Split(authorizationHeader, " ")) == 2 {
		return strings.Split(authorizationHeader, " ")[1]
	}
	return ""
}

func Generate(claim *jwt_schema.JWTAccessClaims) (string, error) {
	secretStr := utils.GetENV("JWT_ACCESS_SECRET", DEFAULT_JWT_ACCESS_SECRET)
	secret := []byte(secretStr)
	lifeSpan := utils.GetENV("TOKEN_HOUR_LIFESPAN", DEFAULT_TOKEN_HOUR_LIFESPAN)
	lifeSpanInt, err := strconv.Atoi(lifeSpan)
	if err != nil {
		log.Println("Conversion error")
		lifeSpanInt = 1
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = claim.ID
	claims["username"] = claim.Username
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(lifeSpanInt)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseBearerToken(tokenString string) (*jwt.Token, error) {
	secretStr := utils.GetENV("JWT_ACCESS_SECRET", DEFAULT_JWT_ACCESS_SECRET)
	secret := []byte(secretStr)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractID(tokenString string) (id string, err error) {
	token, err := ParseBearerToken(tokenString)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		id = fmt.Sprintf("%s", claims["id"])
		return
	}
	return "", nil
}

func ExtractUserName(tokenString string) (username string, err error) {
	token, err := ParseBearerToken(tokenString)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		username = fmt.Sprintf("%s", claims["username"])
		return
	}
	return "", nil
}

func ExtractJWTClaims(tokenString string) (*jwt_schema.JWTAccessClaimsJSON, error) {
	token, err := ParseBearerToken(tokenString)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		id := fmt.Sprintf("%s", claims["id"])
		username := fmt.Sprintf("%s", claims["username"])
		return &jwt_schema.JWTAccessClaimsJSON{
			ID:       id,
			Username: username,
		}, nil
	}
	return &jwt_schema.JWTAccessClaimsJSON{
		ID:       "",
		Username: "",
	}, nil
}

func Validate(tokenString string) error {
	_, err := ParseBearerToken(tokenString)
	return err
}

func ExtractFromGin(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	authorizationHeader := c.Request.Header.Get("Authorization")
	bearerToken := GetBearerToken(authorizationHeader)
	return bearerToken
}

func ExtractIDFromGin(c *gin.Context) (id string, err error) {
	tokenString := ExtractFromGin(c)
	id, err = ExtractID(tokenString)
	return
}

func ExtractUserNameFromGin(c *gin.Context) (username string, err error) {
	tokenString := ExtractFromGin(c)
	username, err = ExtractID(tokenString)
	return
}

func ExtractJWTClaimsFromGin(c *gin.Context) (claims *jwt_schema.JWTAccessClaimsJSON, err error) {
	tokenString := ExtractFromGin(c)
	claims, err = ExtractJWTClaims(tokenString)
	return
}

func ValidateFromGin(c *gin.Context) error {
	tokenString := ExtractFromGin(c)
	err := Validate(tokenString)
	if err != nil {
		return err
	}
	return nil
}
