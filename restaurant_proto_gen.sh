# protoc --proto_path=./../../protos --go_out=./ --grpc_out=./  --plugin=protoc-gen-grpc=grpc_go_plugin ./../../protos/restaurant.proto
#  protoc -I restaurant/ restaurant/restaurant.proto --go_out=plugins=grpc:restaurant
protoc --proto_path=./restaurant --go_out=plugins=grpc:restaurant restaurant.proto
# protoc -I restaurant --go_out=/restaurant /restaurant/restaurant.proto
# protoc --proto_path=./restaurant --go_out=./restaurant --plugin=grpc_go_plugin restaurant.proto
