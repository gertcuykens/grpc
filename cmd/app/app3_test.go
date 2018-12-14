package main

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func Listen3(t *testing.T, server health.HealthServer) {
	t.Logf("Starting %s (%s) server...", os.Getenv("HOSTNAME"), os.Getenv("EMAIL"))
	m := &autocert.Manager{
		Cache:      autocert.DirCache("tls"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(os.Getenv("HOSTNAME")),
		Email:      os.Getenv("EMAIL"),
	}
	go http.ListenAndServe(":http", m.HTTPHandler(nil))
	creds := credentials.NewTLS(&tls.Config{GetCertificate: m.GetCertificate})
	srv := grpc.NewServer(grpc.Creds(creds))
	reflection.Register(srv)
	health.RegisterHealthServer(srv, server)
	l, err := net.Listen("tcp", ":https")
	if err != nil {
		t.Fatalf("could not listen to :https %s", err)
	}
	t.Fatal(srv.Serve(l))
}

func Client3(t *testing.T) (health.HealthClient, *grpc.ClientConn) {
	// pool, err := x509.SystemCertPool()
	// creds := credentials.NewClientTLSFromCert(pool, os.Getenv("HOSTNAME"))
	creds, err := credentials.NewClientTLSFromFile(path.Join("tls", os.Getenv("HOSTNAME")), os.Getenv("HOSTNAME"))
	if err != nil {
		t.Fatalf("Failed to generate credentials %s", err)
	}
	conn, err := grpc.Dial(":https", grpc.WithTransportCredentials(creds), grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	if err != nil {
		t.Fatalf("could not connect to backend: %s\n", err)
	}
	return health.NewHealthClient(conn), conn
}

type healthServer struct{}

func (s healthServer) Check(ctx context.Context, r *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}, nil
}

func (s healthServer) Watch(r *health.HealthCheckRequest, w health.Health_WatchServer) error {
	return nil
}

func TestGRPC3(t *testing.T) {
	go Listen3(t, &healthServer{})
	client, conn := Client3(t)
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	msg, err := client.Check(ctx, &health.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("%v %v: ", client, err)
	}
	t.Log(msg)
}
