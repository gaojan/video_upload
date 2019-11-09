package utils

import "time"

func GetCurrentTime() string {
	UnixTime := time.Now().Unix()
	dataTimeStr := time.Unix(UnixTime, 0).Format("2006-01-02 15:04:05")
	return dataTimeStr
}
