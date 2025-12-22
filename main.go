package main

import (
	"backend/api"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Example usage
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := gin.Default()

	r.GET("/portfolio", api.PortfolioHandler("portfolio.csv"))
	r.GET("/ws", api.PortfolioLiveWS("portfolio.csv"))

	r.Run(":8080")

}
