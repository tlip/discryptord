package history

import "strings"

func HistoMinuteFor(ticker string) string {
	scheme := "https://"
	host := "min-api.cryptocompare.com"
	path := "/data/histominute"
	link := scheme + host + path

	fsym := "?fsym=" + strings.ToUpper(ticker)
	tsym := "&tsym=USD"
	limit := "&limit=60"
	aggregate := "&aggregate=6"
	query := fsym + tsym + limit + aggregate

	return link + query
}
