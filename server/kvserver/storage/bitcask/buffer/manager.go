package buffer

import (
	"os"
	"sync"

	"github.com/mapleFU/KV-Server/server/kvserver/storage/bitcask/options"

	log "github.com/sirupsen/logrus"
	"github.com/gofrs/flock"
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

	opts *options.Options
	SwitchChan chan struct{}
}

//func (pbm *BitcaskBufferManager) getCurrentFile() *os.File {
//	f, err := pbm.fetchFilePointer(dataFileName(pbm.CurrentFileId, pbm.dirName))
//	if err != nil {
//		log.Fatal(err)
//	}
//	return f
//}

func (pbm *BitcaskBufferManager) closeFilePointerWithoutLock(fileName string)  {
	file, exists := pbm.filePool[fileName]
	if !exists {
		return
	}
	file.Close()
}

func (pbm *BitcaskBufferManager) closeFilePointer(fileName string)  {
	pbm.mu.Lock()
	pbm.closeFilePointerWithoutLock(fileName)
	pbm.mu.Unlock()
}

func(pbm *BitcaskBufferManager) closeAllFilePointer()  {
	for _, v := range pbm.filePool {
		v.Close()
	}
}

func (pbm *BitcaskBufferManager) fetchFilePointer(fileName string) (*os.File, error) {
	pbm.mu.RLock()
	if file, exists := pbm.filePool[fileName]; exists {
		defer pbm.mu.RUnlock()
		return file, nil
	} else {
		pbm.mu.RUnlock()
		pbm.mu.Lock()

		defer pbm.mu.Unlock()
		if file, exists := pbm.filePool[fileName]; exists {

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


func (pbm *BitcaskBufferManager) Close() {
	pbm.closeAllFilePointer()
	close(pbm.SwitchChan)
	pbm.currentWriter.Close()
}

// SwitchFile is a function to change current write to another data file, and
// set current active data file "old"
// what it actually do is:
// 1. flush buffer of current writer
// 2. reset the current file id and name
// 3. reset current file w length
// 4. renew a writer for new file
// finally Close current writer
func (pbm *BitcaskBufferManager) SwitchFile()  {
	pbm.mu.Lock()
	defer pbm.mu.Unlock()
	err := pbm.currentWriter.Flush()
	if err != nil {
		log.Fatal(err)
	}
	pbm.currentWriter.Flush()
	defer pbm.currentWriter.Close()

	pbm.CurrentFileId++
	pbm.CurrentFileName = dataFileName(pbm.CurrentFileId, pbm.dirName)
	pbm.fileLength = 0

	pbm.currentWriter, err = newWriter(pbm.opts.Sync, pbm.CurrentFileName)
	if err != nil {
		log.Fatal(err)
	}
}

// FetchBytes are all read operations, will not effect append.
func (pbm *BitcaskBufferManager) FetchBytes(fileId uint32, valueSize uint32, valuePos uint32, timeStamp uint32) ([]byte, error) {
	// TODO: maintain the buffer will myself
	if fileId == uint32(pbm.CurrentFileId) {
		pbm.currentWriter.Flush()
	}

	fileName := dataFileName(int(fileId), pbm.dirName)
	// maybe i should get from a read pool...
	file, err := pbm.fetchFilePointer(fileName)
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


func (pbm *BitcaskBufferManager) AppendRecord(data []byte) (uint32, uint32, uint32, uint32, error) {
	pbm.appendMu.Lock()
	defer pbm.appendMu.Unlock()
	// write data happens here

	n, err := pbm.currentWriter.Write(data)

	oldStart := pbm.fileLength
	pbm.fileLength += int64(n)

	if err != nil {
		log.Infoln("write failed")
		return 0, 0, 0, 0, err
	}

	if oldStart < int64(pbm.opts.MaxFileSize) && oldStart + int64(n) >=  int64(pbm.opts.MaxFileSize) {
		go func() {
			pbm.SwitchChan <- struct{}{}
		}()
	}

	return uint32(pbm.CurrentFileId), uint32(n), uint32(oldStart), 0, nil
}


func Open(dirName string, opts *options.Options) (*BitcaskBufferManager, error) {
	//fileLock := flock.New(path.Join(dirName, "bitcask.lock"))
	if opts == nil {
		opts = options.DefaultOption()
	}

	var fileLock *flock.Flock

	var fileName string
	var currentFileId int
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

	pbm := BitcaskBufferManager{
		flock:fileLock,
		CurrentFileId:currentFileId,
		CurrentFileName:fileName,
		dirName:dirName,
		fileLength:fileLength,
		filePool:make(map[string]*os.File),
		SwitchChan: make(chan struct{}),
		opts:opts,
	}

	pbm.loadOptions(opts)

	return &pbm, nil
}
