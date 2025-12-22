package api

import (
	"backend/loader"
	"backend/services"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type PortfolioResponse struct {
	Particulars   string  `json:"particulars"`
	Symbol        string  `json:"symbol"`
	PurchasePrice float64 `json:"purchasePrice"`
	Quantity      float64 `json:"quantity"`
	Investment    float64 `json:"investment"`
	PortfolioPct  float64 `json:"portfolioPct"`
	CMP           float64 `json:"cmp"`
	PresentValue  float64 `json:"presentValue"`
	GainLoss      float64 `json:"gainLoss"`
	PE            float64 `json:"pe"`
	EPS           float64 `json:"eps"`
}

func PortfolioHandler(path string) gin.HandlerFunc {

	portfolio, _ := loader.LoadCSV(path)

	srv := services.GetSheetsService()
	spreadsheetID := os.Getenv("SPREADSHEET_ID")

	return func(c *gin.Context) {

		// 1️⃣ extract symbols from CSV
		var symbols []string
		for _, r := range portfolio {
			symbols = append(symbols, r.Symbol)
		}

		// 2️⃣ Google call ONCE, JSON returned for all stocks
		data, err := services.GetGoogleFundamentalsBatch(
			srv,
			spreadsheetID,
			symbols,
		)

		if err != nil {
			fmt.Println("GF ERROR =", err)
		}

		var response []PortfolioResponse

		// 3️⃣ Build portfolio output
		for _, r := range portfolio {

			cmp := services.GetCMP(r.Symbol)

			presentValue := cmp * r.Quantity
			gainLoss := presentValue - r.Investment

			// finance lookup
			fin := data[r.Symbol]

			response = append(response, PortfolioResponse{
				Particulars:   r.Particulars,
				Symbol:        r.Symbol,
				PurchasePrice: r.PurchasePrice,
				Quantity:      r.Quantity,
				Investment:    r.Investment,
				PortfolioPct:  r.PortfolioRatio,
				CMP:           cmp,
				PresentValue:  presentValue,
				GainLoss:      gainLoss,
				PE:            fin.PE,
				EPS:           fin.EPS,
			})
		}

		c.JSON(200, response)
	}
}
