package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"lib/api"
	"lib/drawer"
	"lib/types"

	"github.com/bwmarrin/discordgo"
)

// Create :: called once for every message on any channel that the autenticated bot has access to.
func Create(s *discordgo.Session, m *discordgo.MessageCreate) {

	// prevent overflow
	// // //
	if len(m.Content) <= 2 {
		return
	}

	hasBeenInvoked := strings.HasPrefix(m.Content, "!") // check invocation
	authorIsHuman := (m.Author.ID != s.State.User.ID)   // check humanity of invoker

	if authorIsHuman && hasBeenInvoked {
		// separate ticker from invocation
		// // //
		splitCommand := strings.Split(m.Content, " ")
		candle := "minute"
		if len(splitCommand) > 1 {
			flag := splitCommand[len(splitCommand)-1]

			if flag == "-w" {
				candle = "hour"
				splitCommand = splitCommand[:len(splitCommand)-1]
			} else if flag == "-m" {
				candle = "day"
				splitCommand = splitCommand[:len(splitCommand)-1]
			} else if flag == "-3m" {
				candle = "3m"
				splitCommand = splitCommand[:len(splitCommand)-1]
			}

		}

		// prevent overflow
		// // //
		if len(splitCommand) <= 2 {
			// get tickers
			// // //
			var base string
			coin := strings.ToUpper(splitCommand[0][1:])
			if len(splitCommand) == 2 {
				base = strings.ToUpper(splitCommand[1])
			} else {
				base = "USD"
			}

			// fetch data
			// // //
			var histoData types.HistoResponse
			var apiURL, timerange string
			if candle == "hour" {
				apiURL = api.BuildHistoHourURL(coin, base)
				timerange = "7D"
			} else if candle == "minute" {
				apiURL = api.BuildHistoMinuteURL(coin, base)
				timerange = "24H"
			} else if candle == "day" {
				apiURL = api.BuildHistoDayURL(coin, base)
				timerange = "1M"
			} else if candle == "3m" {
				apiURL = api.BuildHistoHourAllURL(coin, base)
				timerange = "3M"
			}

			resp, err := http.Get(apiURL)
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
			if err != nil {
				fmt.Println(err)
				return
			}

			// format axes
			// // //
			axes := drawer.ParsePriceData(histoData.Data)

			// draw chart
			// // //
			buffer, err := drawer.DrawChart(axes)
			if err != nil {
				fmt.Println(err)
				return
			}

			// format closing price
			// // //
			var sym string
			var lastPrice float64

			switch base {
			case "usd":
				sym = "$"
			case "btc":
				sym = "Ƀ"
			case "eth":
				sym = "Ξ"
			default:
				sym = ""
			}

			if len(axes.Y) > 0 {
				lastPrice = axes.Y[len(axes.Y)-1]
			} else {
				lastPrice = 0.0
			}

			//	build message
			//	//
			pairing := fmt.Sprintf("%s/%s", coin, base)
			msg := fmt.Sprintf("`%s :: %s%f                            %s`", pairing, sym, lastPrice, timerange)

			// send response
			//	//
			s.ChannelFileSendWithMessage(m.ChannelID, msg, coin+base+".png", bytes.NewReader(buffer.Bytes()))

		}
	}
}
