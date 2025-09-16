package util

import (
	"math/rand"
	"strings"
	"time"
)

var R *rand.Rand
var currencies = [4]string{"USD", "EUR", "CAD", "CNY"}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	R = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + R.Int63n(max-min+1)
}

// RandomString generates a random string of given length
func RandomString(length int) string {
	var sb strings.Builder
	var size = len(alphabet)

	for i := 0; i < length; i++ {
		c := alphabet[R.Intn(size)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	n := len(currencies)
	return currencies[R.Intn(n)]
}
