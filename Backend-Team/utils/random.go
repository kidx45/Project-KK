package utils

import (
	"math/rand"
)

func RandomLetterString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomNumberString(n int) string {
	const numbers = "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(b)
}

func RandomInt(min int64, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomUserName() string {
	return RandomLetterString(int(RandomInt(5, 10)))
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR"}
	return currencies[rand.Intn(len(currencies))]
}

func RandomEmail() string {
	return RandomLetterString(int(RandomInt(5, 10))) + "@" + RandomLetterString(int(RandomInt(5, 10))) + ".com"
}

func RandomFullName() string {
	return RandomLetterString(int(RandomInt(7, 10))) + " " + RandomLetterString(int(RandomInt(7, 10)))
}

func RandomPassword() string {
	return RandomLetterString(int(RandomInt(8, 16)))
}

func RandomPhoneNumber() string {
	return "+" + RandomNumberString(int(RandomInt(9, 15)))
}