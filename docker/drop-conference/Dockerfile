FROM ubuntu

MAINTAINER gdubina@dataart.com

COPY drop-conference /demo/
COPY demo.gcfg /demo/

WORKDIR /demo
CMD ./drop-conference -host=$EXTERNAL_IP -redis=$REDIS_SERVICE_HOST:$REDIS_SERVICE_PORT -restcomm=$RESTCOMM_SERVICE -r-user=$RESTCOMM_USER -r-pass=$RESTCOMM_PASSWORD
