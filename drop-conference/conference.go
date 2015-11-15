package main

import (
	"strings"
	"tad-demo/common"
)

var cfg = common.NewConfig()
var db = common.NewDbClient(cfg.Service.Redis)
var restcommApi = common.NewRestcommApi(cfg.Service.Restcomm, cfg.Auth.User, cfg.Auth.Pass)

type Conference struct {
}

func (conf Conference) GetParticipants() []string {
	phonesSid := make([]string, 0)

	for _, i := range db.LRange(common.DB_KEY_URI, 0, 1000).Val() {
		uri := i[0 : len(i)-5]
		sid := uri[strings.LastIndex(uri, "/")+1 : len(uri)]
		phonesSid = append(phonesSid, sid)
	}
	return phonesSid
}

func (conf Conference) Drop() {
	common.Info.Println("Drop conference")
	for _, i := range db.LRange(common.DB_KEY_URI, 0, 1000).Val() {
		uri := i[0 : len(i)-5]
		dropped := restcommApi.CompleteCallByUri(uri)
		if dropped {
			sid := uri[strings.LastIndex(uri, "/")+1 : len(uri)]
			conf.NotifySms(sid)
			db.LRem(common.DB_KEY_URI, 0, i)
		} else {
			common.Error.Println("Can't drop call: ", uri)
		}
	}
}

func (conf Conference) NotifySms(sid string) {
	to := db.Get(sid).Val()
	if to == "" {
		common.Info.Println("\t to is EMPTY for", sid)
		return
	}
	common.Info.Println("Notify sms: " + to)
	db.Publish(cfg.Redis.ConfChannel, to)
}
