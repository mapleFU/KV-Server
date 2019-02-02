package storage

import (
	"os"
	"sync"
)

type MyBitcask struct {
	// map of record
	entryMap entryMap

	// the directory under control
	directoryName string
	// current file name
	currentFileName string
	currentFile *os.File


	mu sync.RWMutex
}

func OpenBitcask(dirName string) *MyBitcask  {
	panic("implement me")
	return nil
}

func (*MyBitcask) Close() {
	panic("implement me")
}

func (*MyBitcask) Scan(cursor ScanCursor) ([][]byte, error) {
	panic("implement me")
}

func (bitcask *MyBitcask) Get(key []byte) ([]byte, error) {
	//entry, exists := bitcask.entryMap.getEntry(string(key))
	//if exists {
	//
	//}
	panic("impl me!")
	return nil, nil
}

/**
write disk -- write hashmap
 */
func (*MyBitcask) Put([]byte, []byte) error {
	panic("imple me")
}

/**
like Put...
write disk -- write hashmap
 */
func (*MyBitcask) Del([]byte) error {
	panic("implement me")
}



