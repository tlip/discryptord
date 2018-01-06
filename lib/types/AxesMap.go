package types

import "time"

// AxesMap :: Map of axes and their metadata
type AxesMap struct {
	X        []time.Time
	Y        []float64
	Vol      []float64
	VolFixed []float64
	Xmax     float64
	Ymin     float64
	Ymax     float64
	Volmax   float64
	Volmin   float64
}
