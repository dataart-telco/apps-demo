docker build -t dataart/tad-demo-call-play ./demo-call-play
docker rm demo-call-play
docker run --name demo-call-play --env HOST=192.168.176.220 -p 7092:7092 dataart/tad-demo-call-play
