package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type fundamentalsCacheEntry struct {
	Data      map[string]FinanceData
	FetchedAt time.Time
}

var (
	fundamentalsCache *fundamentalsCacheEntry
	fundamentalsMu    sync.RWMutex
	fundamentalsTTL   = 2 * time.Hour
)

type FinanceData struct {
	PE  float64 `json:"pe"`
	EPS float64 `json:"eps"`
}

func GetSheetsService() *sheets.Service {
	jsonStr := os.Getenv("GOOGLE_SERVICE_ACCOUNT_JSON")
	if jsonStr == "" {
		log.Fatalf("GOOGLE_SERVICE_ACCOUNT_JSON not set")
	}

	conf, err := google.JWTConfigFromJSON(
		[]byte(jsonStr),
		sheets.SpreadsheetsScope,
	)
	if err != nil {
		log.Fatalf("JWT config failed: %v", err)
	}

	srv, err := sheets.NewService(
		context.Background(),
		option.WithHTTPClient(conf.Client(context.Background())),
	)
	if err != nil {
		log.Fatalf("Sheets service failed: %v", err)
	}

	return srv
}

func GetSheetValues(srv *sheets.Service, spreadsheetID, readRange string) [][]interface{} {
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil
	}
	return resp.Values
}

func WriteSheetFormula(
	srv *sheets.Service,
	spreadsheetID string,
	writeRange string,
	formula string,
) error {

	_, err := srv.Spreadsheets.Values.Update(
		spreadsheetID,
		writeRange,
		&sheets.ValueRange{
			Values: [][]interface{}{{formula}},
		},
	).ValueInputOption("USER_ENTERED").Do()

	return err
}

func ReadCell(
	srv *sheets.Service,
	spreadsheetID string,
	cell string,
) (string, error) {

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, cell).Do()
	if err != nil {
		return "", err
	}

	if len(resp.Values) == 0 || len(resp.Values[0]) == 0 {
		return "", fmt.Errorf("empty")
	}

	return fmt.Sprint(resp.Values[0][0]), nil
}

func ClearSheetRange(srv *sheets.Service, spreadsheetID string, rng string) {
	_, _ = srv.Spreadsheets.Values.Clear(
		spreadsheetID,
		rng,
		&sheets.ClearValuesRequest{},
	).Do()
}

func BuildFinanceFormula(symbols []string) string {
	formula := "=\"{\" & "
	for i, s := range symbols {
		formula += fmt.Sprintf(
			`"""%s"":{""eps"": " & IFERROR(GOOGLEFINANCE("%s","EPS"),"null") & ",""pe"": " & IFERROR(GOOGLEFINANCE("%s","PE"),"null") & "}"`,
			s, s, s,
		)
		if i < len(symbols)-1 {
			formula += " & \",\" & "
		}
	}
	formula += " & \"}\""
	return formula
}

func splitSymbols(symbols []string, size int) [][]string {

	var chunks [][]string

	for size < len(symbols) {
		symbols, chunks = symbols[size:], append(chunks, symbols[0:size])
	}

	chunks = append(chunks, symbols)

	return chunks
}

func writeGroupAndRead(
	srv *sheets.Service,
	spreadsheetID string,
	symbols []string,
	row int,
) (map[string]FinanceData, error) {

	targetCell := fmt.Sprintf("Sheet1!A%d", row)

	// write formula
	err := WriteSheetFormula(
		srv,
		spreadsheetID,
		targetCell,
		BuildFinanceFormula(symbols),
	)
	if err != nil {
		return nil, err
	}

	fmt.Println("check 1")

	raw, err := ReadCell(srv, spreadsheetID, targetCell)
	if err != nil {
		return nil, err
	}

	// remove cell formula later
	// ClearSheetRange(srv, spreadsheetID, targetCell)

	var out map[string]FinanceData

	err = json.Unmarshal([]byte(raw), &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func GetGoogleFundamentalsBatch(
	srv *sheets.Service,
	spreadsheetID string,
	symbols []string,
) (map[string]FinanceData, error) {

	// 1️⃣ FAST PATH — cache hit
	fundamentalsMu.RLock()
	if fundamentalsCache != nil &&
		time.Since(fundamentalsCache.FetchedAt) < fundamentalsTTL {

		data := fundamentalsCache.Data
		fundamentalsMu.RUnlock()
		return data, nil
	}
	fundamentalsMu.RUnlock()

	// 2️⃣ CACHE MISS → fetch from Google Sheets
	chunks := splitSymbols(symbols, 10)
	result := make(map[string]FinanceData)

	row := 1
	for _, group := range chunks {

		data, err := writeGroupAndRead(
			srv,
			spreadsheetID,
			group,
			row,
		)
		if err != nil {
			return nil, err
		}

		for sym, v := range data {
			result[sym] = v
		}

		row++
	}

	// 3️⃣ STORE IN CACHE
	fundamentalsMu.Lock()
	fundamentalsCache = &fundamentalsCacheEntry{
		Data:      result,
		FetchedAt: time.Now(),
	}
	fundamentalsMu.Unlock()

	return result, nil
}
