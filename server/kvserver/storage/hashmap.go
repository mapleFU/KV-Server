package storage

// entry map to store entry records
type entryMap map[string]*entry

func newEntryMap() *entryMap {
	newMap := make(map[string]*entry)
	var retMap entryMap
	retMap = entryMap(newMap)
	return &retMap
}
/**
reflesh record
 */
func (entryMap *entryMap) flushRecord(key string, entry *entry) error {
	(*entryMap)[key] = entry
	return nil
}

/**
add record with key
 */
func (entryMap *entryMap) addRecord(key string, entry *entry) error {
	(*entryMap)[key] = entry
	return nil
}

func (entryMap entryMap) scan(cursor ScanCursor) ([]*entry, *ScanCursor, error) {
	ret := make([]*entry, 0)
	endCursor := new(ScanCursor)
	for k, v := range entryMap {
		if k != cursor.cursor {
			continue
		}
		// TODO: make it better later
		ret = append(ret, v)
		endCursor.cursor = k
	}
	return ret, endCursor, nil
}

func (entryMap entryMap) deleteRecord(key string)  {
	delete(entryMap, key)
}

func (entryMap entryMap) getEntry(key string) (*entry, bool) {
	entry, ok := entryMap[key]
	return entry, ok
}