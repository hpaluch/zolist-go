// misc utilities for ZoList application
package zoutils

import (
	"net/http"
	"time"

	"appengine"

	"golang.org/x/text/message"

	"github.com/hpaluch/zolist-go/zolist/zol10n"
)

func RoundDurationToMs(d time.Duration) time.Duration {
	return ((d + time.Millisecond/2) / time.Millisecond) * time.Millisecond
}

type BreadCrumb struct {
	Url         string
	Description string
}

type LangSwitch struct {
	LangCode string
	LangName string
	Url      string
}

// data model for templates/zz_layout.html
type LayoutModel struct {
	NowUTC         time.Time
	RenderTime     string
	ServerSoftware string
	Title          string
	BreadCrumbs    []BreadCrumb
	P              *message.Printer
	HtmlLangCode   string
	UrlBase        string
	LangSwitches   []LangSwitch
}

func CreateLayoutModel(tic time.Time, title string, bc *BreadCrumb, ctx appengine.Context, r *http.Request) (LayoutModel, error) {

	var langIndex, locPrinter, err = zol10n.LocFromUrlBase(ctx, r)
	if err != nil {
		return LayoutModel{}, err
	}

	var urlBase = zol10n.LangUrlBase(langIndex)

	var breadCrumbs = make([]BreadCrumb, 1)
	breadCrumbs[0] = BreadCrumb{
		Url:         urlBase + "/",
		Description: "ZoList",
	}

	if bc != nil {
		breadCrumbs = append(breadCrumbs, *bc)
	}

	var langSwitches = make([]LangSwitch, 0)
	for i, lang := range zol10n.UrlLangs {
		if i != langIndex {
			var langUrlBase = zol10n.LangUrlBase(i)
			var url = langUrlBase + r.URL.Path[3:]
			var langSw = LangSwitch{
				LangCode: lang,
				LangName: zol10n.LangNames[i],
				Url:      url,
			}
			langSwitches = append(langSwitches, langSw)
		}
	}

	return LayoutModel{
		NowUTC:         time.Now(),
		RenderTime:     RoundDurationToMs(time.Since(tic)).String(),
		ServerSoftware: appengine.ServerSoftware(),
		Title:          title,
		BreadCrumbs:    breadCrumbs,
		P:              locPrinter,
		HtmlLangCode:   zol10n.HtmlLangs[langIndex],
		UrlBase:        urlBase,
		LangSwitches:   langSwitches,
	}, nil
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

func NoCacheHeaders(w http.ResponseWriter) {
	// look at headers of www.seznam.cz :-)
	// WARNING! This also sets Cache-Control: no-cache
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
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
