package zolist

import (
    "fmt"
    "net/http"
    "time"
)

func init() {
    http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Future Zomato(TM) menu list, written in Go")
    fmt.Fprintln(w, "Now is: ",time.Now())
    fmt.Fprintln(w)
    fmt.Fprintln(w,"Request headers:")
    for k,v := range r.Header {
      fmt.Fprintf(w,"%s -> %s\r\n",k,v)
    }
}

