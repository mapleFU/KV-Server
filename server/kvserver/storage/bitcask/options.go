package bitcask

import (
	"github.com/mapleFU/KV-Server/server/kvserver/storage/bitcask/options"

	log "github.com/sirupsen/logrus"
)


func (bitcask *Bitcask) optionInit(options *options.Options)  {
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