package kvserver

import (
	"regexp"
	"context"

	log "github.com/sirupsen/logrus"

	pb "github.com/mapleFU/KV-Server/proto"
)

type KVService struct {
	kvMap map[string][]byte
}

var EmptyArray []byte

func init()  {
	EmptyArray = make([]byte, 0)
}

func (s *KVService) Get(ctx context.Context, req *pb.Key) (*pb.Value, error) {
	log.Infof("Call Get with Key(%s)", req.Key)
	if values, ok := s.kvMap[req.Key]; ok {
		return &pb.Value{
			Values:values,
		}, nil
	} else {
		//log.Println("Send " + string(EmptyArray))
		return &pb.Value{
			Values: EmptyArray,
		}, nil
	}
}

func (s *KVService) Set(ctx context.Context, req *pb.KVPair) (*pb.OperationOK, error) {
	log.Infof("Call Get with KVPair(%s, %s)", req.Key.Key, string(req.Value.Values))
	s.kvMap[req.Key.Key] = req.Value.Values
	return &pb.OperationOK{Ok:true}, nil
}

func (s *KVService) Delete(ctx context.Context, req *pb.Key) (*pb.OperationOK, error)  {
	var ok bool
	if _, ok = s.kvMap[req.Key]; ok {
		delete(s.kvMap, req.Key)
	}
	return &pb.OperationOK{Ok:ok}, nil
}

func (s *KVService) Scan(ctx context.Context, req *pb.ScanArgs) (*pb.Values, error) {

	values := make([]*pb.Value, 0)

	reg, err := regexp.Compile(req.Match.Key)
	if err != nil {
		return &pb.Values{}, nil
	}
	index, cnt := 0, 0
	for k, v := range s.kvMap {

		if reg.MatchString(k) {
			if int32(index) > req.Cursor && int32(cnt) < req.Count {
				values = append(values, &pb.Value{Values:v})
				cnt++
			}
			index++
		}
	}
	return &pb.Values{
		Values:values,
	}, nil
}

func NewKVService() *KVService {
	kvMap := make(map[string][]byte)
	return &KVService{
		kvMap:kvMap,
	}
}