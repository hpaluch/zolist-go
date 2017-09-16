package zolist

import (
	"html/template"
	"fmt"
	"net/http"
	"os"
	"time"

	"appengine"
        "appengine/urlfetch"
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

func init() {
	http.HandleFunc("/", handler)
}

// restId = Restaurant ID
func fetchZomatoRestaurant(ctx appengine.Context,api_key string, restId int) (string, error) {
	var client = urlfetch.Client(ctx)
        var url = fmt.Sprintf("%s%s%d",
				"https://developers.zomato.com/api/v2.1",
				"/restaurant?res_id=",
				restId)
	// see https://stackoverflow.com/questions/12864302/how-to-set-headers-in-http-get-request 
	var req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("user-key", api_key)
        resp, err := client.Do(req)
        if err != nil {
		return "",err
        }
	var str = fmt.Sprintf("HTTP GET returned status %v", resp.Status)
	return str,nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	var api_key = os.Getenv("ZOMATO_API_KEY")
	if api_key == "" {
		http.Error(w, "Internal error - missing ZOMATO_API_KEY",
			http.StatusInternalServerError)
	}

	var restId = 18355040; // Lidak
	var ctx = appengine.NewContext(r)

        restStr,err := fetchZomatoRestaurant(ctx,api_key,restId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	homeModel := HomeModel{
		Now:    time.Now(),
		Header: r.Header,
		Info:   restStr,
	}

	if err := homeTemplate.Execute(w, homeModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
