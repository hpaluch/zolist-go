package zolist

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"appengine"

	"github.com/hpaluch/zolist-go/zolist/zoapi"
	"github.com/hpaluch/zolist-go/zolist/zoconsts"
)

type HomeModel struct {
	Now    time.Time
	Header http.Header
	Info   string
}

var homeTemplate = template.Must(template.New("home").Parse(`
<p>Now is: {{ .Now }}</p>
<p>Info: {{ .Info }}</p>
<h2>Request headers</h2>
<table>
   <tr>
      <th>Key</th><th>-</th><th>Value</th>
   </tr>
{{ range $k, $v := .Header }}
   <tr>
     <td>{{ $k }}</td><td>=&gt;</td><td>{{ $v }}</td>
   </tr>
{{ end }}
</table>
<img src='/static/appengine-silver-120x30.gif' alt='GAE' >
`))

func handler(w http.ResponseWriter, r *http.Request) {
	var api_key = os.Getenv("ZOMATO_API_KEY")
	if api_key == "" {
		http.Error(w, "Internal error - missing ZOMATO_API_KEY",
			http.StatusInternalServerError)
	}

	var restId = 18355040 // Lidak
	var ctx = appengine.NewContext(r)

	restStr, err := zoapi.FetchZomatoRestaurant(ctx, api_key, restId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var info = fmt.Sprintf("ZoRest: %s, MaxItems: %d", restStr,
		zoconsts.ZoMaxRestItems)

	homeModel := HomeModel{
		Now:    time.Now(),
		Header: r.Header,
		Info:   info,
	}

	if err := homeTemplate.Execute(w, homeModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// main handler fo Go/GAE application
func init() {
	http.HandleFunc("/", handler)
}
