FROM golang:1.12.5

WORKDIR /go/src/github.com/vavas/packet-analyzer
COPY . /go/src/github.com/vavas/packet-analyzer

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]