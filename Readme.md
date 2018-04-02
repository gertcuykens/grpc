go generate ./...

go test -v ./...

cat data | protoc --decode_raw

hexdump -c data
