package storage

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/mapleFU/KV-Server/server/kvserver/storage/buffer"
	"time"
	"github.com/mapleFU/KV-Server/proto"
)

var emptyBytes []byte

func init()  {
	emptyBytes = make([]byte, 0)
}

type Bitcask struct {
	// map of record
	entryMap entryMap
	bitcaskPoolManager *buffer.BitcaskPoolManager

	// the directory under control
	directoryName string

	mu sync.RWMutex
}

func Open(dirName string) *Bitcask  {
	bitcask := Bitcask{}
	bitcask.directoryName = dirName
	bc, err := buffer.Open(dirName)
	if err != nil {
		log.Fatal(err)
	}
	bitcask.bitcaskPoolManager = bc

	return nil
}

func (bitcask *Bitcask) Close() {
	bitcask.bitcaskPoolManager.Close()
}

func (*Bitcask) Scan(cursor ScanCursor) ([][]byte, error) {
	panic("implement me")
}

func (bitcask *Bitcask) Get(key []byte) ([]byte, error) {
	keyString := string(key)
	entry, exists := bitcask.entryMap.getEntry(string(key))
	if !exists {
		return emptyBytes, &kvstore_methods.UnexistsError{Key:keyString,}
	}
	byteEncodedData, err := bitcask.bitcaskPoolManager.FetchBytes(entry.fileID, entry.valueSize, entry.valuePos, entry.timestamp)
	if err != nil {
		return emptyBytes, err
	}
	_, valueBytes, _ := PersistDecoding(byteEncodedData)
	return valueBytes, nil
}

/**
write disk -- write hashmap
 */
func (bitcask *Bitcask) Put(key []byte, value []byte) error {
	cTime := time.Now()
	log.Infof("Put key(%s)-value(%s)", string(key), string(value))
	dataEntryBytes := PersistEncoding(key, value, cTime)
	fileID, valueSz, valuePos, timeStamp, err := bitcask.bitcaskPoolManager.AppendRecord(dataEntryBytes)
	if err != nil {
		// put error
		return err
	}

	entry := entry{
		fileID:fileID,
		valuePos:valuePos,
		valueSize:valueSz,
		timestamp:timeStamp,
	}
	if _, ok := bitcask.entryMap.getEntry(string(key)); ok {
		// flush
		err = bitcask.entryMap.flushRecord(string(key), &entry)
		if err != nil {
			return err
		}
		return nil
	} else {
		// is it true that...we use string well?
		err = bitcask.entryMap.addRecord(string(key), &entry)
		if err != nil {
			return err
		}
		return nil
	}
}

/**
like Put...
 write hashmap, just remove the record
 */
func (bitcask *Bitcask) Del(key []byte) error {
	bitcask.entryMap.deleteRecord(string(key))
	return nil
}



