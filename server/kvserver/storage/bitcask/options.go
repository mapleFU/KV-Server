package bitcask

import log "github.com/sirupsen/logrus"

// Options is a module to describe the arguments in bitcask
type Options struct {
	UseLog bool		// use wal log to keep messages in memory

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
}