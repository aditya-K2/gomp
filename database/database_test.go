package database

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/aditya-K2/gomp/utils"
)

var (
	a, _     = time.ParseDuration("1h10m10s")
	b, _     = time.ParseDuration("10m10s")
	c, _     = time.ParseDuration("11h11m10s")
	readTMap = map[string]time.Duration{
		"/Hello/What": a,
		"/Wi/What":    b,
		"/io/What":    c,
	}
)

func TestRead(t *testing.T) {
	dbPath = "./testdata/test_db"
	Read()
	for k, v := range fmap {
		if readTMap[k] != v {
			t.Errorf("Mismatch for key: %s, received: %s, expected: %s", k, v, readTMap[k])
		}
	}
	dbPath = "./testdata/test_db_non_existant"
	if utils.FileExists(dbPath) {
		t.Errorf("%s already exists. expected: shouldn't exist", dbPath)
	}
	fmap = map[string]time.Duration{}
	Read()
	if !utils.FileExists(dbPath) {
		t.Errorf("%s was not created. expected: creation", dbPath)
	}
	if err := os.Remove(dbPath); err != nil {
		t.Error(err)
	}
}

func TestWrite(t *testing.T) {
	dbPath = "./testdata/test_db_write"
	fmap = readTMap
	Write()
	content1, cerr1 := os.ReadFile(dbPath)
	_dbPath := "./testdata/test_db"
	content2, cerr2 := os.ReadFile(_dbPath)
	if cerr1 != nil {
		t.Error(cerr1)
	}
	if cerr2 != nil {
		t.Error(cerr2)
	}
	if !bytes.Equal(content1, content2) {
		t.Errorf("Content from %s and %s not equal.", dbPath, _dbPath)
	}
	if err := os.Remove(dbPath); err != nil {
		t.Error(err)
	}
}
