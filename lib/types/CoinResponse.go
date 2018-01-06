package types

// CoinResponse :: Response Object from CryptoCompare HistoTicker API
type CoinResponse struct {
	Response          string
	Aggregated        bool
	Data              []HistoTicker
	TimeTo            int64
	TimeFrom          int64
	FirstValueInArray bool
	ConversionType    map[string]string
}
