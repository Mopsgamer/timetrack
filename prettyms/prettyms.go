package prettyms

import (
	"fmt"
	"time"
)

func Time(t time.Time) string {
	return fmt.Sprintf("%f", time.Since(t).Seconds())
}
