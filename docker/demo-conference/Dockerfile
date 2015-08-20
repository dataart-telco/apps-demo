FROM ubuntu

MAINTAINER gdubina@dataart.com

COPY demo-conference-server /demo/
COPY demo.gcfg /demo/

EXPOSE 7091

WORKDIR /demo
CMD ./demo-conference-server -host=$EXTERNAL_IP -redis=$REDIS_SERVICE_HOST:$REDIS_SERVICE_PORT -restcomm=$RESTCOMM_SERVICE
