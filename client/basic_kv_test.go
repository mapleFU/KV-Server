package client

import (
	"google.golang.org/grpc"
	pb "github.com/mapleFU/KV-Server/proto"
	"golang.org/x/net/context"
	"testing"
	"strings"
	"strconv"
	log "github.com/sirupsen/logrus"
)

const (
	rpcAddress = "localhost:50001"
	testKey = "test-key"
)

func TestBasicKVSetAndGet(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(rpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKVServicerClient(conn)
	testValue := []byte(testKey)
	// Contact the server and print out its response.
	//name := defaultName
	//if len(os.Args) > 1 {
	//	name = os.Args[1]
	//}

	// set testkey and get testkey
	_, err = c.Set(context.Background(), &pb.KVPair{Key:&pb.Key{Key:testKey},
											Value:&pb.Value{Values: testValue}})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	getResp, err := c.Get(context.Background(), &pb.Key{Key:testKey})
	if err != nil {
		log.Fatalln(err)
	}
	if strings.Compare(string(getResp.Values), string(testValue)) != 0 {
		t.Fatalf("error, set %s but got %s", testKey, string(getResp.Values))
	}
}

func TestBasicKVGetNil(t* testing.T)  {
	conn, err := grpc.Dial(rpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKVServicerClient(conn)
	_, err = c.Get(context.Background(), &pb.Key{Key:"non-exists"})
	if err == nil {
		// resp 没有意义？
		t.Fatal("Get an unexists value without err.")
	} else {
		log.Infoln(err)
	}
}

func TestBasicKVDelete(t *testing.T)  {
	conn, err := grpc.Dial(rpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKVServicerClient(conn)
	testValue := []byte("non-exists")
	c.Set(context.Background(), &pb.KVPair{Key:&pb.Key{Key:"non-exists"}, Value:&pb.Value{Values:testValue}})
	c.Delete(context.Background(), &pb.Key{Key:"non-exists"})

	_, err = c.Get(context.Background(), &pb.Key{Key:"non-exists"})
	if err == nil {
		t.Fatal("Get an unexists value without err.")
	}
}

func TestBasicKVScan(t *testing.T)  {
	conn, err := grpc.Dial(rpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKVServicerClient(conn)

	for i := 0; i <= 1000; i++ {
		c.Set(context.Background(), &pb.KVPair{
			Key: &pb.Key{Key:"Test" + strconv.Itoa(i)},
			Value: &pb.Value{Values: []byte(strconv.Itoa(i))},
		})
	}

	scan, err := c.Scan(context.Background(), &pb.ScanArgs{
		Match: &pb.Key{Key:"Test.*9$"},
		Count: 1000,
		Cursor: 0,
	})
	for _, v := range  scan.Values {

		i, err := strconv.Atoi(string(v.Values))
		if err != nil {
			t.Fatal(err)
		}
		//fmt.Println(i)
		if i % 10 != 9 {
			t.Fatalf("%d is not match the regex, unpass", i)
		}
	}
	//fmt.Printf("Length %d\n", len(scan.Values))
}