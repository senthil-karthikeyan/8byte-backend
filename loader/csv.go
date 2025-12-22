package loader

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

type PortfolioRow struct {
	Particulars    string
	PurchasePrice  float64
	Quantity       float64
	Investment     float64
	PortfolioRatio float64
	Symbol         string
}

func LoadCSV(path string) ([]PortfolioRow, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = -1

	// skip header
	_, _ = r.Read()

	var list []PortfolioRow

	for {
		row, err := r.Read()
		if err != nil {
			break
		}

		particulars := strings.TrimSpace(row[0])

		price, _ := strconv.ParseFloat(row[1], 64)
		qty, _ := strconv.ParseFloat(row[2], 64)

		symbol := strings.TrimSpace(row[5])

		list = append(list, PortfolioRow{
			Particulars:    particulars,
			PurchasePrice:  price,
			Quantity:       qty,
			Investment:     price * qty,
			PortfolioRatio: (price * qty) / 100,
			Symbol:         symbol,
		})
	}

	return list, nil
}
