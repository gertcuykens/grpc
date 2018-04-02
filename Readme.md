protoc -I . task.proto --go_out=plugins=grpc:.
cat data | protoc --decode_raw
hexdump -c data

mockgen -destination mock2 google.golang.org/grpc/examples/helloworld/helloworld GreeterClient
mockgen -destination mock1 -source=helloworld/helloworld/helloworld.pb.go
