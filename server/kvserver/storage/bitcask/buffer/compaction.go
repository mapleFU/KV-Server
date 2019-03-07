package buffer

import (
	"sort"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"time"
)

// compaction is a function merge the data of every data file
// it will merge all data file and generate a new hint file
// return bool that shows if we can do the compaction
func (poolManager *BitcaskBufferManager) Compaction() (map[string]*Record, bool, error) {

	indexes, err := listDataFileID(poolManager.dirName)
	if err != nil {
		log.Fatalln(err)
	}
	if len(indexes) <= 2 {
		// not need to merge
		return nil, false, nil
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
			if oldData, e := compactingMap[v]; e {
				if oldData.Header.TimeStamp < v.Header.TimeStamp {
					compactingMap[string(v.Key)] = v
				}
			} else {
				compactingMap[string(v.Key)] = v
			}

		}
	}

	poolManager.mu.Lock()
	return compactingMap, true, nil
}
