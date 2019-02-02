package buffer

import (
	"os"
	"sync"
	"github.com/gofrs/flock"
	"path"
	"fmt"
)

type BitcaskPoolManager struct {
	currentFileId int
	currentFile *os.File


	// mutex for pool, the mu is used to operating on
	mu sync.RWMutex
	flock *flock.Flock
}

func (*BitcaskPoolManager) switchFile()  {

}

func (*BitcaskPoolManager) FetchBytes(fileId uint32, valueSize uint32,
	valuePos uint32, timeStamp uint32)([]byte, error) {
	panic("implement me")
}

func (*BitcaskPoolManager) AppendRecord([]byte) (uint32, uint32, uint32, uint32) {
	panic("implement me")
}

func Open(dirName string) *BitcaskPoolManager {
	fileLock := flock.New(path.Join(dirName, "bitcask.lock"))
	var fileName string
	var currentFileId int
	var currentFile *os.File
	for i := 0; true ; i++ {
		fileName = path.Join(dirName, fmt.Sprintf("%d.data", i))
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			//fileName = f.Name()
			currentFileId = i
			break
		}
	}

	currentFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	poolManager := BitcaskPoolManager{
		flock:fileLock,
		currentFileId:currentFileId,
		currentFile: currentFile,
	}

	return &poolManager
}


