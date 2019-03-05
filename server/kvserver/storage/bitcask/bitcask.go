package bitcask

import (
	"sync"
	"time"
	"regexp"
	"io"

	log "github.com/sirupsen/logrus"

	"github.com/mapleFU/KV-Server/proto"
	"github.com/mapleFU/KV-Server/server/kvserver/storage/bitcask/buffer"
	"github.com/mapleFU/KV-Server/server/kvserver/storage"
)

var emptyBytes []byte

func init()  {
	emptyBytes = make([]byte, 0)
}

type Bitcask struct {
	options *Options
	// map of record

	// change entryMap to scanMap
	//entryMap entryMap
	entryMap scanMap

	bitcaskPoolManager *buffer.BitcaskBufferManager
	redoLogger *redoLogger
	// the directory under control
	directoryName string

	mu sync.RWMutex
}

func (bitcask *Bitcask) currentFileID() int {
	return bitcask.bitcaskPoolManager.CurrentFileId
}

func Open(dirName string, options *Options) *Bitcask  {
	bitcask := Bitcask{}

	bitcask.options = options

	bitcask.entryMap = *newScanMap()
	bitcask.directoryName = dirName
	bc, err := buffer.Open(dirName)
	if err != nil {
		log.Fatal(err)
	}
	bitcask.bitcaskPoolManager = bc

	logger, err := newRedoLogger(bitcask.directoryName)
	if err != nil {
		log.Fatalln(err)
	}
	bitcask.redoLogger = logger

	return &bitcask
}

func (bitcask *Bitcask) Close() {
	bitcask.bitcaskPoolManager.Close()
	bitcask.syncKeyDirToHint()
}


func (bitcask *Bitcask) doWalLog(entry *entry, string string)  {
	if bitcask.options != nil && bitcask.options.UseLog {
		bitcask.redoLogger.writeLog(entry, string)
	}
}

func (bitcask *Bitcask) Scan(cursor storage.ScanCursor) ([][]byte, int, error) {
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
			log.Infof("Match %v with reg %v\n", key, cursor.MatchKeyString)
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
				byteEncodedData, err := bitcask.bitcaskPoolManager.FetchBytes(entry.FileID, entry.ValueSize, entry.ValuePos, entry.Timestamp)
				if err != nil {
					return retBytes, -1, err
				}
				key, valueBytes, _ := buffer.PersistDecoding(byteEncodedData)
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
	entry, exists := bitcask.entryMap.entry(string(key))
	if !exists {
		return emptyBytes, &kvstore_methods.UnexistsError{Key:keyString,}
	}
	byteEncodedData, err := bitcask.bitcaskPoolManager.FetchBytes(entry.FileID, entry.ValueSize, entry.ValuePos, entry.Timestamp)
	if err != nil {
		return emptyBytes, err
	}
	_, valueBytes, _ := buffer.PersistDecoding(byteEncodedData)
	log.Infof("Get %s-%v", string(key), string(valueBytes))
	return valueBytes, nil
}

/**
write disk -- write hashmap
 */
func (bitcask *Bitcask) Put(key []byte, value []byte) error {
	cTime := time.Now()


	log.Infof("Put key(%s)-value(%s)", string(key), string(value))


	dataEntryBytes := buffer.PersistEncoding(key, value, cTime)
	fileID, valueSz, valuePos, timeStamp, err := bitcask.bitcaskPoolManager.AppendRecord(dataEntryBytes)
	if err != nil {
		// put error
		return err
	}

	entry := entry{
		FileID:fileID,
		ValuePos:valuePos,
		ValueSize:valueSz,
		Timestamp:timeStamp,
	}
	bitcask.doWalLog(&entry, string(key))

	if _, ok := bitcask.entryMap.entry(string(key)); ok {
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
	log.Infof("Delete key %s", string(key))
	delValue := emptyBytes
	delEntryBytes := buffer.PersistEncoding(key, delValue, time.Now())
	fileID, valueSz, valuePos, timeStamp, err := bitcask.bitcaskPoolManager.AppendRecord(delEntryBytes)

	if err != nil {
		// put error
		return err
	}
	entry := entry{
		FileID:fileID,
		ValuePos:valuePos,
		ValueSize:valueSz,
		Timestamp:timeStamp,
	}
	bitcask.doWalLog(&entry, string(key))
	bitcask.entryMap.deleteRecord(string(key))
	return nil
}

func (bitcask *Bitcask) Recover() error {
	var bios int64
	fileEof := false
	fileStat, err := bitcask.redoLogger.logFile.Stat()
	if err != nil {
		return err
	}
	fsz := fileStat.Size()
	for !fileEof && bios < fsz {
		log.Infoln(bios, fsz)
		entry, key, biosDelta, err := bitcask.redoLogger.readLog(bios)
		if err != nil {
			if err == io.EOF {
				fileEof = true
			} else {
				return nil
			}
		}
		bios += int64(biosDelta)
		bitcask.entryMap.flushRecord(key, entry)
	}
	return nil
}


