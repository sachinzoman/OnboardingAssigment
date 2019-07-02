FROM golang:1.12.1-alpine

# install git
RUN apk update && apk add git

# Install protocol buffer
#Prerequisite for grpc and protocol buffer
RUN apk add autoconf automake libtool curl make g++ unzip

RUN apk add protobuf

# setup the working directory
RUN export "GO_PATH=~/go" && export "PATH=$PATH:/$GO_PATH/bin"

# install required go packages
RUN go get -u github.com/go-sql-driver/mysql
RUN go get -u google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go

WORKDIR /go/src/restaurant_listing/go_grpc_server/
COPY . .

#Generate stubs
# RUN chmod +x restaurant_proto_gen.sh && ./restaurant_proto_gen.sh
ENV ENV docker
EXPOSE 4000

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o restaurant_server
ENTRYPOINT ["./restaurant_server"]