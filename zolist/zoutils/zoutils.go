// misc utilities for ZoList application
package zoutils

import (
	"net/http"
	"time"

	"appengine"
)

func RoundDurationToMs(d time.Duration) time.Duration {
	return ((d + time.Millisecond/2) / time.Millisecond) * time.Millisecond
}

type BreadCrumb struct {
	Url	string
	Description string
}

// data model for templates/zz_layout.html
type LayoutModel struct {
	NowUTC         time.Time
	RenderTime     string
	ServerSoftware string
	Title          string
	BreadCrumbs	[]BreadCrumb
}

func CreateLayoutModel(tic time.Time, title string,bc *BreadCrumb) LayoutModel {

	var breadCrumbs = make([]BreadCrumb,1)
	breadCrumbs[0] = BreadCrumb{
		Url: "/",
		Description: "ZoList",	
	}

	if bc != nil {
		breadCrumbs = append(breadCrumbs,*bc)
	}


	return LayoutModel{
		NowUTC:         time.Now(),
		RenderTime:     RoundDurationToMs(time.Since(tic)).String(),
		ServerSoftware: appengine.ServerSoftware(),
		Title:          title,
		BreadCrumbs:	breadCrumbs,
	}
}

func VerifyGetMethod(ctx appengine.Context, w http.ResponseWriter, r *http.Request) bool {

	// how to trigger this error:
	// curl -X POST -v http://localhost:8080
	if r.Method != "GET" {
		ctx.Errorf("Method '%s' not allowed for path '%s'",
			r.Method, r.URL.Path)
		http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed)
		return false
	}
	return true
}

// returns -1 when not found
func SearchIntArray(arr []int, key int) int {
	for i, v := range arr {
		if v == key {
			return i
		}
	}
	return -1
}
