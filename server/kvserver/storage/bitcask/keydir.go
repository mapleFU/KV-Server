package bitcask

import (
	"strings"
	"hash/fnv"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/mapleFU/data-structures/hashtable"
	"github.com/timtadh/data-structures/types"
	"github.com/jinzhu/copier"
)

type hashKey string

func (key hashKey) Equals(b types.Equatable) bool {
	var other hashKey
	var ok bool
	if other, ok = b.(hashKey); !ok {
		return false
	}
	return strings.Compare(string(other), string(key)) == 0
}

func (key hashKey) Less(b types.Sortable) bool {
	var other hashKey
	var ok bool
	if other, ok = b.(hashKey); !ok {
		log.Fatal("in hashKey.Less, argument is not a hashKey(or string)")
	}
	return strings.Compare(string(key), string(other)) < 0
}

func (key hashKey) Hash() int {
	return int(hash(string(key)))
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

type scanMap struct {
	table *hashtable.Hash
	//// add immutable table
	//immutable *hashtable.Hash	// immutable is a table which will not
	mu sync.RWMutex
}

func (scanMap *scanMap) flushRecord(key string, entry *entry) error {
	scanMap.mu.Lock()
	defer scanMap.mu.Unlock()

	return scanMap.table.Put(hashKey(key), entry)
}

func (scanMap *scanMap) BucketSize() int {
	scanMap.mu.RLock()
	defer scanMap.mu.RUnlock()

	return scanMap.table.BucketSize()
}


func (scanMap *scanMap) addRecord(key string, entry *entry) error {
	scanMap.mu.Lock()
	defer scanMap.mu.Unlock()

	return scanMap.table.Put(hashKey(key), entry)
}

//func (scanMap *scanMap) scan(cursor storage.ScanCursor) ([]*entry, *storage.ScanCursor, error) {
//	panic("implement me")
//}

func (scanMap *scanMap) bucket(cursor int) ([]interface{}, error) {
	scanMap.mu.RLock()
	defer scanMap.mu.RUnlock()

	return scanMap.table.GetBucketIndexes(cursor)
}

func (scanMap *scanMap) deleteRecord(key string) {
	scanMap.mu.Lock()
	defer scanMap.mu.Unlock()

	scanMap.table.Remove(hashKey(key))
}

func (scanMap *scanMap) entry(key string) (*entry, bool) {
	scanMap.mu.RLock()
	defer scanMap.mu.RUnlock()

	ret, err := scanMap.table.Get(hashKey(key))
	if err != nil {
		return nil, false
	}
	retEntry, ok := ret.(*entry)
	return retEntry, ok
}

// iterator is not thread safe!
func (scanMap *scanMap ) iterator() types.KVIterator {
	return scanMap.table.Iterate()
}

// immutableScanMap is not thread safe!
func (keydir *scanMap ) immutableScanMap() *scanMap {
	retMap := new(scanMap)
	copier.Copy(retMap.table, keydir.table)
	return retMap
}

const (
	initSize = 16
)

func newScanMap() *scanMap {
	table := hashtable.NewHashTable(initSize)
	return &scanMap{
		table:table,
	}
}

