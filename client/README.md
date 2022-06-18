docker build --tag autogo-client-test .
docker run autogo-client-test

sudo docker inspect --format '{{ .NetworkSettings.IPAddress }}'