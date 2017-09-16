package zolist

import (
    "html/template"
    "net/http"
    "time"
)

type HomeModel struct {
	Now    time.Time
        Header http.Header
}

var homeTemplate = template.Must(template.New("home").Parse(`
<p>Now is: {{ .Now }}</p>
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



func handler(w http.ResponseWriter, r *http.Request) {

    homeModel := HomeModel{
       Now:    time.Now(),
       Header: r.Header,
    }

    if err := homeTemplate.Execute(w, homeModel); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

