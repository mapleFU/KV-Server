package buffer

type BitcaskBuffer interface {
	// fetch data from store
	FetchBytes(fileId uint32, valueSize uint32, valuePos uint32, timeStamp uint32) ([]byte, error)
	// return
	AppendRecord([]byte) (uint32, uint32, uint32, uint32)
	Close()
}