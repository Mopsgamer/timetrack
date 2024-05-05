package prettyms

import (
	"fmt"
	"time"
)

func Time(t time.Time) string {
	return fmt.Sprintf("%v", time.Since(t))
}
