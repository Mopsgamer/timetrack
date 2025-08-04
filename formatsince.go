package main

import (
	"fmt"
	"time"
)

func FormatSince(duration time.Duration) string {
	seconds := int(duration.Seconds())
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24
	result := ""

	if days > 0 {
		result += fmt.Sprint(days, "d ")
	}

	result += fmt.Sprintf("%02d:%02d:%02d", hours%24, minutes%60, seconds%60)
	return result
}
