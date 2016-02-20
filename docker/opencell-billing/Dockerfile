FROM ubuntu

MAINTAINER achikin@dataart.com

COPY opencell-billing /demo/
COPY demo.gcfg /demo/

WORKDIR /demo
CMD ./opencell-billing -redis="$REDIS_SERVICE_HOST:$REDIS_SERVICE_PORT" -restcomm="$RESTCOMM_SERVICE" -r-user="$RESTCOMM_USER" -r-pass="$RESTCOMM_PASSWORD" -oc-host="$OPENCELL_HOST" -oc-user="$OPENCELL_USER" -oc-pass="$OPENCELL_PASSWORD"
