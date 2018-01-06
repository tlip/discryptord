package api

import "strings"

// BuildHistoDayURL :: Construct URL for API call
func BuildHisto6mURL(coin, base string) string {
	scheme := "https://"
	host := "min-api.cryptocompare.com"
	path := "/data/histoday"
	link := scheme + host + path

	fsym := "?fsym=" + strings.ToUpper(coin)
	tsym := "&tsym=" + strings.ToUpper(base)
	limit := "&limit=183"
	aggregate := "&aggregate=1"
	query := fsym + tsym + limit + aggregate

	return link + query
}
