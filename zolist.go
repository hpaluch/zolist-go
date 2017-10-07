package zolist

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"appengine"

	"github.com/hpaluch/zolist-go/zolist/zoapi"
	"github.com/hpaluch/zolist-go/zolist/zocache"
	"github.com/hpaluch/zolist-go/zolist/zoconsts"
	"github.com/hpaluch/zolist-go/zolist/zol10n"
	"github.com/hpaluch/zolist-go/zolist/zoutils"
)

func tplCzDateStr(timeArg interface{}) (string, error) {
	// Type Assertion - please see:
	//    https://stackoverflow.com/questions/14289256/cannot-convert-data-type-interface-to-type-string-need-type-assertion
	t, ok := timeArg.(time.Time)
	if !ok {
		var errMsg = fmt.Sprintf("Unsupported argument type: %T, expecting time.Time", timeArg)
		return "", errors.New(errMsg)
	}

	return t.In(zoconsts.CzechLocation).Format("02.01.2006 15:04:05 MST"), nil
}

func tplCzDateStrWithAgo(timeArg interface{}) (string, error) {
	var dateStr, err = tplCzDateStr(timeArg)
	if err != nil {
		return "", err
	}
	// timeArg was verified by tplCzDateStr
	t, _ := timeArg.(time.Time)

	// compute ago
	var czNow = time.Now().In(zoconsts.CzechLocation)
	var duration = czNow.Sub(t)
	// round to millisecond
	// see: http://grokbase.com/t/gg/golang-nuts/1492epp0qb/go-nuts-how-to-round-a-duration
	duration = zoutils.RoundDurationToMs(duration)
	var czAgo = time.Duration(duration).String()

	var str = fmt.Sprintf("%s (%s ago)", dateStr, czAgo)
	return str, nil
}

var (
	tplFn = template.FuncMap{
		"ZoCzDateFormat":        tplCzDateStr,
		"ZoCzDateFormatWithAgo": tplCzDateStrWithAgo,
	}

	// from: https://github.com/golang/appengine/blob/master/demos/guestbook/guestbook.go
	tpl = template.Must(template.New("").Funcs(tplFn).ParseGlob("templates/*.html"))

	zomato_api_key = os.Getenv("ZOMATO_API_KEY")
	str_rest_ids   = os.Getenv("REST_IDS")
	// initialized in init()
	rest_ids []int
)

type DetailMenuModel struct {
	LayoutModel zoutils.LayoutModel
	Restaurant  *zoapi.Restaurant
	Menu        *zoapi.Menu
	NextId      *int
	PrevId      *int
}

var reDetailPath = regexp.MustCompile(`^/[a-z]{2}/menu/(\d{1,12})$`)

