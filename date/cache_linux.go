package date

import (
	"io/ioutil"
	"time"
)

func Load(URL string) string {
	filename := "/tmp/" + getSHA1(URL) + ".log"
	if b, e := ioutil.ReadFile(filename); e == nil {
		return string(b)
	}
	return ""
}

func Save(URL string) {
	filename := "/tmp/" + getSHA1(URL) + ".log"
	t := GetSeoulTime()
	ioutil.WriteFile(filename, []byte(t.Format(time.RFC3339)), 0644)
}
