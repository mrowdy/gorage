package gorage

import (
	"errors"
	"testing"
)

type MockStorage struct {
	HasBeenWritten bool
	HasBeenRead    bool
	Content        FileContent
}

func (s *MockStorage) Write(file *File) error {
	s.HasBeenWritten = true
	return nil
}

func (s *MockStorage) Read(id string) (FileContent, error) {
	s.HasBeenRead = true
	return s.Content, nil
}

type MockPersistance struct {
	FilePersisted bool
	FileExists    bool
	HasBeenSaved  bool
	HasBeenLoaded bool
	Hash          bool
	File          *File
}

func (p *MockPersistance) Save(file *File) error {
	p.HasBeenSaved = true
	return nil
}

func (p *MockPersistance) Load(id string) (File, error) {
	p.HasBeenLoaded = true

	if p.FilePersisted == false {
		return *new(File), errors.New("File not found")
	}

	return *p.File, nil
}

func (p *MockPersistance) HashExists(hash string) bool {
	return p.FileExists
}

func TestSaveWritesToStorage(t *testing.T) {
	s := new(MockStorage)
	p := new(MockPersistance)
	gorage := NewGorage(s, p)
	gorage.Save("./files/derp.html")
	if s.HasBeenWritten == false {
		t.Fatal("Should persist file")
	}
}

func TestDontPersistFiles(t *testing.T) {
	s := new(MockStorage)
	p := new(MockPersistance)
	gorage := NewGorage(s, p)
	gorage.Save("./files/derp.html")
	if p.HasBeenSaved == false {
		t.Fatal("Should save file with storage")
	}
}

func TestNotSavingDuplicateFiles(t *testing.T) {
	s := new(MockStorage)
	p := new(MockPersistance)
	p.FileExists = true

	gorage := NewGorage(s, p)
	gorage.Save("./files/derp.html")

	if s.HasBeenWritten == true {
		t.Fatal("dont write duplicate files")
	}
}

func TestErrorWhenFileIsNotPersisted(t *testing.T) {
	s := new(MockStorage)
	p := new(MockPersistance)
	p.FilePersisted = false
	gorage := NewGorage(s, p)

	_, err := gorage.Load("13123")
	if err == nil {
		t.Fatal("it should return err when files are not persisted")
	}
}

func TestStorageShouldNotReadUnpersistedFiles(t *testing.T) {
	s := new(MockStorage)
	p := new(MockPersistance)
	p.FilePersisted = false
	gorage := NewGorage(s, p)

	gorage.Load("13123")
	if s.HasBeenRead == true {
		t.Fatal("it shouldn't try to read unpersisted fils")
	}
}

func TestReturnsDataFromPersistanceLayer(t *testing.T) {
	s := new(MockStorage)
	p := new(MockPersistance)

	f := new(File)
	f.ID = "1234"

	p.FilePersisted = true
	p.File = f

	gorage := NewGorage(s, p)

	file, err := gorage.Load("13123")

	if err != nil {
		t.Fatal("unexpected error")
	}

	if f.ID != file.ID {
		t.Fatal("file not returned from persistance layer")
	}
}

func TestLoadsExistingFilesFromStorage(t *testing.T) {
	s := new(MockStorage)
	p := new(MockPersistance)
	f := new(File)

	p.FileExists = true
	p.FilePersisted = true
	p.File = f
	gorage := NewGorage(s, p)
	gorage.Load("sadfasdf")
	if s.HasBeenRead == false {
		t.Fatal("id not set")

	}
}

func TestContentwillBeSet(t *testing.T) {
	s := new(MockStorage)
	s.Content = []byte("Content")

	p := new(MockPersistance)

	f := new(File)
	f.ID = "1234"

	p.FilePersisted = true
	p.File = f

	gorage := NewGorage(s, p)

	file, err := gorage.Load("13123")

	if err != nil {
		t.Fatal("unexpected error")
	}

	if !testEq(s.Content, file.Content) {
		t.Fatal("content not set")
	}
}

func testEq(a, b []byte) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
