package time_utils

import "time"

func Now() int64 {
	return time.Now().Unix()
}