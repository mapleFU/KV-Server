/**
实现 bitcask 的 entry, bitcask 将其保存在内存中

crc | tstamp | ksz | value_sz | k | v

key ----> | file_id | value_sz | value_pos | tstamp |
 */
package storage

import "time"

type entry struct {
	FileID uint32
	ValueSize uint32
	ValuePos uint32
	Timestamp uint32
}

/**
return entry with current time
just use for test...
 */
func CreateEntryInCurrentTime(file uint32, valueSize uint32, valuePos uint32) *entry {
	currentTime := time.Now().UnixNano()
	return &entry{
		FileID:file,
		ValueSize:valueSize,
		ValuePos: valuePos,
		Timestamp: uint32(currentTime),
	}
}
