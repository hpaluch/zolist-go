// Zomato API calls to fetch all restaurant/menu info
// Please see https://developers.zomato.com/documentation
package zoapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

// please see https://golang.org/pkg/encoding/json/
// for more info about json: tags
type Restaurant struct {
	Id   int    `json:"id,string"` // Ooops, they have "id":"123" in quotes (should be int)!
	Name string `json:"name"`
	Url  string `json:"url"`
}

// restId = Restaurant ID
func FetchZomatoRestaurant(ctx appengine.Context, api_key string, restId int) (*Restaurant, error) {
	var client = urlfetch.Client(ctx)
	var url = fmt.Sprintf("%s%s%d",
		"https://developers.zomato.com/api/v2.1",
		"/restaurant?res_id=",
		restId)
	// see https://stackoverflow.com/questions/12864302/how-to-set-headers-in-http-get-request
	var req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("user-key", api_key)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	const OkHttpStatus = 200
	// https://blog.alexellis.io/golang-json-api-client/
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != OkHttpStatus {
		return nil, errors.New(fmt.Sprintf("API call %s returned unexpected status %d <> %d, body: %s", url, resp.Status, OkHttpStatus, body))
	}

	var zoApiRest = Restaurant{}
	err = json.Unmarshal(body, &zoApiRest)
	if err != nil {
		return nil, err
	}

	return &zoApiRest, nil
}
