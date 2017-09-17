package zolist

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"appengine"

	"github.com/hpaluch/zolist-go/zolist/zoapi"
)

type HomeRest struct {
	Restaurant *zoapi.Restaurant
}

type HomeModel struct {
	Now         time.Time
	Header      http.Header
	Restaurants []HomeRest
}

var homeTemplate = template.Must(template.New("home").Parse(`
<p>Now is: {{ .Now }}</p>

{{ range $i, $v := .Restaurants }}
<h2><a href="{{ $v.Restaurant.Url  }}" target="zomato">{{ $v.Restaurant.Name }}</a></h2>
<p>Debug Restaurant Id: {{ $v.Restaurant.Id }}</p>
{{ end }}

<hr>
Powered by GAE <img src='/static/appengine-silver-120x30.gif' alt='GAE' >
`))

func handler(w http.ResponseWriter, r *http.Request) {
	var api_key = os.Getenv("ZOMATO_API_KEY")
	if api_key == "" {
		http.Error(w, "Internal error - missing ZOMATO_API_KEY",
			http.StatusInternalServerError)
	}
	var ctx = appengine.NewContext(r)

	restIds := []int{18355040, // Lidak
		16513797} // Na Pude

	restModels := make([]HomeRest, len(restIds))

	for i, id := range restIds {
		restaurant, err := zoapi.FetchZomatoRestaurant(ctx, api_key, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		restModels[i].Restaurant = restaurant
	}

	homeModel := HomeModel{
		Now:         time.Now(),
		Header:      r.Header,
		Restaurants: restModels,
	}

	if err := homeTemplate.Execute(w, homeModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// main handler fo Go/GAE application
func init() {
	http.HandleFunc("/", handler)
}
