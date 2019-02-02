package buffer

import (
	"testing"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"github.com/mapleFU/KV-Server/server/kvserver/storage"
	"fmt"
	"time"
	"sync"
	"strings"
)



func TestOpen(t *testing.T) {
	testDir := "/Users/fuasahi/GoglandProjects/src/github.com/mapleFU/KV-Server/server/testdata"
	bufPool, err := Open(testDir)
	if err != nil {
		t.Fatal(err)
	}
	defer bufPool.Close()

	log.Info(bufPool.currentFileId)
	if bufPool.currentFileId != 0 {
		t.Fatal("Error, testdir not get 0")
	}
	if _, err := os.Stat(path.Join(testDir, "bitcask.lock")); os.IsExist(err) {
		t.Fatal("Error, bitcask.lock not exists")
	}
	os.Remove(path.Join(testDir, "0.data"))
}

type testEntry struct {
	fileID uint32
	valueSz uint32
	valuePos uint32
	timeStamp uint32
}

func TestBitcaskPoolManager_AppendRecord(t *testing.T) {
	testDir := "/Users/fuasahi/GoglandProjects/src/github.com/mapleFU/KV-Server/server/testdata/testAppend"
	bufPool, err := Open(testDir)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup

	dataArray := make([]testEntry, 100)
	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("key-%d", i))
		value := []byte(fmt.Sprintf("value-%d", i))

		datas := storage.PersistEncoding(key, value, time.Now())
		wg.Add(1)
		go func( dataBytes []byte, index int) {

			v1, v2, v3, v4, err := bufPool.AppendRecord(datas)
			if err != nil {
				dataArray[index] = testEntry{
					v1, v2, v3, v4,
				}
			}
			wg.Done()
		} (datas, i)
	}
	wg.Wait()

	for i := 0; i < 100; i++ {
		if dataArray[i].valueSz != 0 {
			datas, err := bufPool.FetchBytes(dataArray[i].fileID, dataArray[i].valueSz, dataArray[i].valuePos, dataArray[i].timeStamp)
			if err != nil {
				t.Fatal(err)
			}
			key, value, _ :=storage.PersistDecoding(datas)
			if strings.Compare(string(key), fmt.Sprintf("key-%d", i)) != 0 {
				t.Fatal("key error")
			}
			if strings.Compare(string(value), fmt.Sprintf("value-%d", i)) != 0 {
				t.Fatal("key error")
			}
		} else {
			log.Infof("Value in %d pos is zero, write failed", i)
		}
	}
}