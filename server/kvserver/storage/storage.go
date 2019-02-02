package storage

type ScanCursor struct {
	cursor string
}

type Storage interface {
	Close()

	Scan(cursor ScanCursor)	([][]byte, error)
	Get([]byte) ([]byte, error)
	Put([]byte, []byte) error
	Del([]byte) error
}
