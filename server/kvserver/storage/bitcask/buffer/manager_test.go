package buffer

import (
	"testing"

	"os"
	"path"
	"sync"
	"time"
	"strings"
	"fmt"

	log "github.com/sirupsen/logrus"
)



func TestOpen(t *testing.T) {
	testDir := "testdata/test-open"

	err := os.Mkdir(testDir, 0777)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	bufPool, err := Open(testDir, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer bufPool.Close()

	log.Info(bufPool.CurrentFileId)
	if bufPool.CurrentFileId != 0 {
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
	testDir := "testdata/testAppend"


	err := os.Mkdir(testDir, 0777)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	bufPool, err := Open(testDir, nil)
	defer bufPool.Close()
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup

	var arrMux sync.Mutex
	dataArray := make([]testEntry, 100)
	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("key-%d", i))
		value := []byte(fmt.Sprintf("value-%d", i))

		datas := PersistEncoding(key, value, time.Now())
		wg.Add(1)
		go func( dataBytes []byte, index int) {

			v1, v2, v3, v4, err := bufPool.AppendRecord(datas)
			if err != nil {
				t.Fatal(err)

			} else {
				arrMux.Lock()
				defer arrMux.Unlock()
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
			key, value, _ := PersistDecoding(datas)
			if strings.Compare(string(key), fmt.Sprintf("key-%d", i)) != 0 {
				t.Fatal("key error")
			}
			if strings.Compare(string(value), fmt.Sprintf("value-%d", i)) != 0 {
				t.Fatal("value error")
			}
		} else {
			log.Infof("Value in %d pos is zero, write failed", i)
		}
	}


}

func TestReadRecords(t *testing.T) {
	testDir := "testdata/testAppend"

	os.Mkdir(testDir, 0777)
	defer os.RemoveAll(testDir)

	bufPool, err := Open(testDir, nil)
	defer bufPool.Close()
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup

	var arrMux sync.Mutex
	dataArray := make([]testEntry, 10)
	for i := 0; i < 10; i++ {
		key := []byte(fmt.Sprintf("key-%d", i))
		value := []byte(fmt.Sprintf("value-%d", i))

		datas := PersistEncoding(key, value, time.Now())
		wg.Add(1)
		go func( dataBytes []byte, index int) {

			v1, v2, v3, v4, err := bufPool.AppendRecord(datas)
			if err != nil {
				t.Fatal(err)

			} else {
				arrMux.Lock()
				defer arrMux.Unlock()
				dataArray[index] = testEntry{
					v1, v2, v3, v4,
				}
			}
			wg.Done()
		} (datas, i)
	}
	wg.Wait()
	f, err := os.Open("testdata/testAppend/0.data")
	if err != nil {
		t.Fatal(err)
	}
	records, err := ReadRecords(f)
	if err != nil {
		t.Fatal(err)
	}
	os.Remove("testdata/testAppend/0.data")
	log.Infoln(len(records))
}

