package bitcask

import log "github.com/sirupsen/logrus"

// Options is a module to describe the arguments in bitcask
type Options struct {
	UseLog bool		// use wal log to keep messages in memory
	Sync Sync
	// max size of the file, remember it's not a standard value, if a file grows larger than
	// MaxFileSize, the bitcask will call switchfile and set the current active file old and create
	// another active file
	MaxFileSize int
	// Merge controls the policy of merge, about how we
	Merge Merge
}

// the strategy of synchronize database
type SyncStrategy int

const (
	None SyncStrategy = iota
	O_Sync
	Interval
)

type Sync struct {
	Strategy SyncStrategy
	Interval int	// sync interval
}

type MergePolicy int

const (
	Always MergePolicy = iota	// when we get a new file, we will
	Never 						// never merge files
	Windows 					// Merge operations occur during specified hours
)

type MergeWindow struct {
	Start int
	End int
}

type Merge struct {
	Policy MergePolicy
	Window MergeWindow
	Interval int
}


// DefaultOption is a function creates default settings for our project
func DefaultOption() *Options {
	option := Options{}
	// maybe I don't know how to use log, so I set it false here
	option.UseLog = false
	// merge policy will be set "Always"
	option.Merge.Policy = Always
	option.Merge.Interval = 3600 // every one hour

	option.MaxFileSize = 1024 * 500 // maybe 500M?

	// every data should be write to disk
	option.Sync.Strategy = O_Sync

	return nil
}


func (bitcask *Bitcask) optionInit(options *Options)  {
	if options == nil {
		return
	}
	if options.UseLog {
		logger, err := newRedoLogger(bitcask.directoryName)
		if err != nil {
			log.Fatalln(err)
		}
		bitcask.redoLogger = logger
	}

	panic("impl me")
}