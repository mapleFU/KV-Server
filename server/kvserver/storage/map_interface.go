package storage

type entryMapper interface {
	flushRecord(key string, entry *entry) error
	addRecord(key string, entry *entry) error
	scan(cursor ScanCursor) ([]*entry, *ScanCursor, error)
	deleteRecord(key string)
	getEntry(key string) (*entry, bool)
}
