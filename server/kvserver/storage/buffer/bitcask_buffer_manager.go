package buffer

import (
	"os"
	"sync"
	"github.com/gofrs/flock"
	"path"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type BitcaskPoolManager struct {
	// immutable
	dirName string

	// status
	currentFileId int
	currentFile *os.File
	fileLength int64

	// mutex for pool, the mu is used to operating on
	mu sync.RWMutex
	flock *flock.Flock
}



func (poolManager *BitcaskPoolManager) Close() {
	poolManager.currentFile.Close()
	poolManager.flock.Unlock()
}

func (*BitcaskPoolManager) switchFile()  {

}

// fetch bytes are all read operations, will not effect append.
func (poolManager *BitcaskPoolManager) FetchBytes(fileId uint32, valueSize uint32, valuePos uint32, timeStamp uint32) ([]byte, error) {
	fileName := path.Join(poolManager.dirName, fmt.Sprintf("%d.data", fileId))
	// maybe i should get from a read pool...
	file, err := os.Open(fileName)
	if err != nil {
		// fileID error, must be...
		log.Fatal(err)
	}
	datas := make([]byte, valueSize)
	readData, err := file.ReadAt(datas , int64(valuePos))
	if err != nil {
		return datas, err
	}
	if readData != int(valueSize) {
		log.Infoln("readData not equal to valueSize in FetchBytes")
	}
	return datas, nil
}


func (poolManager *BitcaskPoolManager) AppendRecord(data []byte) (uint32, uint32, uint32, uint32, error) {
	poolManager.mu.Lock()
	defer poolManager.mu.Unlock()
	n, err := poolManager.currentFile.Write(data)

	oldStart := poolManager.fileLength
	poolManager.fileLength += int64(n)

	if err != nil {
		log.Infoln("write failed")
		return 0, 0, 0, 0, err
	}

	return uint32(poolManager.currentFileId), uint32(n), uint32(oldStart), 0, nil
}

func Open(dirName string) (*BitcaskPoolManager, error) {
	fileLock := flock.New(path.Join(dirName, "bitcask.lock"))
	var fileName string
	var currentFileId int
	var currentFile *os.File
	var fileLength int64
	for i := 0; true ; i++ {
		fileName = path.Join(dirName, fmt.Sprintf("%d.data", i))
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			//fileName = f.Name()
			currentFileId = i
			break
		}
	}
	ok, err := fileLock.TryLock()
	if !ok {
		return nil, err
	}

	currentFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	fileLength = 0
	if err != nil {
		log.Fatal(err)
	}

	poolManager := BitcaskPoolManager{
		flock:fileLock,
		currentFileId:currentFileId,
		currentFile: currentFile,
		dirName:dirName,
		fileLength:fileLength,
	}

	return &poolManager, nil
}


