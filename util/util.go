package util

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("asdfghjklqwertyuiopzxcvbnmASDFGHJKLQWERTYUIOPZXCVBNM")
	result := make([]byte, n)

	rand.New(rand.NewSource(time.Now().Unix()))
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
