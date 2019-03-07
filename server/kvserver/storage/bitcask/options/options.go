package options

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




// DefaultOption is a function creates default settings for our project
func DefaultOption() *Options {
	option := Options{}
	// maybe I don't know how to use log, so I set it false here
	option.UseLog = false
	// merge policy will be set "Always"
	option.Merge = NewDefaultMerge()

	option.MaxFileSize = 1024 * 500 // maybe 500M?

	// every data should be write to disk
	option.Sync = NewDefaultSync()

	return &option
}
