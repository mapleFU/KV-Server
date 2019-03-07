package buffer

import (

	"github.com/mapleFU/KV-Server/server/kvserver/storage/bitcask/options"

	log "github.com/sirupsen/logrus"
)

func (bm *BitcaskBufferManager) loadOptions(op *options.Options) error {
	if op == nil {
		log.Info("options in loadOptions is nil")
		op = options.DefaultOption()
	}
	// do sync
	var err error
	bm.currentWriter, err = newWriter(op.Sync, bm.CurrentFileName)
	if err != nil {
		return err
	}

	// TODO: impl this
	switch op.Merge.Policy {
	case options.Always:

	case options.Never:

	case options.Windows:

	}

	return nil
}