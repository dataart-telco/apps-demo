FROM ubuntu

MAINTAINER gdubina@dataart.com

COPY demo-call-play-server /demo/
COPY demo-call-play /demo/
COPY demo.gcfg /demo/

COPY start.sh /demo/

EXPOSE 7092
EXPOSE 7093

WORKDIR /demo

RUN chmod 777 ./start.sh

CMD ./start.sh
