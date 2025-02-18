package random

import (
	"math/rand"
	"time"
)

func NewRandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	se := make([]rune, length)
	for i := range se {
		se[i] = letters[rnd.Intn(len(letters))]
	}
	return string(se)
}
