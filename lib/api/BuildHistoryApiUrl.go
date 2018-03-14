package api

import "strings"

// BuildHistoryApiUrl :: Construct URL for API call
func BuildHistoryApiUrl(daterange, coin, base string) string {
	var path, limit, aggregate string
	scheme := "https://"
	host := "min-api.cryptocompare.com"

	switch daterange {
	case "3d":
		path, limit, aggregate = "/data/histohour", "&limit=72", "&aggregate=1"
	case "w":
		path, limit, aggregate = "/data/histohour", "&limit=168", "&aggregate=1"
	case "m":
		path, limit, aggregate = "/data/histohour", "&limit=120", "&aggregate=6"
	case "3m":
		path, limit, aggregate = "/data/histohour", "&limit=2000", "&aggregate=6"
	case "6m":
		path, limit, aggregate = "/data/histoday", "&limit=183", "&aggregate=1"
	case "y":
		path, limit, aggregate = "/data/histoday", "&limit=365", "&aggregate=1"
	case "24h":
		fallthrough
	default:
		path, limit, aggregate = "/data/histominute", "&limit=144", "&aggregate=10"
	}

	link := scheme + host + path

	fsym := "?fsym=" + strings.ToUpper(coin)
	tsym := "&tsym=" + strings.ToUpper(base)
	query := fsym + tsym + limit + aggregate

	return link + query
}
