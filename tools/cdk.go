package tools

import (
	"math/rand"
	"time"
)

const (
	chars = "23456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ"
)

func Cdk(num int32, length int) map[string]bool {
	res := make(map[string]bool, num)
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		key := cdk(length, rd)
		if _, ok := res[key]; ok {
			continue
		}

		res[key] = true
		num--

		if num <= 0 {
			break
		}
	}

	return res
}

func cdk(length int, rand *rand.Rand) string {
	max := len(chars)
	res := make([]byte, length)
	for i := 0; i < length; i++ {
		res[i] = chars[rand.Intn(max)]
	}

	return string(res)
}
