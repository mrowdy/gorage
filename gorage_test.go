package gorage

import (
	"testing"
)

type MockStorage struct {
	HasBeenSaved  bool
	HasBeenLoaded bool
}

func (s *MockStorage) Write(file *File) error {
	s.HasBeenSaved = true
	return nil
}

func (s *MockStorage) Read(id string) (File, error) {
	s.HasBeenLoaded = true
	return *new(File), nil
}

func TestSaveCallsStorage(t *testing.T) {
	s := new(MockStorage)
	gorage := NewGorage(s)
	gorage.Save("./files/derp.html")
	if s.HasBeenSaved == false {
		t.Fatal("Should save file with storage")
	}
}

func TestLoadCallsStorage(t *testing.T) {
	s := new(MockStorage)
	gorage := NewGorage(s)
	gorage.Load("sadfasdf")
	if s.HasBeenLoaded == false {
		t.Fatal("Should load file with storage")
	}
}
