package bitcask

import (
	"os"
	"time"
	"path"
	"fmt"

	"github.com/mapleFU/KV-Server/server/kvserver/storage/bitcask/buffer"

	log "github.com/sirupsen/logrus"
)

// Merge is a function to merge data files
func (bc *Bitcask) Merge(directoryName string) error {

	cMap, b, err := bc.bitcaskPoolManager.Compaction()
	if err != nil {
		return err
	}
	if !b {
		return nil
	}

	f, err := os.Create(path.Join(bc.directoryName, fmt.Sprintf("%d.data", 0)))

	var bios int64
	bios = 0
	if err != nil {
		return err
	}


	go func() {
		currentWrite := 0
		for _, v := range cMap {
			nearestTime := time.Unix(0, int64(v.Header.TimeStamp))
			bytesData := buffer.PersistEncoding(v.Key, v.Value, nearestTime)
			n, err := f.Write(bytesData)
			if err != nil {
				log.Info("error, err in f.Write is nil")
				return
			}

			// TODO: map
			if en, e := bc.entryMap.entry(string(v.Key)); e {
				if en.Timestamp <= v.Header.TimeStamp {

					en.FileID = 0
					en.ValueSize = v.Header.ValueSz
					en.ValuePos = uint32(currentWrite)
					bc.entryMap.addRecord(string(v.Key), en)
				}
			}

			currentWrite += n
		}
	}()


	return nil
}

