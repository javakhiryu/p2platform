package util

const (
	USD = "USD"
	EUR = "EUR"
	UZS = "UZS"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case EUR, USD, UZS:
		return true
	}
	return false
}
