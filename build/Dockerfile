FROM golang:1.18 as build

ENV TSN_ROOT=/go/src/tsn-service
ENV CGO_ENABLED=0

RUN mkdir -p $TSN_ROOT/

COPY . $TSN_ROOT

RUN cd $TSN_ROOT && GO111MODULE=on go build -o /go/bin/tsn ./


FROM alpine:3.11
RUN apk add bash
ENV HOME=/home/tsn-service
RUN mkdir $HOME
WORKDIR $HOME

COPY --from=build /go/bin/tsn /usr/local/bin/

CMD ["main"]