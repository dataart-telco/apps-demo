FROM ubuntu

MAINTAINER gdubina@dataart.com

COPY conference-call /demo/
COPY demo.gcfg /demo/

EXPOSE 30791

WORKDIR /demo
CMD ./conference-call -host=$EXTERNAL_IP -redis=$REDIS_SERVICE_HOST:$REDIS_SERVICE_PORT -restcomm=$RESTCOMM_SERVICE -r-user=$RESTCOMM_USER -r-pass=$RESTCOMM_PASSWORD -r-phone-conf=$PHONE_NUMBER -dp=$DID_DOMAIN
