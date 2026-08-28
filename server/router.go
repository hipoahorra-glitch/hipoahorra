package server

import (
	"html/template"
	"io/fs"
	"net/http"
	"mortgage-bonus-calculator/handlers"
	"mortgage-bonus-calculator/webassets"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	staticFS := mustSubFS("static")
	r.StaticFS("/static", http.FS(staticFS))
	r.SetHTMLTemplate(template.Must(template.ParseFS(
		webassets.FS,
		"templates/*.html",
		"templates/partials/*.html",
	)))

	r.GET("/", handlers.Home)
	r.HEAD("/", handlers.Home)
	r.POST("/calcular", handlers.Calculate)
	r.POST("/api/calcular", handlers.CalculateJSON)
	r.POST("/contacto", handlers.Contact)

	return r
}

func mustSubFS(name string) fs.FS {
	subFS, err := fs.Sub(webassets.FS, name)
	if err != nil {
		panic("failed to load embedded assets: " + err.Error())
	}
	return subFS
}
