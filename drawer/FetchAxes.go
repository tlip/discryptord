package drawer

import (
	"sync"
	"time"

	"github.com/flamingyawn/discryptord/types"
)

// FetchAxes :: Fetch axes data
func FetchAxes(data []types.HistoMinute) types.AxesMap {
	var wg, wg2 sync.WaitGroup

	axes := types.AxesMap{
		X:      make([]time.Time, len(data)),
		Y:      make([]float64, len(data)),
		Vol:    make([]float64, len(data)),
		Ymin:   10000000000.0,
		Ymax:   0.0,
		Volmin: 10000000000.0,
		Volmax: 0.0,
	}

	wg.Add(len(data))
	for i, minute := range data {
		go func(i int, minute types.HistoMinute) {
			defer wg.Done()

			axes.X[i] = time.Unix(minute.Time, 0)
			axes.Y[i] = minute.Close
			axes.Vol[i] = minute.Volumeto

			if axes.Y[i] < axes.Ymin {
				axes.Ymin = axes.Y[i]
			} else if axes.Y[i] > axes.Ymax {
				axes.Ymax = axes.Y[i]
			}

			if axes.Vol[i] < axes.Volmin {
				axes.Volmin = axes.Vol[i]
			} else if axes.Vol[i] > axes.Volmax {
				axes.Volmax = axes.Vol[i]
			}

		}(i, minute)
	}
	wg.Wait()

	wg2.Add(len(axes.Vol))
	for i, v := range axes.Vol {
		go func(i int, v float64) {
			defer wg2.Done()
			volRange := (axes.Volmax - axes.Volmin)
			yRange := (axes.Ymax - axes.Ymin)
			axes.Vol[i] = (((v - axes.Volmin) / volRange) * yRange) + axes.Ymin
		}(i, v)
	}
	wg2.Wait()

	return axes
}
