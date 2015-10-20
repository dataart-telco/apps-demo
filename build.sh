mkdir -p out
cp ./common/demo.gcfg ./out

#build bins
go get ./...
go install ./common

echo 'Build feedback'

cd ./feedback-call
go build -o ../out/feedback-call main.go advertisingtcall.go conference.go webserver.go
go build -o ../out/feedback-call-portal portal-main.go advertisingtcall.go conference.go portal-webserver.go
cd ../

echo "Build calls-consumer"
cd ./calls-consumer
go build -o ../out/calls-consumer
cd ../

echo 'Build conference-call'
cd ./conference-call
go build -o ../out/conference-call
cd ../

echo 'Build mailagent'
cd ./mailagent
go build -o ../out/mailagent
cd ../

echo 'Build sms-feedback'
cd ./sms-feedback
go build -o ../out/sms-feedback
cd ../

echo 'Build drop-conference'
cd ./drop-conference
go build -o ../out/drop-conference
cd ../

echo 'Copy files to docker dirs'
#copy do docker directory
cp ./common/demo.gcfg ./docker/conference-call/
cp ./common/demo.gcfg ./docker/calls-consumer/
cp ./common/demo.gcfg ./docker/mailagent/
cp ./common/demo.gcfg ./docker/sms-feedback/
cp ./common/demo.gcfg ./docker/feedback-call/
cp ./common/demo.gcfg ./docker/drop-conference

cp ./out/calls-consumer ./docker/calls-consumer/
cp ./out/conference-call ./docker/conference-call
cp ./out/mailagent ./docker/mailagent
cp ./out/sms-feedback ./docker/sms-feedback
cp ./out/feedback-call ./docker/feedback-call
cp ./out/feedback-call-portal ./docker/feedback-call
cp ./out/drop-conference ./docker/drop-conference

echo 'Build completed!'
