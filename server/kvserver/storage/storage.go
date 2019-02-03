package storage

type ScanCursor struct {
	// cursor
	Cursor int
	// match
	UseMatchKey bool
	MatchKeyString string
	// count has a default: 10
	Count int
}

type Storage interface {
	Close()

	Scan(cursor ScanCursor)	([][]byte, error)
	Get([]byte) ([]byte, error)
	Put([]byte, []byte) error
	Del([]byte) error
}
