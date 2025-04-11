package util

const (
	cash = "cash"
	card = "card"
)

func IsSupportedSource(source string) bool {
	switch source {
	case cash, card:
		return true
	}
	return false
}
