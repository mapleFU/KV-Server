package bitcask

import (
	"testing"
	"os"
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
	createTestDirectory("testdata/test-hint")

	destructTestDirectory("testdata/test-hint")
}