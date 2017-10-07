package history

import "strings"

func HistoMinuteFor(coin, base string) string {
	scheme := "https://"
	host := "min-api.cryptocompare.com"
	path := "/data/histominute"
	link := scheme + host + path

	fsym := "?fsym=" + strings.ToUpper(coin)
	tsym := "&tsym=" + strings.ToUpper(base)
	limit := "&limit=180"
	aggregate := "&aggregate=2"
	query := fsym + tsym + limit + aggregate

	return link + query
}
