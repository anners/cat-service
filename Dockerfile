FROM golang:1.13 AS build
MAINTAINER Ann Wallace annerz@gmail.com
WORKDIR  /go/src/github.com/anners/cat-service
COPY caas.go . 
RUN CGO_ENABLED=0 GOOS=linux go build -o caas .

FROM alpine:latest 
WORKDIR /root/
COPY --from=build /go/src/github.com/anners/cat-service/caas .
CMD ["./caas"]
# Document that the service listens on port 8080.
EXPOSE 8080
