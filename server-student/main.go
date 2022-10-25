package main

import (
	"fmt"
	"log"
	"net"

	"github.com/juancsr/platzi-grpc/database"
	"github.com/juancsr/platzi-grpc/server"
	"github.com/juancsr/platzi-grpc/studentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fmt.Println("starting module")
	list, err := net.Listen("tcp", ":5060")
	if err != nil {
		log.Fatal(err)
	}

	repo, err := database.NewPostgresRepository("postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewServer(repo)

	s := grpc.NewServer()
	studentpb.RegisterStudentServiceServer(s, server)

	reflection.Register(s)

	fmt.Println("Serving")
	if err := s.Serve(list); err != nil {
		log.Fatal(err)
	}
}
