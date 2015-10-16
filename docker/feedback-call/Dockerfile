FROM ubuntu

MAINTAINER gdubina@dataart.com

COPY feedback-call /demo/
COPY feedback-call-portal /demo/
COPY demo.gcfg /demo/

COPY start.sh /demo/

EXPOSE 30792
EXPOSE 30793

WORKDIR /demo

RUN chmod 777 ./start.sh

CMD ./start.sh
