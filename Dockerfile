FROM golang:1.12 as builder

ENV GO111MODULE=off
RUN go get -u google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go

RUN apt-get update
RUN apt-get install -y unzip

WORKDIR /protoc
RUN curl -o protoc.zip -L https://github.com/protocolbuffers/protobuf/releases/download/v3.10.0/protoc-3.10.0-linux-x86_64.zip
RUN unzip protoc.zip
RUN cp /protoc/bin/protoc /usr/local/bin/protoc


WORKDIR /go/src/github.com/ninnemana/rpc-demo
ADD . .
RUN make generate
RUN GO111MODULE=on go build -o /api ./cmd

FROM alpine

COPY --from=builder /api /api
ENTRYPOINT ["/api"]
