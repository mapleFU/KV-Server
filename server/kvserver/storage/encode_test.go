package storage

import (
	"testing"
	"time"
	"bytes"
)

func TestEncodeDecode(t *testing.T)  {
	key := []byte("test-key")
	value := []byte("test-value")

	currentTime := time.Now()

	encodeResult := persistEncoding(key, value, currentTime)
	if len(encodeResult) != PersistHeaderSize + len(key) + len(value) {
		t.Fatal("Length of persist encoding error")
	}

	key2, value2, _ := persistDecoding(encodeResult)
	if bytes.Compare(key, key2) != 0 {
		t.Fatal("key marshall or unmarshall error")
	}

	if bytes.Compare(value, value2) != 0 {
		t.Fatal("value marshall or unmarshal error")
	}

}