package main
import (
	"net/http"
	"fmt"
	"tad-demo/common"
)

const PATH_RUN_PLAY  = "/run-now"
const PATH_STAT  = "/stat"
const PATH_ACTIVE  = "/active"

const KEY_GATHER = "gather"

type HttpHandler func(http.ResponseWriter, *http.Request)

type PortalWebServer struct{
	pages map[string]HttpHandler
}

func NewPortalWebServer()(PortalWebServer){
    return PortalWebServer{pages : make(map[string]HttpHandler)}
}

func (web PortalWebServer)httpHandlerPlay(w http.ResponseWriter, r *http.Request){
	fmt.Println("\t<- http request -", r.URL)

	fmt.Println("\tpages", web.pages)
	f := web.pages[r.URL.Path]
	if f != nil {
		f(w, r)
		return
	}
	web.handleStat(w, r)
}

func (web PortalWebServer)printHeader(w http.ResponseWriter, title string ){
	fmt.Fprintf(w, "<html><head><title>%s</title></head><body><p><b>%s</b></p>", title, title)
	fmt.Fprintf(w,
		"<div>" +
		"<p><a href=\"%s\">%s</a></p>" +
		"<p><a href=\"%s\">%s</a></p>" +
		"<p><a href=\"%s\">%s</a></p>" +
		"</div>" +
		"<div><b>OUTPUT</b></div>",
		PATH_STAT, "Statistics",
		PATH_ACTIVE, "Active conference",
		PATH_RUN_PLAY, "Run script now")
}

func (web PortalWebServer)printBottom(w http.ResponseWriter){
	fmt.Fprintf(w, "</body></html>")
}

func (web PortalWebServer)handleStat(w http.ResponseWriter, r *http.Request){
	web.printHeader(w, "Users list")
	fmt.Fprintf(w, "<ol>")
	for _, e := range db.LRange(KEY_GATHER, 0, 1000).Val(){
		fmt.Fprintf(w, "<li>%s</li>", e)
	}
	fmt.Fprintf(w, "</ol>")
	web.printBottom(w)
}

func (web PortalWebServer)handleActive(w http.ResponseWriter, r *http.Request){
	web.printHeader(w, "Active conference")
	fmt.Fprintf(w, "<ol>")
	for _, e := range db.LRange(common.DB_KEY_URI, 0, 1000).Val(){
		fmt.Fprintf(w, "<li>%s</li>", e)
	}
	fmt.Fprintf(w, "</ol>")
	web.printBottom(w)
}

func (web PortalWebServer)handleRunNow(w http.ResponseWriter, r *http.Request){
	web.printHeader(w, "Script done")

	/*conference := &Conference{}
	participants := conference.GetParticipants()
	conference.Drop()

	adtCall := &AdtCall{}
	adtCall.Call(participants)*/

	fmt.Fprintf(w, "<p>DON'T WORN NOW</p>")
	web.printBottom(w)
}

func (web PortalWebServer) Init() {
	web.pages[PATH_RUN_PLAY] = web.handleRunNow
	web.pages[PATH_STAT] = web.handleStat
	web.pages[PATH_ACTIVE] = web.handleActive
}

func (web PortalWebServer) Start() {
	fmt.Println("Start advertising web server")

	web.Init()

	http.HandleFunc("/", web.httpHandlerPlay)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort.Portal), nil)
	if(err != nil){
		panic(err)
	}
}
