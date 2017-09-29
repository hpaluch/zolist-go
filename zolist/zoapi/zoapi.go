// Zomato API calls to fetch all restaurant/menu info
// Please see https://developers.zomato.com/documentation
package zoapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

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
	// ctx.Infof("Body for %s: %s", url, body)
	return body, nil

}

// please see https://golang.org/pkg/encoding/json/
// for more info about json: tags
type Restaurant struct {
	Id      int    `json:"id,string"` // Ooops, they have "id":"123" in quotes (should be int)!
	Name    string `json:"name"`
	Url     string `json:"url"`
	Changed time.Time
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
		ctx.Errorf("Unable to Unmarshall '%s': %v", body, err)
		return nil, err
	}

	zoApiRest.Changed = time.Now()
	return &zoApiRest, nil
}

/*
Real data (SAPI doc is a bit outdated:

{"daily_menus":[{"daily_menu":{"daily_menu_id":"19148688","start_date":"2017-09-17 00:00:00","end_date":"2017-09-17 23:59:59","name":"","dishes":[{"dish":{"dish_id":"659350170","name":"V\u00fdb\u011br z klasick\u00e9ho j\u00eddeln\u00edho l\u00edstku.","price":""}},{"dish":{"dish_id":"659350171","name":"T\u011b\u0161\u00edme se na Va\u0161i n\u00e1v\u0161t\u011bvu.","price":""}}]}}],"status":"success"}
*/

// Uch, the Json data are a bit weird...
type Dish struct {
	Id    int    `json:"dish_id,string"`
	Name  string `json:"name"`
	Price string `json:"price"` // should be float, but can be empty!!!
}

type Dishes struct {
	Dish Dish `json:"dish"`
}

type MenuItemItem struct {
	Id        int      `json:"daily_menu_id,string"` // Ooops, they have "id":"123" in quotes (should be int)!
	Name      string   `json:"name"`
	StartDate string   `json:"start_date"` // TODO: Date object
	EndDate   string   `json:"end_date"`   // TODO: Date object
	Dishes    []Dishes `json:"dishes"`
}

// zomato date match in format "YYYY-MM-DD"
var reZoDate = regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2}) (\d\d):(\d\d):(\d\d)`)

func (mii MenuItemItem) MenuCompactDate() string {
	// if startTime is 00:00:00 and endTime is 23:59:59 return whole day
	if strings.HasSuffix(mii.StartDate, "00:00:00") && strings.HasSuffix(mii.EndDate, "23:59:59") && len(mii.StartDate) >= 10 {
		if reZoDate.MatchString(mii.StartDate) {
			return reZoDate.ReplaceAllString(mii.StartDate, "$3.$2.$1")
		}
		return mii.StartDate[:10]
	}
	return fmt.Sprintf("%s - %s", mii.StartDate, mii.EndDate)
}

type MenuItem struct {
	MenuItemItem MenuItemItem `json:"daily_menu"`
}

type Menu struct {
	MenuItem []MenuItem `json:"daily_menus"`
	Status   string     `json:"status"`
	Changed  time.Time
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

	zoApiMenu.Changed = time.Now()
	return &zoApiMenu, nil
}
