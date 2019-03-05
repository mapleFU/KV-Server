package bitcask

import (
	"testing"
	"github.com/gpmgo/gopm/modules/log"
	"fmt"
)

func TestScanMap_iterator(t *testing.T)  {
	var sm *scanMap
	sm = newScanMap()

	entryArr := make([]*entry, 0)
	for i := 0; i < 5; i++ {
		u32i := uint32(i)
		centry := createEntryInCurrentTime(u32i, u32i, u32i)
		entryArr = append(entryArr, centry)
		sm.addRecord(fmt.Sprintf("key-%d", i), centry)
	}
	iter := sm.iterator()
	cnt := 0
	for k, v, next := iter(); next != nil; k, v, next = next() {
		log.Info("k-%v, v-%v", k, v)
		cnt ++
	}
	if cnt != 5 {
		t.Fatal("add 5 records but cnt != 5, failed")
	}

}