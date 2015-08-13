FROM ubuntu

MAINTAINER gdubina@dataart.com

COPY demo-conference-server /demo/
COPY demo.gcfg /demo/

EXPOSE 7091

WORKDIR /demo
CMD ./demo-conference-server -host=$HOST -redis=$REDIS -restcomm=${RESTCOMM}
