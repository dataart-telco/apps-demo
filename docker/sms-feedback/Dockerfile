FROM ubuntu

MAINTAINER gdubina@dataart.com

COPY sms-feedback /demo/
COPY demo.gcfg /demo/

EXPOSE 30794

WORKDIR /demo
CMD ./sms-feedback  -host=$EXTERNAL_IP -redis=$REDIS_SERVICE_HOST:$REDIS_SERVICE_PORT -restcomm=$RESTCOMM_SERVICE
