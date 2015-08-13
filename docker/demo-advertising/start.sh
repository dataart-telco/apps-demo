./demo-call-play-server -host=$HOST -redis=$REDIS -restcomm=${RESTCOMM} &> call-play-server.log
./demo-call-play -host=$HOST -redis=$REDIS -restcomm=${RESTCOMM}
echo "Started with daemon mode"
tail -f call-play-server.log
