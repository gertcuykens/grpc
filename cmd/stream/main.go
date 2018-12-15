package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func Listen(server RouteGuideServer) {
	creds, err := credentials.NewServerTLSFromFile("/etc/ssl/certs/tls.pem", "/etc/ssl/private/tls-key.pem")
	if err != nil {
		log.Fatalf("Failed to generate credentials %s", err)
	}
	srv := grpc.NewServer(grpc.Creds(creds))
	RegisterRouteGuideServer(srv, server)
	reflection.Register(srv)
	l, err := net.Listen("tcp", ":8444")
	if err != nil {
		log.Fatalf("could not listen to :8444 %s", err)
	}
	log.Fatal(srv.Serve(l))
}

func Client() (RouteGuideClient, *grpc.ClientConn) {
	creds, err := credentials.NewClientTLSFromFile("/etc/ssl/certs/ca.pem", "localhost")
	if err != nil {
		log.Fatalf("Failed to generate credentials %s", err)
	}
	conn, err := grpc.Dial(":8444", grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to backend: %s\n", err)
		os.Exit(1)
	}
	return NewRouteGuideClient(conn), conn
}

//go:generate protoc -I . stream.proto --go_out=plugins=grpc:.
//go:generate protoc -I . stream.proto --descriptor_set_out=stream.protoset --include_imports
//go:generate mockgen -destination stream.pb_test.go -source=stream.pb.go -package=main
func main() {
	go Listen(newServer())
	client, conn := Client()
	defer conn.Close()

	printFeature(client, &Point{Latitude: 409146138, Longitude: -746188906})
	printFeature(client, &Point{Latitude: 0, Longitude: 0})
	printFeatures(client, &Rectangle{
		Lo: &Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &Point{Latitude: 420000000, Longitude: -730000000},
	})

	runRecordRoute(client)
	runRouteChat(client)
}
