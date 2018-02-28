FROM golang:1.10.0

RUN apt-get update && apt-get install -y --no-install-recommends \
		unzip \
	&& rm -rf /var/lib/apt/lists/*

RUN go get github.com/nlnwa/veidemann-ws-api-gateway \
  ; cd /go/src/github.com/nlnwa/veidemann-ws-api-gateway \
 && go generate \
 && CGO_ENABLED=0 go build -tags netgo



FROM alpine:3.7

MAINTAINER Norsk nettarkiv
EXPOSE 3010
ENV CONTROLLER_HOST=localhost CONTROLLER_PORT=50051

COPY --from=0 /go/src/github.com/nlnwa/veidemann-ws-api-gateway/veidemann-ws-api-gateway /
COPY --from=0 /go/src/github.com/nlnwa/veidemann-ws-api-gateway/html /html

CMD /veidemann-ws-api-gateway -logtostderr -controller_endpoint ${CONTROLLER_HOST}:${CONTROLLER_PORT}
