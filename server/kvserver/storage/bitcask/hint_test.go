package bitcask

import (
	"testing"
	"os"
	"fmt"
	"strings"
)

const directoryPerm = 0766
func createTestDirectory(dirName string) error {
	err := os.Mkdir(dirName, directoryPerm)
	return err
}

func destructTestDirectory(dirName string)  {
	os.RemoveAll(dirName)
}

func TestLoadHint(t *testing.T)  {
	curTestDir := "testdata/test-hint"

	createTestDirectory(curTestDir)
	defer destructTestDirectory(curTestDir)

	bitcask := Open(curTestDir, nil)
	ok, err := bitcask.loadHint()
	if err != nil {
		t.Fatal(err)
	}
	if ok != false {
		t.Fatal("ok is not false when we use load")
	}

	for i := 0; i < 10; i++ {
		bitcask.Put([]byte(fmt.Sprintf("key-%d", i)), []byte(fmt.Sprintf("value-%d", i)))
	}
	bitcask.Close()
	bitcask = Open(curTestDir, nil)
	ok, err = bitcask.loadHint()
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("not ok when we use load")
	}
	for i := 0; i < 10; i++ {
		ret, err := bitcask.Get([]byte(fmt.Sprintf("key-%d", i)))
		if err != nil {
			t.Fatal(err)
		}
		sRet := string(ret)
		if strings.Compare(sRet, fmt.Sprintf("value-%d", i)) != 0 {
			t.Fatalf("in %d, Get not equal with we put", i)
		}
	}
}