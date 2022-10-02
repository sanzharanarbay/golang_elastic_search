package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	v1 "github.com/sanzharanarbay/golang-elastic-search/application/routes/api/v1"
	"os"
)

func main() {
	port := os.Getenv("APP_PORT")

	prefix := os.Getenv("ROUTE_PREFIX")
	fmt.Println("Server started at " + port + "...")

	router := gin.New()
	v1.ApiRoutes(prefix, router)
	router.Run(":" + port)
}
