package bitcask

import log "github.com/sirupsen/logrus"

type Options struct {
	UseLog bool
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