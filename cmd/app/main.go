package main

//go:generate protoc -I . app.proto --go_out=plugins=grpc:.
//go:generate protoc -I . app.proto --descriptor_set_out=app.protoset --include_imports
//go:generate mockgen -destination app_mock/app.go -source=app.pb.go -package=app_mock
func main() {}
