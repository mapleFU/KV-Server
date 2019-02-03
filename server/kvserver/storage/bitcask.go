package storage

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/mapleFU/KV-Server/server/kvserver/storage/buffer"
	"time"
	"github.com/mapleFU/KV-Server/proto"
	"regexp"

)

var emptyBytes []byte

func init()  {
	emptyBytes = make([]byte, 0)
}

type Bitcask struct {
	// map of record

	// change entryMap to scanMap
	//entryMap entryMap
	entryMap scanMap

	bitcaskPoolManager *buffer.BitcaskPoolManager

	// the directory under control
	directoryName string

	mu sync.RWMutex
}

func Open(dirName string) *Bitcask  {
	bitcask := Bitcask{}

	bitcask.entryMap = *newScanMap()
	bitcask.directoryName = dirName
	bc, err := buffer.Open(dirName)
	if err != nil {
		log.Fatal(err)
	}
	bitcask.bitcaskPoolManager = bc

	return &bitcask
}

func (bitcask *Bitcask) Close() {
	bitcask.bitcaskPoolManager.Close()
}

func (bitcask *Bitcask) Scan(cursor ScanCursor) ([][]byte, int, error) {
	var judgeStr func(string string) bool
	retBytes := make([][]byte, 0)

	if cursor.Cursor < 0 {
		return retBytes, -1, &ArgumentError{
			Expected:"not less than 0",
			Value:cursor.Cursor,
		}
	}
	if cursor.Count <= 0 {
		cursor.Count = 10
	}

	if cursor.UseMatchKey {
		reg, err := regexp.Compile(cursor.MatchKeyString)
		if err != nil {
			return retBytes, -1, err
		}
		judgeStr = func(key string) bool {
			return reg.Match([]byte(key))
		}
	} else {
		judgeStr = func(string string) bool {
			return true
		}
	}

	currentCursor := cursor.Cursor
	cursorSize := bitcask.entryMap.BucketSize()
	var lastAdd int = 0
	//buckets := bitcask.entryMap.table.ListBuckets()
	//for k, v := range buckets {
	//	fmt.Printf("number %d(%d)\n", k, v)
	//}
	//fmt.Printf("Size %d\n", cursorSize)
	for i := 0; i + lastAdd < cursor.Count;  {

		lastAdd = 0

		//if buckets[currentCursor] == 1 {
		//	fmt.Println("Attention here")
		//}
		values, err := bitcask.entryMap.bucket(currentCursor)

		//fmt.Printf("Use cursor %d get length %d, and bucket is %d, i is %d\n", currentCursor, len(values), buckets[currentCursor], i)
		if err == nil {
			for _, k := range values {
				entry := k.(*entry)
				byteEncodedData, err := bitcask.bitcaskPoolManager.FetchBytes(entry.fileID, entry.valueSize, entry.valuePos, entry.timestamp)
				if err != nil {
					return retBytes, -1, err
				}
				key, valueBytes, _ := PersistDecoding(byteEncodedData)
				// match the string
				if !judgeStr(string(key)) {
					continue
				}
				//fmt.Println("Append")
				retBytes = append(retBytes, valueBytes)
				i++
				lastAdd++
			}
		}

		// next cursor

		revcur := byteReverse(currentCursor, cursorSize)
		revcur++
		// next cursor
		currentCursor = int(byteReverse(int(revcur), cursorSize))
		if currentCursor == 0 {
			break
		}

	}
	return retBytes, currentCursor, nil
}

func countDataBytes( size uint32) int {
	cnt := 0
	for size > 0 {
		size = size >> 1
		cnt++
	}
	return cnt - 1
}

func byteReverse(num int, size int) uint32 {
	countData := countDataBytes(uint32(size))
	var ret uint32 = 0
	for  i := 0; i < countData; i++ {
		ret *= 2
		if num % 2 == 1 {
			ret += 1
		}


		num = num >> 1
	}
	return ret
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
	delValue := emptyBytes
	delEntryBytes := PersistEncoding(key, delValue, time.Now())
	_, _, _, _, err := bitcask.bitcaskPoolManager.AppendRecord(delEntryBytes)
	if err != nil {
		// put error
		return err
	}
	bitcask.entryMap.deleteRecord(string(key))
	return nil
}



