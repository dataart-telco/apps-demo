package common

import (
	"gopkg.in/redis.v3"
	"fmt"
)

func NewDbClient(ip string)(db *redis.Client){
	db = redis.NewClient(&redis.Options{
		Addr:    ip,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := db.Ping().Result()
	if err != nil {
		fmt.Println(pong, err)
	}
	return db
}