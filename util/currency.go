package util

// List of supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	CNY = "CNY"
)

// IsSupportCurrency checks if the given currency is supported
func IsSupportCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, CNY:
		return true
	}
	return false
}
