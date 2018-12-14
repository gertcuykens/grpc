package main

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/reflection"
)

func Listen1(t *testing.T, server AppServer) {
	creds, err := credentials.NewServerTLSFromFile("/etc/ssl/certs/tls.pem", "/etc/ssl/private/tls-key.pem")
	if err != nil {
		t.Fatalf("Failed to generate credentials %s", err)
	}
	srv := grpc.NewServer(grpc.Creds(creds))
	reflection.Register(srv)
	RegisterAppServer(srv, server)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatalf("could not listen to :8080 %s", err)
	}
	t.Fatal(srv.Serve(l))
}

func Client1(t *testing.T) (AppClient, *grpc.ClientConn) {
	creds, err := credentials.NewClientTLSFromFile("/etc/ssl/certs/ca.pem", "localhost")
	if err != nil {
		t.Fatalf("Failed to generate credentials %s", err)
	}
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(creds), grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	if err != nil {
		t.Fatalf("could not connect to backend: %s\n", err)
	}
	return NewAppClient(conn), conn
}

type testServer1 struct{}

func (s testServer1) Test(ctx context.Context, err *Err) (*Err, error) {
	return &Err{Msg: "OK"}, nil
}

func TestGRPC1(t *testing.T) {
	go Listen1(t, &testServer1{})
	client, conn := Client1(t)
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	msg, err := client.Test(ctx, &Err{})
	if err != nil {
		t.Fatalf("%v %v: ", client, err)
	}
	t.Log(msg)
}
