package schema

import (
	"time"
	"hash/crc32"
	"encoding/binary"
	"bytes"


	log "github.com/sirupsen/logrus"
)

/**
持久化存储的
(uint32) | (uint32) | (uint32) | (uint32) | ...
crc | tstamp | ksz | value_sz | k | v
 */

 /**
 header
| crc | tstamp | ksz | value_sz |
  */
type BitcaskStoreHeader struct {
	Crc uint32
	TimeStamp uint32
	//timeStamp uint64	// u64 unixNano
	KeySz uint32
	ValueSz uint32
}

const PersistHeaderSize  = 32 / 8 + 32 / 8 + 32 / 8 + 32 / 8

/**
key: key of key/value
value: value of key/value
timeStamp: the time to store
 */
func PersistEncoding(key []byte, value []byte, timeStamp time.Time) []byte {
	crc := crc32.ChecksumIEEE(value)

	keySize := uint32(len(key))
	// int64
	unixNano := uint32(timeStamp.UnixNano())
	valueSize := uint32(len(value))
	bufSize := PersistHeaderSize + keySize + valueSize

	buf := make([]byte, bufSize)

	// write data to buffer
	header := BitcaskStoreHeader{
		Crc:crc,
		TimeStamp: unixNano,
		KeySz:keySize,
		ValueSz:valueSize,
	}
	byteBuf := bytes.NewBuffer(make([]byte, 0))
	err := binary.Write(byteBuf, binary.LittleEndian, header)
	if err != nil {
		log.Fatal(err)
	}
	copy(buf[:PersistHeaderSize], byteBuf.Bytes())
	//log.Infof("Length header %d -- %d\n", PersistHeaderSize, byteBuf.Len())
	//binary.LittleEndian.PutUint32(buf[:4], crc)	// crc
	//binary.LittleEndian.PutUint32(buf[4:8], unixNano)
	//binary.LittleEndian.PutUint32(buf[8:12], keySize)
	//binary.LittleEndian.PutUint32(buf[12:16], valueSize)

	copy(buf[PersistHeaderSize:PersistHeaderSize+keySize], key)

	//log.Infof("Range Key %d -- %d\n", keySize, len(key))
	copy(buf[PersistHeaderSize+keySize:PersistHeaderSize+keySize+valueSize], value)

	//log.Infof("Range Value %d -- %d\n", valueSize, len(value))
	return buf
}

func PersistDecoding(data []byte) ([]byte, []byte, time.Time) {
	header := BitcaskStoreHeader{}
	//fmt.Printf("%d -- %d\n", PersistHeaderSize, unsafe.Sizeof(bitcaskStoreHeader{}))

	binary.Read(bytes.NewBuffer(data[:PersistHeaderSize]), binary.LittleEndian, &header)
	ksz, vsz := header.KeySz, header.ValueSz
	unixNano := header.TimeStamp
	UTCfromUnixNano := time.Unix(0, int64(unixNano))

	return data[PersistHeaderSize:PersistHeaderSize+ksz],
	data[PersistHeaderSize+ksz:PersistHeaderSize+ksz+vsz], UTCfromUnixNano
}

