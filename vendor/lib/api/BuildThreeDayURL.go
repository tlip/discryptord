package api

import "strings"

// BuildThreeDayURL :: Construct URL for API call
func BuildThreeDayURL(coin, base string) string {
	scheme := "https://"
	host := "min-api.cryptocompare.com"
	path := "/data/histohour"
	link := scheme + host + path

	fsym := "?fsym=" + strings.ToUpper(coin)
	tsym := "&tsym=" + strings.ToUpper(base)
	limit := "&limit=72"
	aggregate := "&aggregate=1"
	query := fsym + tsym + limit + aggregate

	return link + query
}
