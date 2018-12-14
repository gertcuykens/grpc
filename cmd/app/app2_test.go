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
	"google.golang.org/grpc/reflection"
)

func Listen2(t *testing.T, server AppServer) {
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
	RegisterAppServer(srv, server)
	l, err := net.Listen("tcp", ":https")
	if err != nil {
		t.Fatalf("could not listen to :https %s", err)
	}
	t.Fatal(srv.Serve(l))
}

func Client2(t *testing.T) (AppClient, *grpc.ClientConn) {
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
	return NewAppClient(conn), conn
}

type testServer2 struct{}

func (s testServer2) Test(ctx context.Context, err *Err) (*Err, error) {
	return &Err{Msg: "OK"}, nil
}

func TestGRPC2(t *testing.T) {
	go Listen2(t, &testServer2{})
	client, conn := Client2(t)
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	msg, err := client.Test(ctx, &Err{})
	if err != nil {
		t.Fatalf("%v %v: ", client, err)
	}
	t.Log(msg)
}
