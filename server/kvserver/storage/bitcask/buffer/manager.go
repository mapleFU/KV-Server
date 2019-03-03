package buffer

import (
	"os"
	"sync"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/gofrs/flock"
)

type BitcaskPoolManager struct {
	// immutable
	dirName string

	// status
	CurrentFileId int
	currentFile *os.File
	fileLength int64

	mu sync.RWMutex
	filePool map[string]*os.File	// pool to use map
	// mutex for pool, the mu is used to operating on

	appendMu sync.Mutex			//the mutex to append data
	flock *flock.Flock
}

func (poolManager *BitcaskPoolManager) getCurrentFile() *os.File {
	f, err := poolManager.fetchFilePointer(dataFileName(poolManager.CurrentFileId, poolManager.dirName))
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func (poolManager *BitcaskPoolManager) closeFilePointerWithoutLock(fileName string)  {
	file, exists := poolManager.filePool[fileName]
	if !exists {
		return
	}
	file.Close()
}

func (poolManager *BitcaskPoolManager) closeFilePointer(fileName string)  {
	poolManager.mu.Lock()
	poolManager.closeFilePointerWithoutLock(fileName)
	poolManager.mu.Unlock()
}

func(poolManager *BitcaskPoolManager) closeAllFilePointer()  {
	for _, v := range poolManager.filePool {
		v.Close()
	}
}

func (poolManager *BitcaskPoolManager) fetchFilePointer(fileName string) (*os.File, error) {
	poolManager.mu.RLock()
	if file, exists := poolManager.filePool[fileName]; exists {
		defer poolManager.mu.RUnlock()
		return file, nil
	} else {
		poolManager.mu.RUnlock()
		poolManager.mu.Lock()

		defer poolManager.mu.Unlock()
		if file, exists := poolManager.filePool[fileName]; exists {

			return file, nil
		} else {
			file, err := os.Open(fileName)
			if err != nil {
				return nil, err
			}
			return file, nil
		}
	}
}


func (poolManager *BitcaskPoolManager) Close() {
	poolManager.closeAllFilePointer()
	poolManager.currentFile.Close()
	poolManager.flock.Unlock()
}

func (*BitcaskPoolManager) switchFile()  {

}

// fetch bytes are all read operations, will not effect append.
func (poolManager *BitcaskPoolManager) FetchBytes(fileId uint32, valueSize uint32, valuePos uint32, timeStamp uint32) ([]byte, error) {

	fileName := dataFileName(int(fileId), poolManager.dirName)
	// maybe i should get from a read pool...
	file, err := poolManager.fetchFilePointer(fileName)
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
	poolManager.appendMu.Lock()
	defer poolManager.appendMu.Unlock()
	n, err := poolManager.currentFile.Write(data)

	oldStart := poolManager.fileLength
	poolManager.fileLength += int64(n)

	if err != nil {
		log.Infoln("write failed")
		return 0, 0, 0, 0, err
	}

	return uint32(poolManager.CurrentFileId), uint32(n), uint32(oldStart), 0, nil
}



func Open(dirName string) (*BitcaskPoolManager, error) {
	fileLock := flock.New(path.Join(dirName, "bitcask.lock"))

	var fileName string
	var currentFileId int
	var currentFile *os.File
	var fileLength int64

	currentFileId = 0
	arr, err := listDataFileID(dirName)
	if err != nil {
		return nil, err
	}
	for _, id := range arr {
		if id >= currentFileId {
			currentFileId = id + 1
		}
	}
	fileName = dataFileName(currentFileId, dirName)

	ok, err := fileLock.TryLock()
	if !ok {
		return nil, err
	}
	//
	currentFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	fileLength = 0
	if err != nil {
		log.Fatal(err)
	}

	poolManager := BitcaskPoolManager{
		flock:fileLock,
		CurrentFileId:currentFileId,
		currentFile: currentFile,
		dirName:dirName,
		fileLength:fileLength,
		filePool:make(map[string]*os.File),
	}

	return &poolManager, nil
}


