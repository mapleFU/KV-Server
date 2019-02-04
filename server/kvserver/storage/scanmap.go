package storage

import (
	"strings"
	"hash/fnv"

	log "github.com/sirupsen/logrus"

	"github.com/mapleFU/data-structures/hashtable"
	"github.com/timtadh/data-structures/types"
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
}

func (scanMap *scanMap) flushRecord(key string, entry *entry) error {
	return scanMap.table.Put(hashKey(key), entry)
}

func (scanMap *scanMap) BucketSize() int {
	return scanMap.table.BucketSize()
}


func (scanMap *scanMap) addRecord(key string, entry *entry) error {
	return scanMap.table.Put(hashKey(key), entry)
}

func (scanMap *scanMap) scan(cursor ScanCursor) ([]*entry, *ScanCursor, error) {
	panic("implement me")
}

func (scanMap *scanMap) bucket(cursor int) ([]interface{}, error) {
	return scanMap.table.GetBucketIndexes(cursor)
}

func (scanMap *scanMap) deleteRecord(key string) {
	scanMap.table.Remove(hashKey(key))
}

func (scanMap *scanMap) getEntry(key string) (*entry, bool) {
	ret, err := scanMap.table.Get(hashKey(key))
	if err != nil {
		return nil, false
	}
	retEntry, ok := ret.(*entry)
	return retEntry, ok
}

func (scanMap *scanMap ) iterator() types.KVIterator {

	return scanMap.table.Iterate()
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

