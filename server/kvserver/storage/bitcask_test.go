package storage

import (
	"testing"
	"fmt"
)

func TestOpen(t *testing.T) {
	bitcask := Open("/Users/fuasahi/GoglandProjects/src/github.com/mapleFU/KV-Server/server/data")
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
	bitcask := Open("/Users/fuasahi/GoglandProjects/src/github.com/mapleFU/KV-Server/server/testdata/testScan")
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
}