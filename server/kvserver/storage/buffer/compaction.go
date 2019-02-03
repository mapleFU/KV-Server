package buffer

import (
	"sort"

	log "github.com/sirupsen/logrus"

	"os"
	"sync"
)

func (poolManager *BitcaskPoolManager) compaction()  {
	panic("not implemented.")

	indexes, err := listDataFileID(poolManager.dirName)
	if err != nil {
		log.Fatalln(err)
	}
	sort.Ints(indexes)

 	compactingMap := make(map[string]*Record)
 	mergeRecords := make([][]*Record, len(indexes) - 1)
 	var mgSync sync.Mutex
 	//recordMapChan := make(chan []*Record)

 	var wg sync.WaitGroup
	for index, fileId := range indexes[:len(indexes) - 1] {
		fileName := dataFileName(fileId, poolManager.dirName)
		file, err := poolManager.fetchFilePointer(fileName)
		if err != nil {
			log.Fatalln(err)
		}
		wg.Add(1)
		go func(file *os.File, fIndex int) {

			records, _ := ReadRecords(file)
			//go func() {
			//	recordMapChan <- records
			//}()
			mgSync.Lock()
			mergeRecords[fIndex] = records
			mgSync.Unlock()
			wg.Done()
		}(file, index)

	}
	wg.Wait()
	for _, records := range mergeRecords {
		for _, v := range records {
			compactingMap[string(v.Key)] = v
		}
	}

	poolManager.mu.Lock()
	for _, id := range indexes {
		fileName := dataFileName(id, poolManager.dirName)
		poolManager.closeFilePointerWithoutLock(fileName)
		os.Remove(fileName)
	}
	//f, err := os.Create(dataFileName(0, poolManager.dirName))
	//var bios int64
	//bios = 0
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, v := range compactingMap {
	//	nearestTime := time.Unix(0, int64(v.Header.TimeStamp))
	//	bytesData := schema.PersistEncoding(v.Key, v.Key, nearestTime)
	//
	//}
	panic("not implemented")
	defer poolManager.mu.Unlock()
}
