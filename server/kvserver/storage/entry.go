/**
实现 bitcask 的 entry, bitcask 将其保存在内存中

crc | tstamp | ksz | value_sz | k | v

key ----> | file_id | value_sz | value_pos | tstamp |
 */
package storage

import "time"

type entry struct {
	fileID uint32
	valueSize uint32
	valuePos uint32
	timestamp time.Time
}

/**
return entry with current time
 */
func createEntryInCurrentTime(file uint32, valueSize uint32, valuePos uint32) *entry {
	currentTime := time.Now()
	return &entry{
		fileID:file,
		valueSize:valueSize,
		valuePos: valuePos,
		timestamp: currentTime,
	}
}
