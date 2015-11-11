FROM ubuntu

MAINTAINER gdubina@dataart.com


COPY sms-feedback /demo/
COPY demo.gcfg /demo/

RUN apt-get update  
RUN apt-get install -y ca-certificates

EXPOSE 30794

WORKDIR /demo
CMD ./sms-feedback  -host=$EXTERNAL_IP -redis=$REDIS_SERVICE_HOST:$REDIS_SERVICE_PORT -restcomm=$RESTCOMM_SERVICE -r-user=$RESTCOMM_USER -r-pass=$RESTCOMM_PASSWORD -r-phone-sms=$PHONE_NUMBER -rec=$REC_HOST -dp=$DID_DOMAIN
