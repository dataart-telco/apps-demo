docker build -t docker_mailagent ./mailagent
docker rm docker_mailagent_1
docker run --name docker_mailagent_1 -e GMAIL_USER="tads2015dataart@gmail.com"\
    -e EXTERNAL_IP=192.168.176.220 \
    -e REDIS_SERVICE_HOST=192.168.176.220 \
    -e REDIS_SERVICE_PORT=6379 \
    -e RESTCOMM_SERVICE=192.168.176.220:7070 \
    -e GMAIL_USER="tads2015dataart@gmail.com" \
    -e GMAIL_PASS=gdubina2015 \
    -p 7094:7094 docker_mailagent
