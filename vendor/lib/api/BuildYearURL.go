package api

import "strings"

// BuildYearURL :: Construct URL for API call
func BuildYearURL(coin, base string) string {
	scheme := "https://"
	host := "min-api.cryptocompare.com"
	path := "/data/histoday"
	link := scheme + host + path

	fsym := "?fsym=" + strings.ToUpper(coin)
	tsym := "&tsym=" + strings.ToUpper(base)
	limit := "&limit=365"
	aggregate := "&aggregate=1"
	query := fsym + tsym + limit + aggregate

	return link + query
}
