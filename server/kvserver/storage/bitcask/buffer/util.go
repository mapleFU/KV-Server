/*
 TODO: change the relationship between bitcask datafile
 */
package buffer

import (
	"fmt"
	"path"
	"regexp"
	"path/filepath"
	"strconv"
	"os"
	"unsafe"
	"encoding/binary"
	"bytes"
	"io"

	log "github.com/sirupsen/logrus"
)

/**
	dataFileName is a function to get datafile name from datafile id and dirName
	It's not about buffer logic, just combine the os tools
 */
func dataFileName(dataFileID int, dirName string) string {
	return path.Join(dirName, fmt.Sprintf("%d.data", dataFileID))
}

const fileReg  = `(\d+).data`

func listDataFileID(dirName string) ([]int, error) {
	fileNameReg := regexp.MustCompile(fileReg)
	matches, err := filepath.Glob(path.Join(dirName, "*.data"))

	retArr := make([]int, 0)
	if err != nil {
		return retArr, err
	}

	for _, v := range matches {
		fmt.Println(v)
		s := fileNameReg.FindStringSubmatch(v)
		id, err := strconv.Atoi(s[1])
		if err != nil {
			log.Fatal(err)
		}
		retArr = append(retArr, id)
		// set currentfile id maxest

	}
	return retArr, nil
}

type Record struct {
	Key []byte							// key
	Value []byte						// value
	Header *BitcaskStoreHeader
}

/**
read records from a data file
 */
func ReadRecords(file *os.File) ([]*Record, error) {

	retRecords := make([]*Record, 0)

	header := BitcaskStoreHeader{}
	buf := make([]byte, unsafe.Sizeof(header))
	var bios int64
	bios = 0
	stat, err := file.Stat()
	if err != nil {
		return retRecords, err
	}
	fileSize := stat.Size()
	for fileSize > bios {
		log.Infoln(fileSize, bios)
		curRecord := new(Record)
		read, err := file.ReadAt(buf, bios)
		if err != nil {
			// meet eof
			if err == io.EOF {
				return retRecords, nil
			}
			return retRecords, err
		}
		bios += int64(read)
		binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &header)
		curRecord.Header = &header
		keySize := header.KeySz
		key := make([]byte, keySize)
		read, err = file.ReadAt(buf, bios)
		bios += int64(read)
		if err != nil {
			return retRecords, err
		}
		curRecord.Key = key

		// not delete
		if header.ValueSz != 0 {
			valueSize := header.ValueSz
			value := make([]byte, valueSize)
			read, err = file.ReadAt(buf, bios)
			if err != nil {
				if err == io.EOF {
					return retRecords, nil
				}
				return retRecords, err
			}
			bios += int64(read)
			curRecord.Value = value
		}

		retRecords = append(retRecords, curRecord)
	}
	return retRecords, nil
}