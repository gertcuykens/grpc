protoc -I . task.proto --go_out=plugins=grpc:.
cat data | protoc --decode_raw
hexdump -c data
