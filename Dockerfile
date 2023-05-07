FROM golang:1.16 AS base

    RUN apt-get update 
    RUN apt-get install --yes --no-install-recommends `cat requirements.os`
    RUN apt-get autoremove --yes 
    RUN apt-get clean 
    RUN rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

    WORKDIR /app
    
#FROM base as client

    COPY .env /app/.env

    RUN go get github.com/eclipse/paho.mqtt.golang && \
        go get github.com/gorilla/websocket && \
        go get github.com/spf13/viper

    WORKDIR /app/config
    COPY ./config ./
    RUN go mod tidy

    WORKDIR /app/infrastructure
    COPY ./infrastructure ./
    RUN go mod tidy

    WORKDIR /app/client
    COPY ./client ./
    RUN go mod tidy
 
    RUN go build 

    RUN chmod +x ./client

    EXPOSE 8080

    CMD ["./client"]