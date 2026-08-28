package main

import (
	"log"
	"os"

	"mortgage-bonus-calculator/server"

	"github.com/gin-gonic/gin"
)

func main() {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := server.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "111"
		port = "8080"
	}

	log.Printf("Server listening on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
