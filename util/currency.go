package util

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	NGN = "NGN"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, NGN:
		return true
	}
	return false
}
