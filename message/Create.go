package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/flamingyawn/discryptord/drawer"
	"github.com/flamingyawn/discryptord/history"
	"github.com/flamingyawn/discryptord/types"
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
			if err != nil {
				fmt.Println(err)
				return
			}

			// format axes
			// // //
			axes := drawer.FetchAxes(histoData.Data)

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
			msg := fmt.Sprintf("`%s :: %s%f                            24h`", pairing, sym, lastPrice)

			// send response
			//	//
			s.ChannelFileSendWithMessage(m.ChannelID, msg, coin+base+".png", bytes.NewReader(buffer.Bytes()))

		}
	}
}
