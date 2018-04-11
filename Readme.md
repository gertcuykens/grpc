```go
go generate ./...

go test -v ./...

protoc --decode message grpc.proto < data.pbf > data.txt

protoc --decode_raw < data.pbf > data.txt

hexdump -c data.pbf

grpcurl -cacert /Users/gert/go/src/github.com/gertcuykens/tls/ca.crt localhost:8444 list grpc.RouteGuide

grpcurl -cacert /Users/gert/go/src/github.com/gertcuykens/tls/ca.crt localhost:8444 describe grpc.RouteGuide.ListFeatures

grpcurl -cacert /Users/gert/go/src/github.com/gertcuykens/tls/ca.crt -msg-template localhost:8444 describe grpc.Rectangle

cat msg.json | grpcurl -cacert /Users/gert/go/src/github.com/gertcuykens/tls/ca.crt -d @ localhost:8444 grpc.RouteGuide.ListFeatures
```