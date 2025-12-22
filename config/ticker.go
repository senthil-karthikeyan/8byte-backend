package config

var NSEBSEMap = map[string]string{
	"LTIM":       "LTIM.NS",
	"HDFCBANK":   "HDFCBANK.NS",
	"DMART":      "DMART.NS",
	"BAJFINANCE": "BAJFINANCE.NS",
	"ASTRAL":     "ASTRAL.NS",
	"AFFLE":      "AFFLE.NS",
	"544252":     "BAJAJHFL.NS",
	"544107":     "BLSE.NS",
	"544028":     "TATATECH.NS",
	"543517":     "HARIOMPIPE.NS",
	"543318":     "CLEAN.NS",
	"542851":     "GENSOL.NS",
	"542652":     "POLYCAB.NS",
	"542651":     "KPITTECH.NS",
	"542323":     "KPIGREEN.NS",
	"541557":     "FINEORG.NS",
	"540719":     "SBILIFE.NS",
	"533282":     "GRAVITA.NS",
	"532790":     "TANLA.NS",
	"532667":     "SUZLON.NS",
	"532540":     "TATACONSUM.NS",
	"532174":     "ICICIBANK.NS",
	"511577":     "SAVFI.BO",
	"506401":     "DEEPAKNTR.NS",
	"500400":     "TATAPOWER.NS",
	"500331":     "PIDILITIND.NS",
}

func GetNSEBSE(key string) string {
	if val, ok := NSEBSEMap[key]; ok {
		return val
	}
	return ""
}
