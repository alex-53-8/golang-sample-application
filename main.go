package main

import (
	"rest_app/cache"
	"rest_app/cfg"
	"rest_app/database"
	"rest_app/rest"
)

func main() {
	var cfg = cfg.ConfigurationKeys{}

	var _cache = cache.RedisCache{RedisConnectionUrl: cfg.RedisConnectionUrl()}
	_cache.Init()

	var db = database.DatabaseService{DbConnectionUrl: cfg.DbConnectionUrl()}
	db.Init()

	rest.CreateRestServer(&_cache, &db, &cfg).Start()
}
