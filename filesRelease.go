//go:build release

package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

//go:embed templates/* static/*
var content embed.FS

func SetupTemplatesAndStatic(r *gin.Engine) {
	r.SetHTMLTemplate(template.Must(template.ParseFS(content, "templates/*")))

	r.GET("/static/*file", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(content))
	})
}
