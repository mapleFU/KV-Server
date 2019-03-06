package buffer

import (
	"os"
	"bufio"
	"time"

	"github.com/mapleFU/KV-Server/server/kvserver/storage/bitcask/options"

	log "github.com/sirupsen/logrus"
	"github.com/pkg/errors"
)

type BufWriter interface {
	Write([]byte) (int, error)
	Flush() error
	Close() error
}

func newWriter(sync options.Sync, fileName string) (BufWriter, error) {

	switch sync.Strategy {
	case options.Interval:
		return newWriterInterval(fileName, sync.Interval)
	case options.O_Sync:
		return newWriterOSync(fileName)
	case options.None:
		return newWriterNone(fileName)
	default:
		log.Fatal("sync.Strategy got unexcepted value")

		return nil, errors.New("sync.Strategy got unexcepted value")
	}

}

type WriterNone struct {
	currentFile *os.File
}

func newWriterNone(fileName string) (*WriterNone, error) {
	currentFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)

	if err != nil {
		return nil, err
	}
	return &WriterNone{
		currentFile: currentFile,
	}, nil
}

func (wn *WriterNone) Write(b []byte) (int, error) {

	return wn.currentFile.Write(b)
}

func (wn *WriterNone) Flush() error {
	return wn.currentFile.Sync()
}

func (wn *WriterNone) Close() error {
	return wn.currentFile.Close()
}

type WriterOSync struct {
	currentFile *os.File
}

func newWriterOSync(fileName string) (*WriterOSync, error) {
	// create  with o_sync
	currentFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_SYNC, 0777)

	if err != nil {
		return nil, err
	}
	return &WriterOSync{
		currentFile: currentFile,
	}, nil
}

func (wos *WriterOSync) Write(b []byte) (int, error) {
	return wos.currentFile.Write(b)
}

func (wos *WriterOSync) Flush() error {
	// doesn't need to flush, hahaha
	return nil
}

func (wos *WriterOSync) Close() error {
	return wos.currentFile.Close()
}


type WriterInterval struct {
	currentFile *os.File
	bufWriter *bufio.Writer

	closeBuf chan struct{}
}

func (wi *WriterInterval) Write(b []byte) (int, error) {
	return wi.bufWriter.Write(b)
}

func (wi *WriterInterval) Flush() error {
	err := wi.bufWriter.Flush()
	if err != nil {
		return err
	}
	return wi.currentFile.Sync()
}

func (wi *WriterInterval) Close() error {
	// wait for it close well
	wi.closeBuf <- struct{}{}
	close(wi.closeBuf)
	return wi.currentFile.Close()
}

func newWriterInterval(fileName string, interval int) (*WriterInterval, error) {
	// create  with o_sync
	currentFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriterSize(currentFile, 4092 * 2)

	closeBuf := make(chan struct{})

	wi := &WriterInterval{
		currentFile: currentFile,
		bufWriter:writer,
		closeBuf: closeBuf,
	}

	go func() {
		inLoop := true
		for inLoop  {
			select {
			case <-closeBuf:
				inLoop = false
			case <- time.After(time.Second * time.Duration(interval)):
				// clear buf
				wi.Flush()
			}
		}

	}()

	return wi, nil
}
