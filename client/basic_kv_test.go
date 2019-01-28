package client

import (
	"google.golang.org/grpc"
	"log"
	pb "github.com/mapleFU/KV-Server/proto"
	"golang.org/x/net/context"
	"testing"
	"strings"
)

func TestBasicKV(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKVServicerClient(conn)

	// Contact the server and print out its response.
	//name := defaultName
	//if len(os.Args) > 1 {
	//	name = os.Args[1]
	//}

	_, err = c.Set(context.Background(), &pb.KVPair{Key:&pb.Key{Key:"hello"},
											Value:&pb.Value{Values: []byte("hello")}})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	getResp, err := c.Get(context.Background(), &pb.Key{Key:"hello"})
	if err != nil {
		log.Fatalln(err)
	}
	if strings.Compare(string(getResp.Values), "hello") != 0 {
		t.Fatalf("error, set %s but got %s", "hello", string(getResp.Values))
	}
}
