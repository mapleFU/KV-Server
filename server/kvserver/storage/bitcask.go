package storage

import (
	"github.com/laohanlinux/bitcask"
	log "github.com/sirupsen/logrus"
)

type Bitcask struct {
	bc *bitcask.BitCask
}

func Open(dirName string) *Bitcask {

	bc, err := bitcask.Open(dirName, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &Bitcask{
		bc:bc,
	}
}

func (bc *Bitcask) Close()  {
	defer bc.bc.Close()
}

func (bc *Bitcask) Get(key []byte) ([]byte, error) {
	return bc.bc.Get(key)
}

func (bc *Bitcask) Set(key []byte, value []byte) error {
	
	return bc.bc.Put(key, value)
}

func (bc *Bitcask) Scan()  {
	panic("impl me!")
}

func (bc *Bitcask) Delete(key []byte) error {
	return bc.bc.Del(key)
}