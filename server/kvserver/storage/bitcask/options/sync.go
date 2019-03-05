package options

const (
	// default
	None SyncStrategy = iota	// lets the operating system manage syncing writes (default)
	O_Sync						//
	Interval
)

type Sync struct {
	Strategy SyncStrategy
	Interval int	// sync interval
}

// the strategy of synchronize database
type SyncStrategy int

func NewDefaultSync() Sync {
	return Sync{
		Strategy: None,
		Interval: 0,
	}
}