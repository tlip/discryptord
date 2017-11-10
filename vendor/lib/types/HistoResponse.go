package types

// HistoResponse :: Response Object from CryptoCompare HistoTicker API
type HistoResponse struct {
	Response          string
	Aggregated        bool
	Data              []HistoTicker
	TimeTo            int64
	TimeFrom          int64
	FirstValueInArray bool
	ConversionType    map[string]string
}
