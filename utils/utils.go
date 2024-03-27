package utils

import "math/rand"

func RandString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	var total = len(charset)
	var ret = ""
	for i := 0; i < length; i++ {
		var index = rand.Intn(total)
		ret = ret + charset[index:index+1]
	}
	return ret
}
