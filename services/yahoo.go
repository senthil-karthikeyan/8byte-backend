package services

import (
	"backend/config"
	"sync"
	"time"

	yfa "github.com/oscarli916/yahoo-finance-api"
)

type priceCacheEntry struct {
	Price     float64
	FetchedAt time.Time
}

var (
	cmpCache = make(map[string]priceCacheEntry)
	cmpMu    sync.RWMutex
	inFlight = make(map[string]bool)
	cmpTTL   = 2 * time.Minute
)

func GetCMP(rawSymbol string) float64 {
	symbol := config.GetNSEBSE(rawSymbol)
	if symbol == "" {
		symbol = rawSymbol + ".NS"
	}

	// 1️⃣ Cache read
	cmpMu.RLock()
	entry, ok := cmpCache[symbol]
	if ok && time.Since(entry.FetchedAt) < cmpTTL {
		cmpMu.RUnlock()
		return entry.Price
	}
	cmpMu.RUnlock()

	// 2️⃣ Prevent duplicate fetches
	cmpMu.Lock()
	if inFlight[symbol] {
		cmpMu.Unlock()
		time.Sleep(200 * time.Millisecond)
		return entry.Price
	}
	inFlight[symbol] = true
	cmpMu.Unlock()

	// 3️⃣ Fetch from Yahoo
	t := yfa.NewTicker(symbol)
	q, err := t.Quote()

	cmpMu.Lock()
	delete(inFlight, symbol)

	if err == nil {
		cmpCache[symbol] = priceCacheEntry{
			Price:     q.Close,
			FetchedAt: time.Now(),
		}
	}
	cmpMu.Unlock()

	if err != nil {
		return entry.Price // fallback
	}

	return q.Close
}
