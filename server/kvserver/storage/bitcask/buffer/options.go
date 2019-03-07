package buffer

import (

	"github.com/mapleFU/KV-Server/server/kvserver/storage/bitcask/options"

	log "github.com/sirupsen/logrus"
	"time"
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
	var conditionF func(time2 time.Time) bool

	if options.Windows == bm.opts.Merge.Policy {
		conditionF = func(time2 time.Time) bool {
			return time2.Hour() < bm.opts.Merge.Window.End && time2.Hour() > bm.opts.Merge.Window.Start
		}
	} else if bm.opts.Merge.Policy == options.Always {
		conditionF = func(time2 time.Time) bool {
			return true
		}
	}

	go func() {
		inLoop := true
		for inLoop {
			if !conditionF(time.Now()) {
				select {
				case bm.MergeChan <- struct{}{}:
					break
				default:
					inLoop = false
				}
			}
			time.Sleep(time.Minute * 30)
		}

	}()

	return nil
}