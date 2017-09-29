// misc utilities for ZoList application
package zoutils

import (
	"time"

	"appengine"
)

func RoundDurationToMs( d time.Duration ) time.Duration {
	return  ((d + time.Millisecond/2) / time.Millisecond ) * time.Millisecond
}

// data model for templates/zz_layout.html
type LayoutModel struct {
	NowUTC         time.Time
	RenderTime     string
	ServerSoftware string
	Title		string
}

func CreateLayoutModel(tic time.Time, title string) LayoutModel {
	return LayoutModel{
		NowUTC:         time.Now(),
		RenderTime:     RoundDurationToMs(time.Since(tic)).String(),
		ServerSoftware: appengine.ServerSoftware(),
		Title:	title,
	}
}
