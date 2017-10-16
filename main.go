package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/common-nighthawk/go-figure"
	"github.com/flamingyawn/discryptord/drawer"
	"github.com/flamingyawn/discryptord/history"
	"github.com/flamingyawn/discryptord/types"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Make sure cryptograph isn't responding to any seedy bots or females
	if len(m.Content) <= 2 {
		return
	}
	authorIsHuman := (m.Author.ID != s.State.User.ID)
	hasAPenis := (m.Content[:2] == "# ")

	if authorIsHuman && hasAPenis {
		// Split the command to separate ticker from penis
		splitCommand := strings.Split(m.Content, " ")

		if len(splitCommand) == 2 || len(splitCommand) == 3 {
			var histoData types.HistoResponse
			// var priceMultiData types.PriceMultiFull
			var base string
			// var priceMultiData interface{}

			coin := splitCommand[1]
			// ymin, ymax, volmin, volmax := 100000000000.0, 0.0, 100000000000.0, 0.0

			// build uri
			if len(splitCommand) == 3 {
				base = splitCommand[2]
			} else {
				base = "usd"
			}

			// // //
			// // //

			// resp, err := http.Get(volume.VolumeFor(coin, base))
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }
			// defer resp.Body.Close()

			// //
			// body, err := ioutil.ReadAll(resp.Body)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			// //
			// err = json.Unmarshal(body, &priceMultiData)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			// // //
			// // //

			resp, err := http.Get(history.HistoMinuteFor(coin, base))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()

			//
			histoMinuteBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			//
			err = json.Unmarshal(histoMinuteBody, &histoData)
			if err != nil {
				fmt.Println(err)
				return
			}

			// // //
			// // //

			axes := drawer.FetchAxes(histoData.Data)

			// for i, v := range axes.Vol {
			// 	// go func(i int, v float64) {
			// 	// defer wg.Done()
			// 	volRange := (axes.Volmax - axes.Volmin)
			// 	yRange := (axes.Ymax - axes.Ymin)
			// 	axes.Vol[i] = (((v - axes.Volmin) / volRange) * yRange) + axes.Ymin
			// 	// }(i, v)
			// }

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

			// bbSeries := &chart.BollingerBandsSeries{
			// 	Name: "SPY - Bol. Bands",
			// 	Style: chart.Style{
			// 		Show:        true,
			// 		StrokeColor: drawing.ColorFromHex("ffffff").WithAlpha(30),
			// 		FillColor:   drawing.ColorFromHex("ffffff").WithAlpha(1),
			// 	},
			// 	InnerSeries: priceSeries,
			// }

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
			err = graph.Render(chart.PNG, buffer)

			// img, _, _ := image.Decode(bytes.NewReader(buffer.Bytes()))
			// out, err3 := os.Create("./img/graph.png")
			// var out io.Writer
			// if err3 != nil {
			// 	fmt.Println(err3)
			// }
			// err = png.Encode(out, img)

			// Read image
			// finalImg, err4 := os.Open("./img/graph.png")
			// defer finalImg.Close()
			// if err4 != nil {
			// 	fmt.Println(err4)
			// }

			// Send image
			// var sym string
			var price string

			// if base == "usd" {
			// 	sym = "$"
			// } else if base == "btc" {
			// 	sym = "Ƀ"
			// } else if base == "eth" {
			// 	sym = "Ξ"
			// } else {
			// 	sym = ""
			// }

			ticker := fmt.Sprintf("%s/%s (24h)", strings.ToUpper(coin), strings.ToUpper(base))
			// pRange := fmt.Sprintf("%s%f - %s%f", sym, axes.Ymin, sym, axes.Ymax)
			// head := ticker + fmt.Sprintf("(24h :: %s)\n", pRange)
			closePrice := fmt.Sprintf("%f", axes.Y[len(axes.Y)-1])

			if len(closePrice) > 6 {
				closePrice = closePrice[:6]
			}

			asciiPrice := figure.NewFigure(closePrice, "banner3", true).Slicify()

			for _, row := range asciiPrice {
				price = price + row + "\n"
			}

			fmt.Println(price)
			msg := "```go\n" + ticker + "\n\n" + price + "```"

			s.ChannelFileSendWithMessage(m.ChannelID, msg, splitCommand[1]+base+".png", bytes.NewReader(buffer.Bytes()))

		}
	}
}
