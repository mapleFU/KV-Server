package storage

import "testing"

func TestOpen(t *testing.T) {
	bitcask := Open("/Users/fuasahi/GoglandProjects/src/github.com/mapleFU/KV-Server/server/data")
	if bitcask == nil {
		t.Fatal("bitcask is nil")
	}
	if bitcask.bitcaskPoolManager == nil {
		t.Fatal("bitcask.bitcaskPoolManager is nill")
	}

	if bitcask.entryMap == nil {
		t.Fatal("bitcask.entryMap is nil")
	}
}