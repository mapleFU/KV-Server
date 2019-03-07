package bitcask

// Merge is a function to merge data files
func (bc *Bitcask) Merge(directoryName string) error {

	b, err := bc.bitcaskPoolManager.Compaction()
	if err != nil {
		return err
	}
	if !b {
		return nil
	}
	// do hint
	oldHash := bc.entryMap.immutableScanMap()

	// write oldHash to hint
	go func() {
		bc.buildHint(0, oldHash)
	}()
}