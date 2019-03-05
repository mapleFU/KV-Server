/*
Hint file has structure in paper is:
(uint32) | (uint32) | (uint32) | (uint32) | ...
tstamp | ksz | valuesz | value_pos | key

I think maybe I can make some difference:


we use hint file as a backup of keyDir log
 */

package bitcask

import (
	"os"
	"bytes"
	"unsafe"
	"encoding/binary"

	log "github.com/sirupsen/logrus"
	"io"
)

// hint file name in the system
const HintFileName = "hint"
const HintFilePerm = 0644

type hintHeader struct {
	Crc uint32
	TimeStamp uint32
	ValuePos uint32
	KeySz uint32
	ValueSize uint32

	// entry structure, it was an extra field
	FileID uint32

}

var hintHeaderSize int

func init()  {
	hintHeaderSize = int(unsafe.Sizeof(hintHeader{}))
}

// writeHint is a function called when closing file
// it backup the keyDir
func (bitcask *Bitcask) syncKeyDirToHint()  {
	hintFile, err := os.OpenFile(HintFileName, os.O_TRUNC|os.O_CREATE|os.O_RDWR, HintFilePerm)
	if err != nil {
		// Errors should not happens here
		log.Fatal(err)
	}
	//var iter types.KIterator
	iter := bitcask.entryMap.iterator()

	for k, v, iter := iter(); iter != nil; k, v, iter = iter() {
		keyBytes := []byte(string(k.(hashKey)))
		value := v.(*entry)

		// crc not implemented
		header := hintHeader{
			TimeStamp:value.Timestamp,
			ValuePos:value.ValuePos,
			KeySz: uint32(len(keyBytes)),
			ValueSize: value.ValueSize,
			FileID:value.FileID,
		}
		err = appendHintToFile(header, keyBytes, hintFile)
		if err != nil {
			log.Fatal(err)
		}
	}

}

// loadHint is a function to load the hint file to keyDir when we Open a bitcask directory and load
// the project.
// if the hintfile does not exists(like we just begin a Bitcask), no errors will be throw,
// the return will be whether we have a HintFile
func (bitcask *Bitcask) loadHint() (bool, error) {
	hintFile, err := os.Open(HintFileName)
	if err != nil {
		log.Info(err)
		return false, nil
	}
	bios := 0
	loop := true
	for loop {
		curKey, curEntry, curRead, err := loadHintFromFile(bios, hintFile)
		if err != nil {
			if err == io.EOF {
				loop = false
			} else {
				return false, err
			}
		}
		bitcask.entryMap.addRecord(curKey, &curEntry)

		bios += curRead
	}

	return true, nil
}

// loadHintFromFile will load the hint in the system
// TODO: it will well handle EOF.
func loadHintFromFile(biosBeg int, file *os.File) (key string, entry entry, nRead int, err error) {
	// read
	nRead = 0
	readBuf := make([]byte, hintHeaderSize)
	n, err := file.ReadAt(readBuf, int64(biosBeg))
	if err != nil {
		return
	}
	var header hintHeader
	err = binary.Read(bytes.NewBuffer(readBuf), binary.LittleEndian, &header)
	if err != nil {
		return
	}
	nRead += n
	// read key
	readBuf = make([]byte, header.KeySz)
	n, err = file.ReadAt(readBuf, int64(biosBeg + n))
	if err != nil {
		if err != io.EOF {
			return
		}
	}

	key = string(readBuf)
	nRead += n
	entry.ValuePos = header.ValuePos
	entry.Timestamp = header.TimeStamp
	entry.ValueSize = header.ValueSize
	// extra op on fileID
	entry.FileID = header.FileID

	return
}

func appendHintToFile(header hintHeader, keyBytes []byte, file *os.File) error {
	// append
	length := hintHeaderSize + len(keyBytes)
	buf := make([]byte, length)
	wbuf := new(bytes.Buffer)

	err := binary.Write(wbuf, binary.LittleEndian, header)
	if err != nil {
		return err
	}
	copy(buf[:hintHeaderSize], wbuf.Bytes())
	copy(buf[hintHeaderSize:], keyBytes)

	_, err = file.Write(buf)
	if err != nil {
		return err
	} else {
		return nil
	}
}