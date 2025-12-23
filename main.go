package main

import (
	"backend/api"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/portfolio", api.PortfolioHandler("portfolio.csv"))
	r.GET("/ws", api.PortfolioLiveWS("portfolio.csv"))

	r.Run(":8080")

}
