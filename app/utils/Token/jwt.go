package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	utils "github.com/emnopal/odoo-golang-restapi/app/utils/env"
	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_JWT_ACCESS_SECRET   = "nbEKnc3E1p5a8wK74PkkmJjAsHVxfJIi"
	DEFAULT_TOKEN_HOUR_LIFESPAN = "1"
)

func Generate(id string) (string, error) {

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
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(lifeSpanInt)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func Extract(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractID(c *gin.Context) (string, error) {

	secretStr := utils.GetENV("JWT_ACCESS_SECRET", DEFAULT_JWT_ACCESS_SECRET)
	secret := []byte(secretStr)

	tokenString := Extract(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		id := fmt.Sprintf("%.0f", claims["id"])
		return id, nil
	}

	return "", nil
}

func Validate(c *gin.Context) error {

	secretStr := utils.GetENV("JWT_ACCESS_SECRET", DEFAULT_JWT_ACCESS_SECRET)
	secret := []byte(secretStr)

	tokenString := Extract(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return err
	}

	return nil
}
