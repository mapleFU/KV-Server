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