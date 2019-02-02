package kvserver

import (
	//"regexp"
	"context"

	log "github.com/sirupsen/logrus"

	pb "github.com/mapleFU/KV-Server/proto"
	"github.com/mapleFU/KV-Server/server/kvserver/storage"
)

type KVService struct {
	bitcask *storage.Bitcask
}

var EmptyArray []byte

func init()  {
	EmptyArray = make([]byte, 0)

}

func (s *KVService) Get(ctx context.Context, req *pb.Key) (*pb.Value, error) {
	log.Infoln("Get")
	bytes, err := s.bitcask.Get([]byte(req.Key))
	if err != nil {
		log.Infof("Get nil, length of bytes is %d\n", len(bytes))
		bytes = EmptyArray
	}
	return &pb.Value{Values:bytes}, err
}

func (s *KVService) Set(ctx context.Context, req *pb.KVPair) (*pb.OperationOK, error) {
	log.Infoln("Set")
	return &pb.OperationOK{Ok:true}, s.bitcask.Put([]byte(req.Key.Key), req.Value.Values)
}

func (s *KVService) Delete(ctx context.Context, req *pb.Key) (*pb.OperationOK, error)  {
	log.Infoln("Delete")
	err := s.bitcask.Del([]byte(req.Key))
	ok := true
	if err != nil {
		ok = false
	}

	return &pb.OperationOK{
		Ok:ok,
	}, err
}

func (s *KVService) Scan(ctx context.Context, req *pb.ScanArgs) (*pb.Values, error) {
	log.Infoln("Scan")
	panic("impl me")
	//values := make([]*pb.Value, 0)
	//
	//reg, err := regexp.Compile(req.Match.Key)
	//if err != nil {
	//	return &pb.Values{}, nil
	//}
	//index, cnt := 0, 0
	//for k, v := range s.kvMap {
	//
	//	if reg.MatchString(k) {
	//		if int32(index) > req.Cursor && int32(cnt) < req.Count {
	//			values = append(values, &pb.Value{Values:v})
	//			cnt++
	//		}
	//		index++
	//	}
	//}
	//return &pb.Values{
	//	Values:values,
	//}, nil
}

func NewKVService() *KVService {
	return NewKVServiceWithDir("/Users/fuasahi/GoglandProjects/src/github.com/mapleFU/KV-Server/server/data")
}

func NewKVServiceWithDir(dirName string) *KVService {
	bc := storage.Open(dirName)
	//if err != nil {
	//	log.Fatal(err)
	//}
	return &KVService{
		bitcask:bc,
	}
}