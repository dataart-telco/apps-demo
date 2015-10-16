FROM ubuntu

MAINTAINER gdubina@dataart.com

COPY calls-consumer /demo/
COPY demo.gcfg /demo/

EXPOSE 7090

WORKDIR /demo
CMD ./calls-consumer -host=$EXTERNAL_IP -redis=$REDIS_SERVICE_HOST:$REDIS_SERVICE_PORT -restcomm=$RESTCOMM_SERVICE
