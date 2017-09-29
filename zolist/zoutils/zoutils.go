// misc utilities for ZoList application
package zoutils

import (
	"time"
)

func RoundDurationToMs( d time.Duration ) time.Duration {
	return  ((d + time.Millisecond/2) / time.Millisecond ) * time.Millisecond
}

