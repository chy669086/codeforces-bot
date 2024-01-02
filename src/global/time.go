package global

import (
	"fmt"
	"time"
)

var lst int64 = 0

func CheckTime() bool {
	if time.Now().Unix()-lst >= 2 {
		fmt.Println(lst)
		lst = time.Now().Unix()
		return true
	}
	return false
}
