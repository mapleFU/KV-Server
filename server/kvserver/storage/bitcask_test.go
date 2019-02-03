package storage

import (
	"testing"
	"fmt"
	"math"
	"bytes"
	"path/filepath"
	"os"
	"path"
)

func TestOpen(t *testing.T) {
	bitcask := Open("/Users/fuasahi/GoglandProjects/src/github.com/mapleFU/KV-Server/server/data", nil)
	defer bitcask.Close()
	if bitcask == nil {
		t.Fatal("bitcask is nil")
	}
	if bitcask.bitcaskPoolManager == nil {
		t.Fatal("bitcask.bitcaskPoolManager is nill")
	}

	//if bitcask.entryMap == nil {
	//	t.Fatal("bitcask.entryMap is nil")
	//}
}

func TestBitcask_Scan(t *testing.T) {
	bitcask := Open("/Users/fuasahi/GoglandProjects/src/github.com/mapleFU/KV-Server/server/testdata/testScan", nil)
	defer bitcask.Close()

	for i := 0; i < 1000; i++ {
		key := []byte(fmt.Sprintf("key-%d", i))
		value := []byte(fmt.Sprintf("value-%d", i))
		err := bitcask.Put(key, value)
		if err != nil {
			t.Fatal(err)
		}
	}

	// not use match key
	var recordCnt int
	ret, cursor, err := bitcask.Scan(ScanCursor{
		UseMatchKey:false,
		Cursor: 0,
		Count: -1,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cursor)
	for _, data := range ret {
		fmt.Println(string(data), len(data))

	}
	recordCnt = len(ret)

	// cursor test
	for cursor != 0 {
		ret, cursor, err = bitcask.Scan(ScanCursor{
			UseMatchKey:false,
			Cursor: cursor,
			Count: -1,
		})
		recordCnt += len(ret)
		fmt.Printf("Add %d records, sum up to %d, now cursor is %d\n", len(ret), recordCnt, cursor)
	}
	if recordCnt != 1000 {
		t.Fatalf("cursor didn't traverse all the map, recordCnt == %d", recordCnt)
	}
}

func TestByteReverse(t *testing.T)  {
	begin := 0
	size := 16
	//if byteReverse(begin, size) != 0 {
	//	t.Fatal("byteReverse(begin, size) != 0")
	//}

	if byteReverse(15, size) != 15 {
		t.Fatal("byteReverse(15, 4) != 15")
	}

	if byteReverse(16, size) != 0 {
		t.Fatal("byteReverse(16, 4) != 0")
	}

	// LOOP
	begin = 0
	for i := 0; i <= int(math.Pow(float64(2), float64(size)) - 1); i++ {
		//fmt.Printf("Now we are from %d\n", begin)
		rev := int(byteReverse(begin, size))
		rev += 1
		begin = int(byteReverse(rev, size))
	}
	if begin != 0 {
		t.Fatal("0000 not change to 0000 during loops")
	}

	cursorSize := 2048
	currentCursor := 0
	isStart := true
	round := 0
	for true  {

		if currentCursor == 0 && isStart == false {
			break
		}
		isStart = false
		rev := int(byteReverse(currentCursor, cursorSize))
		rev += 1
		currentCursor = int(byteReverse(rev, cursorSize))
		fmt.Printf("Get %d\n", currentCursor)
		round++
	}
	fmt.Println(round)
}

func TestBitcask_Recover(t *testing.T) {
	workDir := "/Users/fuasahi/GoglandProjects/src/github.com/mapleFU/KV-Server/server/testdata/testRecover"
	bitcask := Open(workDir, &Options{UseLog:true})
	defer RemoveContents(workDir)

	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("key-%d", i))
		value := []byte(fmt.Sprintf("value-%d", i))
		err := bitcask.Put(key, value)
		if err != nil {
			t.Fatal(err)
		}
	}
	bitcask.Close()

	bitcask = Open(workDir, &Options{UseLog:true})
	err := bitcask.Recover()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := bitcask.Get([]byte("key-1"))
	if err != nil {

		t.Fatal(err)
	}
	if bytes.Compare(resp, []byte("value-1")) != 0 {
		t.Fatal("recover failed")
	}

}

func RemoveContents(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestWalRead(t *testing.T)  {
	workDir := "/Users/fuasahi/GoglandProjects/src/github.com/mapleFU/KV-Server/server/testdata/testLogger"
	defer os.Remove(path.Join(workDir, "log.hint"))
	logger, err := newRedoLogger(workDir)
	if err != nil {
		t.Fatal(err)
	}
	entryLog := &entry{
		FileID:0,
		ValueSize:12,
		ValuePos: 13,
		Timestamp:15,
	}
	logger.writeLog(entryLog, "key-1")
	data := writeWalLog(entryLog, "key-1")
	f, _ := os.Open(path.Join(workDir, "log.hint"))
	readData := make([]byte, len(data))
	f.Read(readData)
	if bytes.Compare(readData, data) != 0 {
		t.Fatal("bytes.Compare(readData, data) != 0")
	}
	f.ReadAt(readData, 0)
	if bytes.Compare(readData, data) != 0 {
		t.Fatal("bytes.Compare(readData, data) != 0")
	}

	entry, key, _, err := readLog(f, 0)
	if entry.ValuePos != entryLog.ValuePos {
		t.Fatal("ReadLog parse error")
	}
	_, key, _, _= logger.readLog(0 )
	if key != "key-1" {
		t.Fatal("logger error with key " + key)
	}

}