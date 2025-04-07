package util

const (
	USD = "USD" // Доллар США
	EUR = "EUR" // Евро
	UZS = "UZS" // Узбекский сум
	RUB = "RUB" // Российский рубль
	GBP = "GBP" // Британский фунт стерлингов
	JPY = "JPY" // Японская иена
	CHF = "CHF" // Швейцарский франк
	CNY = "CNY" // Китайский юань
	AUD = "AUD" // Австралийский доллар
	CAD = "CAD" // Канадский доллар
	SGD = "SGD" // Сингапурский доллар
	AED = "AED" // Дирхам ОАЭ
	TRY = "TRY" // Турецкая лира
	KZT = "KZT" // Казахстанский тенге
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case EUR, USD, UZS:
		return true
	}
	return false
}
