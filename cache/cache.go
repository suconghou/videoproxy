package cache

import (
	"sync"
	"time"

	"github.com/suconghou/videoproxy/db"
	"github.com/suconghou/videoproxy/util"
)

// string : cacheItem
var caches = sync.Map{}

type cacheItem struct {
	t int64
	v bool
}

func init() {
	go func() {
		for {
			var now = time.Now().Unix()
			caches.Range(func(k interface{}, v interface{}) bool {
				t, ok := v.(*cacheItem)
				if !ok || t.t < now {
					caches.Delete(k)
				}
				return true
			})
			time.Sleep(time.Minute)
		}
	}()
}

func InWhiteList(vid string) bool {
	v, ok := caches.Load(vid)
	if ok {
		return v.(*cacheItem).v
	}
	retId, exist, err := db.FindId(vid, db.TABLE_WHITELIST)
	if err != nil {
		util.Log.Print(err)
	}
	caches.Store(retId, &cacheItem{time.Now().Unix() + 3600, exist})
	return exist
}
