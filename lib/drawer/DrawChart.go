package drawer

import (
	"bytes"

	"github.com/flamingyawn/discryptord/lib/types"
	"github.com/vertis/gota/ta"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

// DrawChart :: Draws chart and writes to buffer
func DrawChart(axes types.AxesMap, rsiEnabled bool) (*bytes.Buffer, error) {

	var rsiSeries, rsiTopLine, rsiBottomLine chart.TimeSeries
	var rsi, rsi70Line, rsi30Line []float64
	VOL, Y := axes.Vol, axes.Y

	if rsiEnabled == true {
		rsi = ta.Rsi(axes.Y, 14)
		Yrange := (axes.Ymax*1.005 - axes.Ymin*0.995)

		for len(rsi) < len(axes.X) {
			rsi = append([]float64{30}, rsi...)
		}

		for i, n := range rsi {
			rsi[i] = 0.2*(n/100)*Yrange + (axes.Ymin * 0.995) - (.2 * .2 * Yrange)
			rsi30Line = append(rsi30Line, (0.2*0.3*Yrange)+(axes.Ymin*0.995)-(.2*.2*Yrange))
			rsi70Line = append(rsi70Line, (0.2*0.7*Yrange)+(axes.Ymin*0.995)-(.2*.2*Yrange))
		}

		for i, n := range Y {
			Y[i] = 0.8*(n-axes.Ymin*0.995) + (axes.Ymin * 0.995) + (.2 * Yrange)
		}

		for i, n := range VOL {
			VOL[i] = 0.8*(n-axes.Ymin*0.995) + (axes.Ymin * 0.995) + (.2 * Yrange)
		}
	}

	priceSeries := chart.TimeSeries{
		Name: "PRICE",
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorFromHex("4DE786"),
		},
		XValues: axes.X,
		YValues: Y,
	}

	volumeSeries := chart.TimeSeries{
		Name: "PRICE - VOL",
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorFromHex("00A1E7").WithAlpha(70),
		},
		XValues: axes.X,
		YValues: VOL,
	}

	smaSeries := chart.SMASeries{
		Name: "SMA",
		Style: chart.Style{
			Show:            true,
			StrokeColor:     drawing.ColorFromHex("AE73FF"),
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: priceSeries,
	}

	if rsiEnabled == true {
		rsiSeries = chart.TimeSeries{
			Name: "RSI",
			Style: chart.Style{
				Show:        true,
				StrokeColor: drawing.ColorFromHex("EE933E"),
			},
			XValues: axes.X,
			YValues: rsi,
		}

		rsiTopLine = chart.TimeSeries{
			Name: "RSI",
			Style: chart.Style{
				Show:        true,
				StrokeColor: drawing.ColorFromHex("FFFFFF").WithAlpha(40),
				FillColor:   drawing.ColorFromHex("34373C"),
			},
			XValues: axes.X,
			YValues: rsi70Line,
		}

		rsiBottomLine = chart.TimeSeries{
			Name: "RSI",
			Style: chart.Style{
				Show:        true,
				StrokeColor: drawing.ColorFromHex("FFFFFF").WithAlpha(40),
				FillColor:   drawing.ColorFromHex("36393E"),
			},
			XValues: axes.X,
			YValues: rsi30Line,
		}
	}

	chartSeries := []chart.Series{
		volumeSeries,
		priceSeries,
		smaSeries,
	}

	if rsiEnabled == true {
		chartSeries = append(chartSeries, rsiTopLine, rsiBottomLine, rsiSeries)
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
		Series: chartSeries,
	}

	buffer := bytes.NewBuffer([]byte{})

	// render and save chart
	err := graph.Render(chart.PNG, buffer)

	return buffer, err
}
