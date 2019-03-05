/*
Hint file has structure:
(uint32) | (uint32) | (uint32) | (uint32) | ...
tstamp | ksz | valuesz | value_pos | key

we use hint file as a backup of keyDir log
 */

package bitcask

import (
	"os"

	log "github.com/sirupsen/logrus"

	//"github.com/timtadh/data-structures/types"
	"bytes"
	"encoding/binary"
	"unsafe"
)

// hint file name in the system
const HintFileName = "hint"
const HintFilePerm = 0644

type hintHeader struct {
	Crc uint32
	TimeStamp uint32
	ValuePos uint32
	KeySz uint32
}

type hint struct {
	hintHeader
	key []byte
}

var hintHeaderSize int

func init()  {
	hintHeaderSize = int(unsafe.Sizeof(hintHeader{}))
}

// writeHint is a function called when closing file
// it backup the keyDir
func (bitcask *Bitcask) writeHint()  {
	hintFile, err := os.OpenFile(HintFileName, os.O_TRUNC|os.O_CREATE|os.O_RDWR, HintFilePerm)
	if err != nil {
		// Errors should not happens here
		log.Fatal(err)
	}
	//var iter types.KIterator
	iter := bitcask.entryMap.iterator()
	bitcask.entryMap.addRecord()

	for k, v, iter := iter(); iter != nil; k, v, iter = iter() {
		keyBytes := []byte(string(k.(hashKey)))
		value := v.(*entry)

		// crc not implemented
		header := hintHeader{
			TimeStamp:value.Timestamp,
			ValuePos:value.ValuePos,
			KeySz: uint32(len(keyBytes)),
		}
		appendHintToFile(header, keyBytes, hintFile)
	}

}

// loadHint is a function to load the hint file
// if the hintfile does not
func (bitcask *Bitcask) loadHint()  {
	_, err := os.Open(HintFileName)
	if err != nil {
		log .Info(err)
		return
	}
	panic("impl me")
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