package main

import (
	"net"
	"os"
)

import (
	"google.golang.org/grpc"
)

import (
	log "github.com/pengswift/gamelibs/nsq-logger"
	_ "github.com/pengswift/gamelibs/statsd-pprof"
)

import (
	pb "proto"
)

const (
	_port = ":50004"
)

func main() {
	log.SetPrefix(SERVICE)

	lis, err := net.Listen("tcp", _port)
	if err != nil {
		log.Critical(err)
		os.Exit(-1)
	}
	log.Info("listening on ", lis.Addr())

	s := grpc.NewServer()
	ins := &server{}
	ins.init()
	pb.RegisterBgSaveServiceServer(s, ins)

	s.Serve(lis)
}
