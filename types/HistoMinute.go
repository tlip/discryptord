package types

// HistoMinute :: Individual record from HistoResponse
type HistoMinute struct {
	Time       int64
	Close      float64
	Low        float64
	Open       float64
	Volumefrom float64
	Volumeto   float64
}
