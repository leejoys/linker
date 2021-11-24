package generator

import (
	"math/rand"
	"time"
)

const (
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
	shortLength = 10
)

func Do() string {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Nanosecond)
	abc := []byte(alphabet)
	short := []byte{}
	for i := 1; i <= shortLength; i++ {
		short = append(short, abc[rand.Intn(len(abc)-1)])
	}
	return string(short)
}
