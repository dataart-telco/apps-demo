cp ./common/demo.gcfg ./docker/demo-advertising/
cp ./common/demo.gcfg ./docker/demo-conference/
cp ./common/demo.gcfg ./docker/demo-main/

go install ./common
go build -o ./docker/demo-advertising/demo-call-play ./advertising/demo-call-play.go ./advertising/advertisingtcall.go ./advertising/conference.go ./advertising/webserver.go
go build -o ./docker/demo-advertising/demo-call-play-server ./advertising/demo-call-play-server.go ./advertising/advertisingtcall.go ./advertising/conference.go ./advertising/portal-webserver.go
go build -o ./docker/demo-conference/demo-conference-server ./conference/demo-conference-server.go ./conference/conference.go ./conference/storage.go ./conference/webserver.go
go build -o ./docker/demo-main/demo-main-server ./main/demo-main-server.go ./main/sms.go
