package main

import (
	"google.golang.org/grpc"
	"log"
	"net"

	pb "proto"
)

const (
	_port = ":50004"
)

func main() {
	lis, err := net.Listen("tcp", _port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening on ", lis.Addr())

	s := grpc.NewServer()
	ins := &server{}
	ins.init()
	pb.RegisterBgSaveServiceServer(s, ins)

	s.Serve(lis)
}
