package util

import (
	"math/rand"
	"time"
)

func RandomString(i int) string {
	var letters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, i)
	rand.Seed(time.Now().Unix())
	for k := range result {
		result[k] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
