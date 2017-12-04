FROM golang:1.9.2-alpine

ADD ./src/clustertruck /go/src/clustertruck
ADD ./src/main /go/src/main

RUN go install main

ENTRYPOINT /go/bin/main

EXPOSE 8090
