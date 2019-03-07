package main

import (
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"

	pb "github.com/mapleFU/KV-Server/proto"
	kv "github.com/mapleFU/KV-Server/server/kvserver"
	"os"
)

const (
	port = ":50001"
)



func main()  {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	os.Mkdir("data", 0777)
	pb.RegisterKVServicerServer(s, kv.NewKVServiceWithDir("data"))

	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Infoln("KV-Server started.")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Infoln("KV-Server end.")
}
