package grpc

//go:generate protoc -I . grpc.proto --go_out=plugins=grpc:$GOPATH/src
//go:generate protoc -I . grpc.proto --descriptor_set_out=grpc.protoset --include_imports

//go:generate mockgen -destination mock_grpc/todo.go github.com/gertcuykens/grpc TodoServer,TodoClient
//go:generate mockgen -destination mock_grpc/rg.go github.com/gertcuykens/grpc RouteGuideClient,RouteGuide_RouteChatClient
