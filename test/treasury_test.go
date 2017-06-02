package test

import (
	"testing"

	"github.com/slemgrim/treasury/treasury"
)

type MockStorage struct {
	HasBeenSaved  bool
	HasBeenLoaded bool
}

func (s *MockStorage) Write(file *treasury.File) error {
	s.HasBeenSaved = true
	return nil
}

func (s *MockStorage) Read(id string) (treasury.File, error) {
	s.HasBeenLoaded = true
	return *new(treasury.File), nil
}

func TestSaveCallsStorage(t *testing.T) {
	s := new(MockStorage)
	treasury := treasury.NewTreasury(s)
	treasury.Save("./../files/derp.html")
	if s.HasBeenSaved == false {
		t.Fatal("Should save file with storage")
	}
}

func TestLoadCallsStorage(t *testing.T) {
	s := new(MockStorage)
	treasury := treasury.NewTreasury(s)
	treasury.Load("sadfasdf")
	if s.HasBeenLoaded == false {
		t.Fatal("Should load file with storage")
	}
}
