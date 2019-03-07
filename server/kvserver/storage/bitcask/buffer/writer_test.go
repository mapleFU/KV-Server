package buffer

import (
	"testing"
	"os"
	"fmt"
	"io/ioutil"
	"bytes"
	"path"

	"github.com/mapleFU/KV-Server/server/kvserver/storage/bitcask/options"

	log "github.com/sirupsen/logrus"
	"io"
)

func TestBufWriter(t *testing.T)  {
	testDir := "testdata/writer"
	err := os.Mkdir(testDir, 0777)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	mode := []options.SyncStrategy{ options.Interval, options.None, options.O_Sync,}
	var sync options.Sync
	sync.Interval = 60
	for i, v := range mode {
		log.Infof("Loop %d started", i)
		sync.Strategy = v
		curFileName := path.Join(testDir, fmt.Sprintf("test-%d", i))

		curW, err := newWriter(sync, curFileName)
		if err != nil {
			t.Fatal(err)
		}
		wdata := []byte(fmt.Sprintf("text-%d", i))
		curW.Write(wdata)

		err = curW.Flush()
		if err != nil {
			t.Fatal(err)
		}
		curW.Close()
		f, err := os.Open(curFileName)
		if err != nil {
			t.Fatal(err)
		}

		data, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(data, wdata) != 0 {
			t.Fatalf("In loop %d, compare got %d, error", i, bytes.Compare(data, wdata))
		}
	}
}

// (bug)
func TestWriterInterval_Write(t *testing.T) {
	testDir := "testdata/writerInterval"
	err := os.Mkdir(testDir, 0777)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)
	var sync options.Sync
	sync.Strategy = options.Interval
	sync.Interval = 60

	curFileName := path.Join(testDir, "test-data")


	curW, err := newWriter(sync, curFileName)

	if err != nil {
		t.Fatal(err)
	}
	wi, f := curW.(*WriterInterval)
	if !f {
		t.Fatal("cannot cast object to writer-interval")
	}

	wdata := []byte("test-value")
	curW.Write(wdata)
	rbuf := make([]byte, len(wdata))

	n, err := wi.currentFile.ReadAt(rbuf, 0)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}

	log.Infof("Read %d bytes: %v", n, string(rbuf))
	if bytes.Compare(rbuf, wdata) != 0 {
		wi.Flush()
		debugBuf := make([]byte, len(wdata))
		_, err = wi.currentFile.ReadAt(debugBuf, 0)
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		t.Fatalf("read in file cannot match write, read %v, debugBuf is %v", string(rbuf), string(debugBuf))

	}
}