FROM ubuntu

MAINTAINER gdubina@dataart.com

COPY mailagent /demo/
COPY demo.gcfg /demo/

RUN apt-get update  
RUN apt-get install -y ca-certificates

WORKDIR /demo
CMD ./mailagent -host=$EXTERNAL_IP -redis=$REDIS_SERVICE_HOST:$REDIS_SERVICE_PORT -restcomm=$RESTCOMM_SERVICE -user=$GMAIL_USER -pass=$GMAIL_PASS -timer=$DUMP_TIMER

