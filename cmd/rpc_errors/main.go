package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	pb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Listen(server GreeterServer) {
	srv := grpc.NewServer()
	RegisterGreeterServer(srv, server)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("could not listen to :8080 %s", err)
	}
	log.Fatal(srv.Serve(l))
}

func Client() GreeterClient {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to backend: %s\n", err)
		os.Exit(1)
	}
	return NewGreeterClient(conn)
}

type server struct {
	mu    sync.Mutex
	count map[string]int
}

func (s *server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("meta: %v", md)
	}

	header := metadata.Pairs("Authorization", "Bearer <token server header>")
	grpc.SendHeader(ctx, header)
	trailer := metadata.Pairs("Authorization", "Bearer <token server trailer>")
	grpc.SetTrailer(ctx, trailer)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.count[in.Name]++
	if s.count[in.Name] > 1 {
		st := status.New(codes.ResourceExhausted, "Request limit exceeded.")
		ds, err := st.WithDetails(
			&pb.QuotaFailure{
				Violations: []*pb.QuotaFailure_Violation{{
					Subject:     fmt.Sprintf("name:%s", in.Name),
					Description: "Limit one greeting per person",
				}},
			},
		)
		if err != nil {
			return nil, st.Err()
		}
		return nil, ds.Err()
	}
	return &HelloReply{Message: "Hello " + in.Name}, nil
}

//go:generate protoc -I . rpc_errors.proto --go_out=plugins=grpc:.
//go:generate protoc -I . rpc_errors.proto --descriptor_set_out=rpc_errors.protoset --include_imports
//go:generate mockgen -destination rpc_errors.pb_test.go -source=rpc_errors.pb.go -package=main
func main() {
	go Listen(&server{count: make(map[string]int)})
	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", "Bearer <token client>")
	var header, trailer metadata.MD
	client := Client()
	for {
		r, err := client.SayHello(ctx, &HelloRequest{Name: "world"}, grpc.Header(&header), grpc.Trailer(&trailer))
		if err != nil {
			s := status.Convert(err)
			for _, d := range s.Details() {
				switch info := d.(type) {
				case *pb.QuotaFailure:
					log.Printf("Quota failure: %s", info)
				default:
					log.Printf("Unexpected type: %s", info)
				}
			}
			os.Exit(1)
		}
		log.Printf("Greeting: %s", r.Message)
		log.Printf("Header: %s", header.Get("Authorization"))
		log.Printf("Trailer: %s", trailer.Get("Authorization"))
	}
}
