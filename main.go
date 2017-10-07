package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
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
	if len(m.Content) <= 6 {
		return
	}
	authorIsHuman := (m.Author.ID != s.State.User.ID)
	hasAPenis := (m.Content[:6] == "?price")

	if authorIsHuman && hasAPenis {
		// Split the command to separate ticker from penis
		splitCommand := strings.Split(m.Content, " ")

		if len(splitCommand) == 2 || len(splitCommand) == 3 {
			var histoData types.HistoResponse
			var coin string = splitCommand[1]
			var base string

			// build uri
			if len(splitCommand) == 3 {
				base = splitCommand[2]
			} else {
				base = "usd"
			}
			resp, err := http.Get(history.HistoMinuteFor(coin, base))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()

			//
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			//
			err = json.Unmarshal(body, &histoData)
			if err != nil {
				fmt.Println(err)
				return
			}

			////
			var xv []time.Time
			var yv []float64
			var ymin, ymax float64 = 1000000, 0

			for _, m := range histoData.Data {
				xv = append(xv, time.Unix(m.Time, 0))
				yv = append(yv, m.Close)
				if m.Close < ymin {
					ymin = m.Close
				}
				if m.Close > ymax {
					ymax = m.Close
				}
			}

			priceSeries := chart.TimeSeries{
				Name: "SPY",
				Style: chart.Style{
					Show:        true,
					StrokeColor: drawing.ColorFromHex("92FF92"),
				},
				XValues: xv,
				YValues: yv,
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

			bbSeries := &chart.BollingerBandsSeries{
				Name: "SPY - Bol. Bands",
				Style: chart.Style{
					Show:        true,
					StrokeColor: drawing.ColorFromHex("bcbcbc").WithAlpha(50),
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
						Max: ymax * 1.005,
						Min: ymin * 0.995,
					},
				},
				Series: []chart.Series{
					bbSeries,
					priceSeries,
					smaSeries,
				},
			}
			buffer := bytes.NewBuffer([]byte{})

			// render and save chart
			err = graph.Render(chart.PNG, buffer)
			img, _, _ := image.Decode(bytes.NewReader(buffer.Bytes()))
			out, err3 := os.Create("./img/graph.png")
			if err3 != nil {
				fmt.Println(err3)
			}
			err = png.Encode(out, img)

			// Read image
			finalImg, err4 := os.Open("./img/graph.png")
			defer finalImg.Close()
			if err4 != nil {
				fmt.Println(err4)
			}

			// Send image
			msg := "`" + strings.ToUpper(splitCommand[1]) + "/USD (Last 6h)`"
			s.ChannelFileSendWithMessage(m.ChannelID, msg, splitCommand[1]+"usd.png", finalImg)

		}
	}
}
