package utils

import (
	utils "github.com/emnopal/odoo-golang-restapi/app/utils/Checking"
	"github.com/gin-gonic/gin"
)

func GetLangFromHeader(c *gin.Context) string {
	lang := c.Request.Header.Get("Lang")
	supportedLang := []string{
		"en_US", "ja_JP",
	}
	isLangSupported := utils.Contains(supportedLang, lang)
	if lang == "" || !isLangSupported {
		lang = "en_US"
	}
	return lang
}
