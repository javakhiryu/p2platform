package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

func VerifyTelegramAuth(data map[string]string, hash string, botToken string) bool {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dataCheckStrings []string
	for _, key := range keys {
		dataCheckStrings = append(dataCheckStrings, key+"="+data[key])
	}
	dataCheckString := strings.Join(dataCheckStrings, "\n")

	secretKey := sha256.Sum256([]byte(botToken))

	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	return strings.EqualFold(hash, calculatedHash)
}
