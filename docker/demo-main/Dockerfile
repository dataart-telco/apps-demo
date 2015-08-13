FROM ubuntu

MAINTAINER gdubina@dataart.com

COPY demo-main-server /demo/
COPY demo.gcfg /demo/

EXPOSE 7090

WORKDIR /demo
CMD ./demo-main-server -host=$HOST -redis=$REDIS -restcomm=${RESTCOMM}