// menu detail for rest_id
func handlerDetail(w http.ResponseWriter, r *http.Request) {
	var tic = time.Now()
	zoutils.NoCacheHeaders(w)
	var ctx = appengine.NewContext(r)

	if !zoutils.VerifyGetMethod(ctx, w, r) {
		return
	}
	// necessary to verify relaxed lang url prefix
	var langIndex, err = zol10n.LangFromUrlBase(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !reDetailPath.MatchString(r.URL.Path) {
		ctx.Errorf("Path '%s' does not match regexp", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	var groups = reDetailPath.FindStringSubmatch(r.URL.Path)
	if len(groups) != 2 {
		ctx.Errorf("Got unexpected number of groups %d <> 2: %v", len(groups), groups)
		http.Error(w, "regexp internal error", http.StatusInternalServerError)
		return

	}
	var strRestId = groups[1]

	id, err := strconv.Atoi(strRestId)
	if err != nil {
		ctx.Errorf("Unable to convert '%s' to int: %v", strRestId, err)
		http.Error(w, "Can't parse ID to int", http.StatusInternalServerError)
	}

	var list_index = zoutils.SearchIntArray(rest_ids, id)
	if list_index == -1 {
		ctx.Errorf("id (%d) not found in list: %v", id, rest_ids)
		http.NotFound(w, r)
		return
	}
	var next_id *int = nil
	if list_index+1 < len(rest_ids) {
		next_id = &rest_ids[list_index+1]
	}
	var prev_id *int = nil
	if list_index > 0 {
		prev_id = &rest_ids[list_index-1]
	}

	restaurant, err := zocache.FetchZomatoRestaurant(ctx, zomato_api_key, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	menu, err := zocache.FetchZomatoDailyMenu(ctx, zomato_api_key, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var urlBase = zol10n.LangUrlBase(langIndex)
	var bc = zoutils.BreadCrumb{
		Url:         fmt.Sprintf("%s/menu/%d", urlBase, id),
		Description: "Menu Detail",
	}

	var loc = zol10n.LocFromIndex(langIndex)
	var title = loc.Sprintf("Detail of %s", restaurant.Name)

	layoutModel, err := zoutils.CreateLayoutModel(tic, title, &bc, ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var model = DetailMenuModel{
		LayoutModel: layoutModel,
		Restaurant:  restaurant,
		Menu:        menu,
		NextId:      next_id,
		PrevId:      prev_id,
	}

	if err := tpl.ExecuteTemplate(w, "detail.html", model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type HomeRest struct {
	Restaurant *zoapi.Restaurant
	Menu       *zoapi.Menu
}

type HomeModel struct {
	LayoutModel zoutils.LayoutModel
	Restaurants []HomeRest
}

func handlerList(w http.ResponseWriter, r *http.Request) {
	// tic code got from https://github.com/golang/appengine/blob/master/demos/guestbook/guestbook.go
	var tic = time.Now()
	var ctx = appengine.NewContext(r)
	var langIndex, err = zol10n.LangFromUrlBase(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var urlBase = zol10n.LangUrlBase(langIndex)
	var myPath = urlBase + "/"

	zoutils.NoCacheHeaders(w)
	// report 404 for other path than "/"
	// see https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine_flexible/helloworld/helloworld.go
	if r.URL.Path != myPath {
		ctx.Errorf("Unexpected path '%s' <> '%s'", r.URL.Path, myPath)
		http.NotFound(w, r)
		return
	}

	if !zoutils.VerifyGetMethod(ctx, w, r) {
		return
	}

	restModels := make([]HomeRest, len(rest_ids))

	for i, id := range rest_ids {
		restaurant, err := zocache.FetchZomatoRestaurant(ctx, zomato_api_key, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		menu, err := zocache.FetchZomatoDailyMenu(ctx, zomato_api_key, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		restModels[i].Restaurant = restaurant
		restModels[i].Menu = menu
	}

	var loc = zol10n.LocFromIndex(langIndex)
	layoutModel, err := zoutils.CreateLayoutModel(tic, loc.Sprintf("Favorite Restaurants menu"), nil, ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	homeModel := HomeModel{
		LayoutModel: layoutModel,
		Restaurants: restModels,
	}

	if err := tpl.ExecuteTemplate(w, "home.html", homeModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// reirect to /xx/ URL where xx is one of supported langs
func handlerHome(w http.ResponseWriter, r *http.Request) {

	zoutils.NoCacheHeaders(w)
	var ctx = appengine.NewContext(r)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if !zoutils.VerifyGetMethod(ctx, w, r) {
		return
	}

	var langIndex = zol10n.LangFromHeader(ctx, r)

	var urlBase = fmt.Sprintf("/%s/", zol10n.UrlLangs[langIndex])
	http.Redirect(w, r, urlBase, http.StatusFound)
}

// main handler fo Go/GAE application
func init() {

	if zomato_api_key == "" {
		panic("Fatal error - missing/empty ZOMATO_API_KEY in app.yaml")
	}

	if str_rest_ids == "" {
		panic("Fatal error - missing/empty REST_IDS in app.yaml")
	}

	var arrIds = strings.Split(str_rest_ids, ",")
	if len(arrIds) == 0 {
		panic("No id found in REST_IDS")
	}
	rest_ids = make([]int, len(arrIds))
	for i, v := range arrIds {
		id, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("Unable to parse '%s' as Int: %v",
				v, err))
		}
		rest_ids[i] = id
	}

	for _, lang := range zol10n.UrlLangs {
		var urlBase = fmt.Sprintf("/%s", lang)
		http.HandleFunc(urlBase+"/menu/", handlerDetail)
		http.HandleFunc(urlBase+"/", handlerList)
	}

	http.HandleFunc("/", handlerHome)
}
