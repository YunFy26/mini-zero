package logx

import "time"

func getTimestamp() string {
	return time.Now().Format(timeFormat)
}
