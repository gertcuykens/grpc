package main

//go:generate protoc -I . app.proto --go_out=plugins=grpc:.
//go:generate protoc -I . app.proto --descriptor_set_out=app.protoset --include_imports
//go:generate mockgen -destination app.pb_test.go -source=app.pb.go -package=main
func main() {}
