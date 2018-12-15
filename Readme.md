https://github.com/protocolbuffers/protobuf/releases

go get -u github.com/golang/protobuf/protoc-gen-go

hexdump -c data.pbf

protoc -I . grpc.proto --decode grpc.Message < data.pbf > data.txt

protoc -I . grpc.proto --encode grpc.Message < data.txt > data.pbf

protoc --decode_raw < data.pbf > data.txt

grpcurl -cacert /etc/ssl/certs/ca.pem localhost:8444 list main.RouteGuide

grpcurl -cacert /etc/ssl/certs/ca.pem localhost:8444 describe main.RouteGuide.ListFeatures

grpcurl -cacert /etc/ssl/certs/ca.pem -msg-template localhost:8444 describe main.Rectangle

cat msg.json | grpcurl -cacert /etc/ssl/certs/ca.pem -d @ localhost:8444 main.RouteGuide.ListFeatures

grpcurl -protoset stream.protoset list main.RouteGuide
