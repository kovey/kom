package tools

import "time"

func DateTime(unix int64) string {
	if unix <= 0 {
		return ""
	}

	return time.Unix(unix, 0).Format(time.DateTime)
}
