package util


//Constance for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	COR = "COR"
)

//Func help to valid the account currency
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, COR:
		return true
	}
	return false
}