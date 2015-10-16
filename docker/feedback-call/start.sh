./feedback-call-portal -host=$EXTERNAL_IP -redis=$REDIS_SERVICE_HOST:$REDIS_SERVICE_PORT -restcomm=$RESTCOMM_SERVICE &> call-play-server.log
./feedback-call -host=$EXTERNAL_IP -redis=$REDIS_SERVICE_HOST:$REDIS_SERVICE_PORT -restcomm=$RESTCOMM_SERVICE
echo "Started with daemon mode"
