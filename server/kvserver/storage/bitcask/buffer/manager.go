package buffer

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/gofrs/flock"
	"github.com/mapleFU/KV-Server/server/kvserver/storage/bitcask/options"
)

type BitcaskBufferManager struct {
	// immutable
	dirName string	// name of directory of buffer manager

	// status
	CurrentFileId int	// current file id, manages the current r/w on the file
	CurrentFileName string
	currentWriter BufWriter	// use BufWriter to replace os.File, thus it allows multiple write mode
	//currentFile *os.File	// current write file

	fileLength int64		// length of current active file

	mu sync.RWMutex			// mu on write
	filePool map[string]*os.File	// pool to use map

	// mutex for pool, the mu is used to operating on
	appendMu sync.Mutex			//the mutex to append data
	flock *flock.Flock


}

//func (poolManager *BitcaskBufferManager) getCurrentFile() *os.File {
//	f, err := poolManager.fetchFilePointer(dataFileName(poolManager.CurrentFileId, poolManager.dirName))
//	if err != nil {
//		log.Fatal(err)
//	}
//	return f
//}

func (poolManager *BitcaskBufferManager) closeFilePointerWithoutLock(fileName string)  {
	file, exists := poolManager.filePool[fileName]
	if !exists {
		return
	}
	file.Close()
}

func (poolManager *BitcaskBufferManager) closeFilePointer(fileName string)  {
	poolManager.mu.Lock()
	poolManager.closeFilePointerWithoutLock(fileName)
	poolManager.mu.Unlock()
}

func(poolManager *BitcaskBufferManager) closeAllFilePointer()  {
	for _, v := range poolManager.filePool {
		v.Close()
	}
}

func (poolManager *BitcaskBufferManager) fetchFilePointer(fileName string) (*os.File, error) {
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


func (poolManager *BitcaskBufferManager) Close() {
	poolManager.closeAllFilePointer()
	poolManager.currentWriter.Close()
}

func (*BitcaskBufferManager) switchFile()  {

}

// fetch bytes are all read operations, will not effect append.
func (poolManager *BitcaskBufferManager) FetchBytes(fileId uint32, valueSize uint32, valuePos uint32, timeStamp uint32) ([]byte, error) {

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


func (poolManager *BitcaskBufferManager) AppendRecord(data []byte) (uint32, uint32, uint32, uint32, error) {
	poolManager.appendMu.Lock()
	defer poolManager.appendMu.Unlock()
	// write data happens here

	n, err := poolManager.currentWriter.Write(data)

	oldStart := poolManager.fileLength
	poolManager.fileLength += int64(n)

	if err != nil {
		log.Infoln("write failed")
		return 0, 0, 0, 0, err
	}

	return uint32(poolManager.CurrentFileId), uint32(n), uint32(oldStart), 0, nil
}


func Open(dirName string, opts *options.Options) (*BitcaskBufferManager, error) {
	//fileLock := flock.New(path.Join(dirName, "bitcask.lock"))
	if opts == nil {
		opts = options.DefaultOption()
	}

	var fileLock *flock.Flock

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

	//ok, err := fileLock.TryLock()
	//if !ok {
	//	return nil, err
	//}

	//currentFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	//
	//if err != nil {
	//	return nil, err
	//}

	fileLength = 0

	poolManager := BitcaskBufferManager{
		flock:fileLock,
		CurrentFileId:currentFileId,
		CurrentFileName:fileName,
		dirName:dirName,
		fileLength:fileLength,
		filePool:make(map[string]*os.File),
	}

	poolManager.loadOptions(opts)

	return &poolManager, nil
}
