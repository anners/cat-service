FROM golang:1.5
EXPOSE 8080
WORKDIR /go/src/app
COPY . /go/src/app

