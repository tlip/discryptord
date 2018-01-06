package api

import "strings"

// BuildHistoMinuteURL :: Construct URL for API call
func BuildHistoMinuteURL(coin, base string) string {
	scheme := "https://"
	host := "min-api.cryptocompare.com"
	path := "/data/histominute"
	link := scheme + host + path

	fsym := "?fsym=" + strings.ToUpper(coin)
	tsym := "&tsym=" + strings.ToUpper(base)
	limit := "&limit=144"
	aggregate := "&aggregate=10"
	query := fsym + tsym + limit + aggregate

	return link + query
}
