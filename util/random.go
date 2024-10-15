package util

import (
	"math/rand"
	"strings"
	"time"
)

// Const for indicate the soport alphabet
const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RamdomInit is used to generate a random number between min and max
func RamdomInit(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString is used to generate a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner is used to generate a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney is used to generate a random amount of money
func RandomMoney() int64 {
	return RamdomInit(0, 1000)
}

// RandomCurrency is used to generate a random currency
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD", "GBP", "COR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
