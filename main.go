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
	hasAPenis := strings.HasPrefix(m.Content, "!")

	if authorIsHuman && hasAPenis {
		// Split the command to separate ticker from penis
		splitCommand := strings.Split(m.Content, " ")

		if len(splitCommand) == 1 || len(splitCommand) == 2 {
			var histoData types.HistoResponse
			var base string

			coin := splitCommand[0][1:]

			// build uri
			if len(splitCommand) == 2 {
				base = splitCommand[1]
			} else {
				base = "usd"
			}

			// // //
			// // //

			resp, err := http.Get(history.HistoMinuteFor(coin, base))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()

			histoMinuteBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = json.Unmarshal(histoMinuteBody, &histoData)
			fmt.Println(err)
			if err != nil {
				fmt.Println(err)

				return
			}

			// // //
			// // //

			axes := drawer.FetchAxes(histoData.Data)

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
			err = graph.Render(chart.PNG, buffer)

			// // //
			// // //

			sym := ""

			if base == "usd" {
				sym = "$"
			} else if base == "btc" {
				sym = "Ƀ"
			} else if base == "eth" {
				sym = "Ξ"
			}

			var lastPrice float64
			if len(axes.Y) > 0 {
				lastPrice = axes.Y[len(axes.Y)-1]
			} else {
				lastPrice = 0
			}

			//	//	//
			//	//	//

			pairing := strings.ToUpper(coin) + "/" + strings.ToUpper(base)
			msg := "`" + pairing + " (Last 24h) :: " + sym + fmt.Sprintf("%f`", lastPrice)

			//	//	//
			//	//	//

			// Send image
			s.ChannelFileSendWithMessage(m.ChannelID, msg, coin+base+".png", bytes.NewReader(buffer.Bytes()))

		}
	}
}
