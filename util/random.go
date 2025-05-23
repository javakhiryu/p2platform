package util

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomTgUsername() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(1, 1000) * 100
}

func RandomCurrency() string {
	currencies := []string{EUR, USD, UZS}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomUUID() uuid.UUID {
	return uuid.New()
}