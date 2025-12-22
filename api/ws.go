package api

import (
	"backend/loader"
	"backend/services"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type LiveUpdate struct {
	Symbol       string  `json:"symbol"`
	CMP          float64 `json:"cmp"`
	PresentValue float64 `json:"presentValue"`
	GainLoss     float64 `json:"gainLoss"`
}

func PortfolioLiveWS(path string) gin.HandlerFunc {

	portfolio, _ := loader.LoadCSV(path)

	return func(c *gin.Context) {

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {

			<-ticker.C

			var updates []LiveUpdate

			for _, r := range portfolio {

				cmp := services.GetCMP(r.Symbol)

				presentValue := cmp * r.Quantity
				gainLoss := presentValue - r.Investment

				updates = append(updates, LiveUpdate{
					Symbol:       r.Symbol,
					CMP:          cmp,
					PresentValue: presentValue,
					GainLoss:     gainLoss,
				})
			}

			payload, _ := json.Marshal(updates)
			err := conn.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				return
			}
		}
	}
}
