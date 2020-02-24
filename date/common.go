package date

import (
	"crypto/sha1"
	"fmt"
	"time"
)

func getSHA1(URL string) string {
	h := sha1.New()
	h.Write([]byte(URL))
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum
}

func GetSeoulTime() time.Time {
	l, e := time.LoadLocation("Asia/Seoul")
	if e == nil {
		return time.Now().In(l)
	}
	return time.Now()
}
