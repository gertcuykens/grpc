```go
go generate ./...

go test -v ./...

protoc --decode message grpc.proto < data.pbf > data.txt

protoc --decode_raw < data.pbf > data.txt

hexdump -c data.pbf
```