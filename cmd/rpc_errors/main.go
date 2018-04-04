package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	pb "github.com/gertcuykens/grpc"
	"golang.org/x/net/context"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Listen(server pb.GreeterServer) {
	srv := grpc.NewServer()
	pb.RegisterGreeterServer(srv, server)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("could not listen to :8080 %s", err)
	}
	log.Fatal(srv.Serve(l))
}

func Client() pb.GreeterClient {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to backend: %s\n", err)
		os.Exit(1)
	}
	return pb.NewGreeterClient(conn)
}

type server struct {
	mu    sync.Mutex
	count map[string]int
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count[in.Name]++
	if s.count[in.Name] > 1 {
		st := status.New(codes.ResourceExhausted, "Request limit exceeded.")
		ds, err := st.WithDetails(
			&epb.QuotaFailure{
				Violations: []*epb.QuotaFailure_Violation{{
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
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	go Listen(&server{count: make(map[string]int)})
	client := Client()
	for {
		r, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "world"})
		if err != nil {
			s := status.Convert(err)
			for _, d := range s.Details() {
				switch info := d.(type) {
				case *epb.QuotaFailure:
					log.Printf("Quota failure: %s", info)
				default:
					log.Printf("Unexpected type: %s", info)
				}
			}
			os.Exit(1)
		}
		log.Printf("Greeting: %s", r.Message)
	}
}
