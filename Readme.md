protoc -I . task.proto --go_out=plugins=grpc:.

mockgen -destination mock_grpc/task.go github.com/gertcuykens/grpc TodoServer,TodoClient

mockgen -destination mock_grpc/task.go -source=task.pb.go

cat data | protoc --decode_raw

hexdump -c data
