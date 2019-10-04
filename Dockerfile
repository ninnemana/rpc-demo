FROM golang:1.12

RUN go get -u google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go

RUN apt-get update
RUN apt-get install -y unzip

WORKDIR /protoc
RUN curl -o protoc.zip -L https://github.com/protocolbuffers/protobuf/releases/download/v3.10.0/protoc-3.10.0-linux-x86_64.zip
RUN unzip protoc.zip
RUN cp /protoc/bin/protoc /usr/local/bin/protoc


WORKDIR /go/src/github.com/ninnemana/rpc-demo
COPY vinyltap.proto vinyltap.proto
COPY prototool.yaml prototool.yaml
COPY Makefile Makefile
COPY cmd/main.go cmd/main.go

RUN make generate
RUN go run ./cmd/main.go
