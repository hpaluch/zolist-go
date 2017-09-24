// ZoList constants
package zoconsts

import (
	"fmt"
	"time"
)

// maximum restaurants per user
const ZoMaxRestItems = 10

var (
	// initialized in init()
	CzechLocation *time.Location
)

func init() {
	var err error
	CzechLocation, err = time.LoadLocation("Europe/Prague")
	if err != nil {
		panic(fmt.Sprintf("Fatal error - unable to load timezone: %v", err))
	}
}
