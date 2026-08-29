package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"mortgage-bonus-calculator/handlers"
	"mortgage-bonus-calculator/webassets"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	staticFS := mustSubFS("static")
	r.StaticFS("/static", http.FS(staticFS))
	r.SetHTMLTemplate(template.Must(template.New("index.html").Funcs(template.FuncMap{
		"formatNumber": formatNumber,
		"printf":       formatPrintf,
		"mul":          func(left, right float64) float64 { return left * right },
	}).ParseFS(
		webassets.FS,
		"templates/*.html",
		"templates/partials/*.html",
	)))

	r.GET("/", handlers.Home)
	r.HEAD("/", handlers.Home)
	r.GET("/politica-proteccion-datos", handlers.DataProtectionPolicy)
	r.HEAD("/politica-proteccion-datos", handlers.DataProtectionPolicy)
	r.POST("/calcular", handlers.Calculate)
	r.POST("/api/calcular", handlers.CalculateJSON)
	r.POST("/contacto", handlers.Contact)

	return r
}

func formatPrintf(format string, values ...interface{}) string {
	if format == "%.2f" && len(values) == 1 {
		switch value := values[0].(type) {
		case float64:
			return formatNumber(value)
		case float32:
			return formatNumber(float64(value))
		}
	}
	return fmt.Sprintf(format, values...)
}

func formatNumber(value float64) string {
	formatted := fmt.Sprintf("%.2f", value)
	parts := strings.SplitN(formatted, ".", 2)
	integer := parts[0]
	for index := len(integer) - 3; index > 0; index -= 3 {
		integer = integer[:index] + "." + integer[index:]
	}
	if len(parts) == 1 {
		return integer + ",00"
	}
	return integer + "," + parts[1]
}

func mustSubFS(name string) fs.FS {
	subFS, err := fs.Sub(webassets.FS, name)
	if err != nil {
		panic("failed to load embedded assets: " + err.Error())
	}
	return subFS
}
