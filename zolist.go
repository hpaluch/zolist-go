package zolist

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"appengine"

	"github.com/hpaluch/zolist-go/zolist/zoapi"
	"github.com/hpaluch/zolist-go/zolist/zocache"
)

var (
	// from: https://github.com/golang/appengine/blob/master/demos/guestbook/guestbook.go
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	zomato_api_key = os.Getenv("ZOMATO_API_KEY")

)

type HomeRest struct {
	Restaurant *zoapi.Restaurant
	Menu       *zoapi.Menu
}

type HomeModel struct {
	Now            time.Time
	Header         http.Header
	Restaurants    []HomeRest
	RenderTime     string
	ServerSoftware string
}


func handler(w http.ResponseWriter, r *http.Request) {

	// tic code got from https://github.com/golang/appengine/blob/master/demos/guestbook/guestbook.go
	tic := time.Now()
	var ctx = appengine.NewContext(r)
	// report 404 for other path than "/"
	// see https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine_flexible/helloworld/helloworld.go
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// how to trigger this error:
	// curl -X POST -v http://localhost:8080
	if r.Method != "GET" {
		ctx.Errorf("Method '%s' not allowed for path '%s'",
			r.Method, r.URL.Path)
		http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed)
		return
	}

	restIds := []int{18355040, // Lidak
		16513797} // Na Pude

	restModels := make([]HomeRest, len(restIds))

	for i, id := range restIds {
		restaurant, err := zocache.FetchZomatoRestaurant(ctx, zomato_api_key, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		menu, err := zoapi.FetchZomatoDailyMenu(ctx, zomato_api_key, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		restModels[i].Restaurant = restaurant
		restModels[i].Menu = menu
	}

	homeModel := HomeModel{
		Now:            time.Now(),
		Header:         r.Header,
		Restaurants:    restModels,
		RenderTime:     time.Since(tic).String(),
		ServerSoftware: appengine.ServerSoftware(),
	}

	if err := tpl.ExecuteTemplate(w, "home.html", homeModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// main handler fo Go/GAE application
func init() {
	if zomato_api_key == "" {
		panic("Fatal error - missing ZOMATO_API_KEY in app.yaml")
	}

	http.HandleFunc("/", handler)
}
