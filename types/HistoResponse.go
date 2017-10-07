package types

// HistoResponse :: Response Object from CryptoCompare HistoMinute API
type HistoResponse struct {
	Response          string
	Aggregated        bool
	Data              []HistoMinute
	TimeTo            int64
	TimeFrom          int64
	FirstValueInArray bool
	ConversionType    map[string]string
}
