package drawer

import (
	"bytes"

	"lib/types"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

// DrawChart :: Draws chart and writes to buffer
func DrawChart(axes types.AxesMap) (*bytes.Buffer, error) {
	priceSeries := chart.TimeSeries{
		Name: "SPY",
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorFromHex("4DE786"),
		},
		XValues: axes.X,
		YValues: axes.Y,
	}

	volumeSeries := chart.TimeSeries{
		Name: "SPY - VOL",
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorFromHex("00A1E7").WithAlpha(70),
		},
		XValues: axes.X,
		YValues: axes.Vol,
	}

	smaSeries := chart.SMASeries{
		Name: "SPY - SMA",
		Style: chart.Style{
			Show:            true,
			StrokeColor:     drawing.ColorFromHex("AE73FF"),
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: priceSeries,
	}

	graph := chart.Chart{
		Canvas: chart.Style{
			FillColor: drawing.ColorFromHex("36393E"),
		},
		Background: chart.Style{
			FillColor: drawing.ColorFromHex("36393E"),
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				StrokeColor: drawing.ColorFromHex("ffffff"),
				Show:        false,
			},
			TickPosition: chart.TickPositionBetweenTicks,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: false},
			Range: &chart.ContinuousRange{
				Max: axes.Ymax * 1.005,
				Min: axes.Ymin * 0.995,
			},
		},
		Series: []chart.Series{
			volumeSeries,
			// bbSeries,
			priceSeries,
			smaSeries,
		},
	}

	buffer := bytes.NewBuffer([]byte{})

	// render and save chart
	err := graph.Render(chart.PNG, buffer)

	return buffer, err
}
