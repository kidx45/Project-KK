package utils

const (
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
)

func IsCurrencyValid(currency string) bool {
	switch currency {
	case CurrencyUSD, CurrencyEUR:
		return true
	}
	return false
}