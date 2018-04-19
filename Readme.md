```go
go generate ./...

go test -v ./...

hexdump -c data.pbf

protoc -I . grpc.proto --decode grpc.Message < data.pbf > data.txt

protoc -I . grpc.proto --encode grpc.Message < data.txt > data.pbf

protoc --decode_raw < data.pbf > data.txt

grpcurl -cacert /Users/gert/go/src/github.com/gertcuykens/tls/ca.crt localhost:8444 list grpc.RouteGuide

grpcurl -cacert /Users/gert/go/src/github.com/gertcuykens/tls/ca.crt localhost:8444 describe grpc.RouteGuide.ListFeatures

grpcurl -cacert /Users/gert/go/src/github.com/gertcuykens/tls/ca.crt -msg-template localhost:8444 describe grpc.Rectangle

cat msg.json | grpcurl -cacert /Users/gert/go/src/github.com/gertcuykens/tls/ca.crt -d @ localhost:8444 grpc.RouteGuide.ListFeatures

grpcurl -protoset grpc.protoset list grpc.RouteGuide
```