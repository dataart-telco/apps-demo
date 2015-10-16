package main
import (
	"tad-demo/common"
	"strings"
	"fmt"
)

type Conference struct{

}

func (conf Conference)GetParticipants()([]string){
	phonesSid := make([]string, 0)

	for _, i := range db.LRange(common.DB_KEY_URI, 0, 1000).Val(){
		uri := i[0:len(i) - 5]
		sid := uri[strings.LastIndex(i, "/") + 1 : len(uri)]
		phonesSid = append(phonesSid, sid)
	}
	return phonesSid;
}

func (conf Conference)Drop(){
	fmt.Println("Drop conference")
	i := db.LPop(common.DB_KEY_URI).Val()
	for i != "" {
		uri := i[0:len(i) - 5]
		restcommApi.CompleteCallByUri(uri)
		i = db.LPop(common.DB_KEY_URI).Val()
	}
}