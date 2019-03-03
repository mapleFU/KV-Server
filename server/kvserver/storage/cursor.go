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