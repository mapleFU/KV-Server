package bitcask

// switchActive
func (bc *Bitcask) switchActive()  {

	// get old data file name
	oldFileId := bc.bitcaskPoolManager.CurrentFileId

	bc.entryMap.mu.Lock()
	defer bc.entryMap.mu.Unlock()

	oldHash := bc.entryMap.immutableScanMap()
	bc.bitcaskPoolManager.SwitchFile()

	// write oldHash to hint
	go func() {
		bc.buildHint(oldFileId, oldHash)
	}()
}