package client

import (
	"context"

	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"

	pb "github.com/mapleFU/KV-Server/proto"
	"github.com/mapleFU/KV-Server/proto"
)

type KVClient struct {
	conn *grpc.ClientConn
	c kvstore_methods.KVServicerClient
}


func Open(rpcAddress string) (*KVClient, error) {
	conn, err := grpc.Dial(rpcAddress)
	c := pb.NewKVServicerClient(conn)
	return &KVClient{
		conn:conn,
		c:c,
	}, err
}

func (kvClient *KVClient) Close()  {
	if kvClient == nil {
		log.Fatalln("kvClient in (kvClient *KVClient) Close() is nil")
	}
	kvClient.conn.Close()
}

func (kvClient *KVClient) Get(key string) ([]byte, error) {
	resp, err := kvClient.c.Get(context.Background(), &pb.Key{Key:key})
	return resp.Values, err
}

func (kvClient *KVClient) Set(key string, value []byte) error {
	_, err := kvClient.c.Set(context.Background(), &pb.KVPair{Key:&pb.Key{Key:key}, Value:&pb.Value{Values:value}})
	return err
}

func (kvClient *KVClient) Delete(key string) error {
	_, err := kvClient.c.Delete(context.Background(), &pb.Key{Key:key})
	return err
}

func (kvClient *KVClient) Scan(cursor int, useMatch bool, matchKey string, count int) (int, []string, error) {
 	scanArg := kvstore_methods.ScanArgs{
 		Cursor:int32(cursor),
 		UseKey:useMatch,
 		Match:&pb.Key{Key:matchKey},
 		Count:int32(count),
	}
	retMap := make([]string, 0)
	resp, err := kvClient.c.Scan(context.Background(), &scanArg)
	for _, v := range resp.Values {
		retMap = append(retMap, string(v.Values))
	}
	return int(resp.Cursor), retMap, err
}

