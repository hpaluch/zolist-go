// caching some API calls
package zocache

import (
	"fmt"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/memcache"

	"github.com/hpaluch/zolist-go/zolist/zoapi"
)

// we version key to avoid possible class deserialization incompatibilities
const restaurantClassVersion = 1

func genRestaurantKey(ctx appengine.Context, restId int) *datastore.Key {
	var strKind = fmt.Sprintf("Restaurant%d", restaurantClassVersion)
	var strKey = fmt.Sprintf("%d", restId)
	// Do NOT use int key - it may collide with other enitites!!!
	return datastore.NewKey(ctx, strKind, strKey, 0, nil)
}

func FetchZomatoRestaurant(ctx appengine.Context, api_key string, restId int) (*zoapi.Restaurant, error) {

	var key = genRestaurantKey(ctx, restId)
	entity := new(zoapi.Restaurant)
	err := datastore.Get(ctx, key, entity)
	if err != nil && err != datastore.ErrNoSuchEntity {
		ctx.Errorf("Error getting Restaurant for key '%v': %v",
			key, err)
		return nil, err
	}
	// err nil = entity Get() successfully
	if err == nil {
		// ctx.Debugf("Cache hit %v", entity)
		return entity, nil
	}

	ctx.Warningf("Cache MISS for restaurant_id=%d", restId)
	// not in database - fetch using Zomato API
	// ??? why I can't use "entity" anymore?
	entity2, err := zoapi.FetchZomatoRestaurant(ctx, api_key, restId)
	if err != nil {
		return nil, err
	}

	// TODO: theoretically we could use our data anyway...
	if _, err := datastore.Put(ctx, key, entity2); err != nil {
		ctx.Errorf("Error putting Restaurant for key '%v': %v",
			key, err)
		return nil, err
	}
	return entity2, nil
}

const menuMinTtlSecs = 3600 // one hour

func FetchZomatoDailyMenu(ctx appengine.Context, api_key string, restId int) (*zoapi.Menu, error) {
	var cacheKey = fmt.Sprintf("menu_by_rest_id:%d", restId)
	var menuGet zoapi.Menu
	var menuApi *zoapi.Menu

	// NOTE: we use Gob (binary codec), because some fields
	//       have not JSON mapping.
	_, err := memcache.Gob.Get(ctx, cacheKey, &menuGet)
	if err == nil {
		return &menuGet, nil
	}

	if err != memcache.ErrCacheMiss {
		ctx.Errorf("memcacheGet('%s'): %v", cacheKey, err)
		return nil, err
	}
	// now we are sure there sis ErrCacheMiss...
	menuApi, err = zoapi.FetchZomatoDailyMenu(ctx, api_key, restId)
	if err != nil {
		return nil, err
	}

	var item = memcache.Item{
		Key:        cacheKey,
		Object:     menuApi,
		Expiration: time.Duration(menuMinTtlSecs * time.Second),
	}
	// use Info instead of warning - we have expiration - so it is common
	ctx.Infof("Menu Cache MISS for restaurant_id=%d, duration=%v", restId, item.Expiration)

	err = memcache.Gob.Set(ctx, &item)
	if err != nil {
		ctx.Errorf("Unable to memcache.Gob.Set('%s'): %v", cacheKey, err)
		return nil, err
	}

	return menuApi, nil
}
