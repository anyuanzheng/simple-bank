package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(len int) string {
	builder := strings.Builder{}

	for i := 0; i < len; i += 1 {
		builder.WriteByte(alphabet[rand.Intn(26)])
	}

	return builder.String()
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

func RandomName() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{USD, EUR, RMB}
	return currencies[rand.Intn(len(currencies))]
}
