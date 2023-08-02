package utils

import (
	"encoding/json"
	"log"

	utils "github.com/emnopal/odoo-golang-restapi/app/utils/env"
	"github.com/mervick/aes-everywhere/go/aes256"
)

func DecryptAes256(data string, v interface{}) {
	key := utils.GetENV("LOGIN_KEY", "nbEKnc3E1p5a8wK74PkkmJjAsHVxfJIi")
	decrypted := aes256.Decrypt(data, key)
	jsonData := []byte(decrypted)
	err1 := json.Unmarshal(jsonData, &v)
	if err1 != nil {
		log.Println(err1)
	}
}
