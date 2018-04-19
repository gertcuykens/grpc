package grpc_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"testing"
	"time"

	pb "github.com/gertcuykens/grpc"
	tls2 "github.com/gertcuykens/tls"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/reflection"
)

func Listen2(t *testing.T, server pb.TlsServer) {
	creds, err := credentials.NewServerTLSFromFile(tls2.CRT, tls2.KEY)
	if err != nil {
		t.Fatalf("Failed to generate credentials %s", err)
	}
	srv := grpc.NewServer(grpc.Creds(creds))
	reflection.Register(srv)
	pb.RegisterTlsServer(srv, server)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatalf("could not listen to :8080 %s", err)
	}
	t.Fatal(srv.Serve(l))
}

func Client2(t *testing.T) (pb.TlsClient, *grpc.ClientConn) {
	creds, err := credentials.NewClientTLSFromFile(tls2.CA, "localhost")
	if err != nil {
		t.Fatalf("Failed to generate credentials %s", err)
	}
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(creds), grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	if err != nil {
		t.Fatalf("could not connect to backend: %s\n", err)
	}
	return pb.NewTlsClient(conn), conn
}

func Listen1(t *testing.T, server pb.TlsServer) {
	m := &autocert.Manager{
		Cache:      autocert.DirCache("tls"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("example.com"),
	}
	go http.ListenAndServe(":http", m.HTTPHandler(nil))
	creds := credentials.NewTLS(&tls.Config{GetCertificate: m.GetCertificate})
	srv := grpc.NewServer(grpc.Creds(creds))
	reflection.Register(srv)
	pb.RegisterTlsServer(srv, server)
	l, err := net.Listen("tcp", ":https")
	if err != nil {
		t.Fatalf("could not listen to :https %s", err)
	}
	t.Fatal(srv.Serve(l))
}

func Client1(t *testing.T) (pb.TlsClient, *grpc.ClientConn) {
	pool, err := x509.SystemCertPool()
	creds := credentials.NewClientTLSFromCert(pool, "example.com")
	if err != nil {
		t.Fatalf("Failed to generate credentials %s", err)
	}
	conn, err := grpc.Dial(":https", grpc.WithTransportCredentials(creds), grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	if err != nil {
		t.Fatalf("could not connect to backend: %s\n", err)
	}
	return pb.NewTlsClient(conn), conn
}

type testServer struct{}

func (s testServer) Test(ctx context.Context, err *pb.Err) (*pb.Err, error) {
	return &pb.Err{Msg: "OK"}, nil
}

func TestGRPC(t *testing.T) {
	go Listen2(t, &testServer{})
	client, conn := Client2(t)
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	msg, err := client.Test(ctx, &pb.Err{})
	if err != nil {
		t.Fatalf("%v %v: ", client, err)
	}
	t.Log(msg)
}
