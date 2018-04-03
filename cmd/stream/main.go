package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	pb "github.com/gertcuykens/grpc"
	"github.com/gertcuykens/tls"
	"golang.org/x/net/context"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func Listen(server pb.RouteGuideServer) {
	creds, err := credentials.NewServerTLSFromFile(tls.Path("server.pem"), tls.Path("server-key.pem"))
	if err != nil {
		log.Fatalf("Failed to generate credentials %s", err)
	}
	srv := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterRouteGuideServer(srv, server)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("could not listen to :8080 %s", err)
	}
	log.Fatal(srv.Serve(l))
}

func Client() pb.RouteGuideClient {
	creds, err := credentials.NewClientTLSFromFile(tls.Path("ca.pem"), "localhost")
	if err != nil {
		log.Fatalf("Failed to generate credentials %s", err)
	}
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to backend: %s\n", err)
		os.Exit(1)
	}
	return pb.NewRouteGuideClient(conn)
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

//go:generate protoc -I ../../ ../../grpc.proto --go_out=plugins=grpc:../../
//go:generate mockgen -destination ../../mock_grpc/rg.go github.com/gertcuykens/grpc RouteGuideClient,RouteGuide_RouteChatClient
func main() {
	go Listen(newServer())
	client := Client()
	// defer conn.Close()

	printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})
	printFeature(client, &pb.Point{Latitude: 0, Longitude: 0})
	printFeatures(client, &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	})

	runRecordRoute(client)
	runRouteChat(client)
}
