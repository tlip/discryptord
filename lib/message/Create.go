package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"

	"github.com/flamingyawn/discryptord/lib/api"
	"github.com/flamingyawn/discryptord/lib/drawer"
	"github.com/flamingyawn/discryptord/lib/types"

	"github.com/bwmarrin/discordgo"
)

// Create :: called once for every message on any channel that the autenticated bot has access to.
func Create(s *discordgo.Session, m *discordgo.MessageCreate) {

	// // //
	// prevent overflow
	if len(m.Content) <= 2 {
		return
	}

	hasBeenInvoked := strings.HasPrefix(m.Content, "!") // check invocation
	authorIsHuman := (m.Author.ID != s.State.User.ID)   // check humanity of invoker
	if authorIsHuman && hasBeenInvoked {
		// // //
		// separate ticker from invocation
		// and set default RSI presence
		splitCommand := strings.Split(m.Content, " ")
		rsiEnabled, logEnabled := false, false
		candle := "minute"

		if len(splitCommand) > 1 {
			for strings.HasPrefix(splitCommand[len(splitCommand)-1], "-") || splitCommand[len(splitCommand)-1] == "" {
				flag := splitCommand[len(splitCommand)-1]
				splitCommand = splitCommand[:len(splitCommand)-1]

				if flag == "-3d" {
					candle = "3d"
				} else if flag == "-w" {
					candle = "hour"
				} else if flag == "-m" {
					candle = "day"
				} else if flag == "-3m" {
					candle = "3m"
				} else if flag == "-6m" {
					candle = "6m"
				} else if flag == "-y" {
					candle = "y"
				} else if flag == "-rsi" || flag == "-RSI" {
					rsiEnabled = true
				} else if flag == "-log" || flag == "-LOG" {
					logEnabled = true
				}

			}
		}

		// // //
		// prevent overflow
		if len(splitCommand) <= 2 {
			// // //
			// get tickers
			var base string
			coin := strings.ToUpper(splitCommand[0][1:])
			if len(splitCommand) == 2 {
				base = strings.ToUpper(splitCommand[1])
			} else {
				base = "USD"
			}

			// // //
			// fetch data
			var histoData types.HistoResponse
			var apiURL, timerange string
			if candle == "3d" {
				apiURL = api.BuildThreeDayURL(coin, base)
				timerange = "3D"
			} else if candle == "hour" {
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
			} else if candle == "6m" {
				apiURL = api.BuildHisto6mURL(coin, base)
				timerange = "6M"
			} else if candle == "y" {
				apiURL = api.BuildYearURL(coin, base)
				timerange = "1Y"
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

			// // //
			// format axes
			axes := drawer.ParsePriceData(histoData.Data, logEnabled)

			// // //
			// draw chart
			buffer, err := drawer.DrawChart(axes, rsiEnabled, logEnabled)
			if err != nil {
				fmt.Println(err)
				return
			}

			// // //
			// format closing price
			var sym string
			color := 0x5dff9f
			lastPrice, firstPrice := 0.0, 0.0
			hi, lo := axes.Ymax, axes.Ymin
			logText := ""
			changeSign := "-"

			switch base {
			case "usd":
			case "usdt":
			case "USDT":
			case "USD":
				sym = "$"
			case "btc":
			case "BTC":
				sym = "Ƀ"
			case "eth":
			case "ETH":
				sym = "Ξ"
			default:
				sym = ""
			}

			if len(axes.Y) > 0 {
				lastPrice = axes.Y[len(axes.Y)-1]
				firstPrice = axes.Y[0]
			}

			if logEnabled && len(axes.Y) > 0 {
				logText = "Logarithmic "

				if lastPrice > 0 {
					lastPrice = math.Pow(10, lastPrice)
				}
				if firstPrice > 0 {
					firstPrice = math.Pow(10, firstPrice)
				}
				if hi > 0 {
					hi = math.Pow(10, hi)
				}
				if lo > 0 {
					lo = math.Pow(10, lo)
				}

			}

			if lastPrice >= firstPrice {
				changeSign = "+"
			}

			//	//
			//	build message
			delta := lastPrice - firstPrice
			pairing := fmt.Sprintf("%s/%s %sPrice Chart (%s)", coin, base, logText, timerange)
			deltaPct := fmt.Sprintf("%.2f%%", delta/firstPrice*100)

			if delta < 0 {
				color = 0xE94335
			}

			embed := NewEmbed().
				SetAuthor(pairing, "https://cdn.discordapp.com/app-icons/359564584564293632/21fb4ad276ed1ddc3318ce0b1a663395.png").
				AddField("Last", fmt.Sprintf("%s%f", sym, lastPrice)).
				AddField("Hi", fmt.Sprintf("%s%f", sym, hi)).
				AddField("Lo", fmt.Sprintf("%s%f", sym, lo)).
				AddField("∆", fmt.Sprintf("`%s%s%f (%s%s)`", changeSign, sym, math.Abs(delta), changeSign, deltaPct)).
				InlineAllFields().
				SetColor(color).
				MessageEmbed

			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			s.ChannelFileSendWithMessage(m.ChannelID, "", coin+base+".png", bytes.NewReader(buffer.Bytes()))

		}
	}
}
