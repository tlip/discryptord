package volume

import "strings"

// VolumeFor :: Build API call for pricemultifull
func VolumeFor(coin, base string) string {
	scheme := "https://"
	host := "min-api.cryptocompare.com"
	path := "/data/pricemultifull"
	link := scheme + host + path

	fsym := "?fsyms=" + strings.ToUpper(coin)
	tsym := "&tsyms=" + strings.ToUpper(base)
	query := fsym + tsym

	return link + query
}
