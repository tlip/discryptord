package api

import "strings"

// BuildHistoHourURL :: Construct URL for API call
func BuildHistoHourURL(coin, base string) string {
	scheme := "https://"
	host := "min-api.cryptocompare.com"
	path := "/data/histohour"
	link := scheme + host + path

	fsym := "?fsym=" + strings.ToUpper(coin)
	tsym := "&tsym=" + strings.ToUpper(base)
	limit := "&limit=168"
	aggregate := "&aggregate=1"
	query := fsym + tsym + limit + aggregate

	return link + query
}
