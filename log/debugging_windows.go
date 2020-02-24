package wmlog

import "log"

func DLog(v ...interface{}) {
	log.Println(v...)
}
