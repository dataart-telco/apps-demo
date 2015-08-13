./demo-call-play-server -host=$HOST -redis=$REDIS &> call-play-server.log
./demo-call-play -host=$HOST -redis=$REDIS
echo "Started with daemon mode"
tail -f call-play-server.log
