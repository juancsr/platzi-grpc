package main

import (
	"log"
	"net"

	"github.com/juancsr/platzi-grpc/database"
	"github.com/juancsr/platzi-grpc/server"
	"github.com/juancsr/platzi-grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	list, err := net.Listen("tcp", "0.0.0.0:5050")
	if err != nil {
		log.Fatal(err)
	}

	repo, err := database.NewPostgresRepository("postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewTestServer(repo)

	s := grpc.NewServer()
	testpb.RegisterTestServiceServer(s, server)

	reflection.Register(s)

	if err := s.Serve(list); err != nil {
		log.Fatal(err)
	}
}
