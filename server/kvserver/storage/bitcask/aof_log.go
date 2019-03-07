/**
Log
uint32 | uint32  | uint32    | uint32   | ...
tstamp | keySize | valueSize | valuePos | key
 */
package bitcask

import (

	"unsafe"
	"encoding/binary"
	"bytes"
	"os"
	"syscall"
	"path"
	"io"
)

type AofLogHeader struct {
	Timestamp uint32
	KeySize uint32
	ValueSize uint32
	FileID uint32
	ValuePos uint32
}

type aofLogger struct {
	dirName string
	logFile *os.File
}

func newRedoLogger(dirName string) (*aofLogger, error) {
	fileName := path.Join(dirName,"log")
	var currentFile *os.File

	currentFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &aofLogger{
		dirName:dirName,
		logFile:currentFile,
	}, nil
}

func (redoLogger *aofLogger) writeLog(entry *entry, key string)  {
	bytesData := writeWalLog(entry, key)
	writeLogToFile(redoLogger.logFile, bytesData)
}

func (redoLogger *aofLogger) readLog(bios int64) (*entry, string, int, error) {
	return readLog(redoLogger.logFile, bios)
}

func writeWalLog(entry *entry, key string) []byte {
	keyBytes := []byte(key)
	walHeader := AofLogHeader{
		KeySize:uint32(len(keyBytes)),
		FileID:entry.FileID,
		ValuePos:entry.ValuePos,
		Timestamp:entry.Timestamp,
		ValueSize:entry.ValueSize,
	}
	headerSize := unsafe.Sizeof(walHeader)

	length := int(headerSize) + len(keyBytes)
	buf := make([]byte, length)
	wbuf := new(bytes.Buffer)

	binary.Write(wbuf, binary.LittleEndian, walHeader)
	copy(buf[:int(headerSize)], wbuf.Bytes())
	copy(buf[int(headerSize):], keyBytes)
	return buf
}

func writeLogToFile(file *os.File, log []byte) (int, error) {
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	bios := stat.Size()
	data, err := file.WriteAt(log, bios)
	syscall.Fsync(int(file.Fd()))
	return data, err
}

type EntryKV struct {
	key string
	entry *entry
}


func readLog(file *os.File, bios int64) (*entry, string, int, error) {
	var retErr error
	retErr = nil
	var walHeader AofLogHeader
	headerSize := int(unsafe.Sizeof(walHeader))
	headerBuf := make([]byte, headerSize)
	file.ReadAt(headerBuf, bios)

	err := binary.Read(bytes.NewBuffer(headerBuf), binary.LittleEndian, &walHeader)
	if err != nil {
		return nil, "", 0, err
	}
	keyBuf := make([]byte, walHeader.KeySize)
	read, err := file.ReadAt(keyBuf, bios + int64(headerSize))
	if err != nil {
		if err != io.EOF {
			return nil, "", 0, err
		} else {
			retErr = io.EOF
		}

	}

	return &entry{
		FileID:walHeader.FileID,
		ValuePos:walHeader.ValuePos,
		ValueSize:walHeader.ValueSize,
		Timestamp:walHeader.Timestamp,
	}, string(keyBuf), read + headerSize, retErr
}