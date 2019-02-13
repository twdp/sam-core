package cache

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	cache2 "github.com/goburrow/cache"
	"tianwei.pro/beego-guava"
	"tianwei.pro/sam-core/config"
	"time"
)

func NewCache() cache.Cache {
	cacheName := config.Conf.DefaultString("cache.type", "guava")
	if cacheName == "guava" {
		cache := cache2.NewLoadingCache(func(key cache2.Key) (value cache2.Value, e error) {
			return nil, nil
		}, cache2.WithMaximumSize(1000),
			cache2.WithExpireAfterAccess(30 * time.Minute),)
		return beego_guava.NewGuava(cache)
	} else {
		cacheConfig := beego.AppConfig.String("cacheConfig")
		if cache, err := cache.NewCache(cacheName, cacheConfig); err != nil {
			logs.Error("init cache failed. err: %v", err)
			panic(err)
		} else {
			return cache
		}
	}
}

func NewCacheDiyFunc(loaderFunc cache2.LoaderFunc) cache.Cache {
	cacheName := config.Conf.DefaultString("cache.type", "guava")
	if cacheName == "guava" {
		cache := cache2.NewLoadingCache(loaderFunc, cache2.WithMaximumSize(1000),
			cache2.WithExpireAfterAccess(30 * time.Minute),)
		return beego_guava.NewGuava(cache)
	} else {
		panic("非guavaCache不支持添加loaderFunc")
	}
}