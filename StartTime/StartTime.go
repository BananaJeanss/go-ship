package StartTime

import "time"

var StartTime int64

func SetStartTime() {
	// set start time to current time in seconds since epoch
	StartTime = time.Now().Unix()
}

func GetStartTime() int64 {
	return StartTime
}