package buffer

import (
	"testing"
	log "github.com/sirupsen/logrus"
)


func TestOpen(t *testing.T) {
	bufPool := Open("/Users/fuasahi/GoglandProjects/src/github.com/mapleFU/KV-Server/server/testdata")
	log.Info(bufPool.currentFileId)
	if bufPool.currentFileId != 0 {
		t.Fatal("Error, testdir not get 0")
	}
}