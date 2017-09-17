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

const zomatoUrlBase = "https://developers.zomato.com/api/v2.1"

// internal helper function to call Zomato API and returns body as string
func fetchZomatoBody(ctx appengine.Context, api_key string, urlAppend string) ([]byte, error) {
	var client = urlfetch.Client(ctx)
	var url = fmt.Sprintf("%s%s",
		zomatoUrlBase,
		urlAppend)
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
	ctx.Infof("Body for %s: %s", url, body)
	return body, nil

}

// please see https://golang.org/pkg/encoding/json/
// for more info about json: tags
type Restaurant struct {
	Id   int    `json:"id,string"` // Ooops, they have "id":"123" in quotes (should be int)!
	Name string `json:"name"`
	Url  string `json:"url"`
}

// restId = Restaurant ID
func FetchZomatoRestaurant(ctx appengine.Context, api_key string, restId int) (*Restaurant, error) {

	var urlAppend = fmt.Sprintf("%s%d",
		"/restaurant?res_id=",
		restId)

	body, err := fetchZomatoBody(ctx, api_key, urlAppend)

	var zoApiRest = Restaurant{}
	err = json.Unmarshal(body, &zoApiRest)
	if err != nil {
		return nil, err
	}

	return &zoApiRest, nil
}

/*
{
  "daily_menu": [
    {
      "daily_menu_id": "16507624",
      "name": "Vinohradský pivovar",
      "start_date": "2016-03-08 11:00",
      "end_date": "2016-03-08 15:00",
      "dishes": [
        {
          "dish_id": "104089345",
          "name": "Tatarák ze sumce s toustem",
          "price": "149 Kč"
        }
      ]
    }
  ]
}
*/

// Uch, the Json data are a bit weird...
type MenuItemItem struct {
	Id        int    `json:"daily_menu_id,string"` // Ooops, they have "id":"123" in quotes (should be int)!
	Name      string `json:"name"`
	StartDate string `json:"start_date"` // TODO: Date object
	EndDate   string `json:"end_date"`   // TODO: Date object
}

type MenuItem struct {
	MenuItemItem MenuItemItem `json:"daily_menu"`
}

type Menu struct {
	MenuItem []MenuItem `json:"daily_menus"`
}

// restId = Restaurant ID
func FetchZomatoDailyMenu(ctx appengine.Context, api_key string, restId int) (*Menu, error) {

	var urlAppend = fmt.Sprintf("%s%d",
		"/dailymenu?res_id=",
		restId)

	body, err := fetchZomatoBody(ctx, api_key, urlAppend)

	var zoApiMenu = Menu{}
	err = json.Unmarshal(body, &zoApiMenu)
	if err != nil {
		return nil, err
	}

	return &zoApiMenu, nil
}
