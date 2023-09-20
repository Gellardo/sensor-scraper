//go:build !release

package main

import (
	"github.com/gin-gonic/gin"
)

func SetupTemplatesAndStatic(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")

	r.Static("/static", "./static")
}
