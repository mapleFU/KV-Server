package main

import (
	"net"
	log "github.com/sirupsen/logrus"
	pb "github.com/mapleFU/KV-Server/proto"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	"context"
)

const (
	port = ":50001"
)

type KVService struct {
	kvMap map[string][]byte
}

func (s *KVService) Get(ctx context.Context, req *pb.Key) (*pb.Value, error) {
	log.Infof("Call Get with Key(%s)", req.Key)
	if values, ok := s.kvMap[req.Key]; ok {
		return &pb.Value{
			Values:values,
		}, nil
	} else {
		return nil, nil
	}
}

func (s *KVService) Set(ctx context.Context, req *pb.KVPair) (*pb.OperationOK, error) {
	log.Infof("Call Get with KVPair(%s, %s)", req.Key.Key, string(req.Value.Values))
	s.kvMap[req.Key.Key] = req.Value.Values
	return &pb.OperationOK{Ok:true}, nil
}

func (s *KVService) Delete(ctx context.Context, req *pb.Key) (*pb.OperationOK, error)  {
	delete(s.kvMap, req.Key)
	return &pb.OperationOK{Ok:true}, nil
}

func (s *KVService) Scan(ctx context.Context, req *pb.ScanArgs) (*pb.Values, error) {
	return nil, nil
}

func NewKVService() *KVService {
	kvMap := make(map[string][]byte)
	return &KVService{
		kvMap:kvMap,
	}
}



func main()  {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.

	pb.RegisterKVServicerServer(s, NewKVService())

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
