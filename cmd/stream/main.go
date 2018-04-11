package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/gertcuykens/grpc"
	"github.com/gertcuykens/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func Listen(server pb.RouteGuideServer) {
	creds, err := credentials.NewServerTLSFromFile(tls.CRT, tls.KEY)
	if err != nil {
		log.Fatalf("Failed to generate credentials %s", err)
	}
	srv := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterRouteGuideServer(srv, server)
	reflection.Register(srv)
	l, err := net.Listen("tcp", ":8444")
	if err != nil {
		log.Fatalf("could not listen to :8444 %s", err)
	}
	log.Fatal(srv.Serve(l))
}

func Client() (pb.RouteGuideClient, *grpc.ClientConn) {
	creds, err := credentials.NewClientTLSFromFile(tls.CA, "localhost")
	if err != nil {
		log.Fatalf("Failed to generate credentials %s", err)
	}
	conn, err := grpc.Dial(":8444", grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to backend: %s\n", err)
		os.Exit(1)
	}
	return pb.NewRouteGuideClient(conn), conn
}

func main() {
	go Listen(newServer())
	client, conn := Client()
	defer conn.Close()

	printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})
	printFeature(client, &pb.Point{Latitude: 0, Longitude: 0})
	printFeatures(client, &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	})

	runRecordRoute(client)
	runRouteChat(client)
}
