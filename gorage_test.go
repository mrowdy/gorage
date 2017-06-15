package gorage

import (
	"errors"
	"testing"
)

type MockStorage struct {
	HasBeenWritten bool
	HasBeenRead    bool
	HasBeenDeleted bool
	DeletedId      string
	Content        FileContent
	File           File
}

func (s *MockStorage) Write(file File) error {
	s.HasBeenWritten = true
	s.File = file
	return nil
}

func (s *MockStorage) Read(hash string) (FileContent, error) {
	s.HasBeenRead = true
	return s.Content, nil
}

func (s *MockStorage) Delete(hash string) error {
	s.HasBeenDeleted = true
	s.DeletedId = hash
	return nil
}

type MockMetaRepo struct {
	HasBeenSaved   bool
	HasBeenLoaded  bool
	HasBeenDeleted bool
	SaveID         string
	LoadID         string
	DeleteID       string
	MetaExists     bool
	References     int
}

func (m *MockMetaRepo) Save(file Meta) (Meta, error) {
	m.HasBeenSaved = true
	return Meta{}, nil
}
func (m *MockMetaRepo) Load(id string) (Meta, error) {
	m.HasBeenLoaded = true

	if m.MetaExists == false {
		return Meta{}, errors.New("file doesn't exist")
	}

	return Meta{}, nil
}
func (m *MockMetaRepo) Delete(id string) error {
	m.HasBeenDeleted = true
	m.DeleteID = id
	return nil
}
func (m *MockMetaRepo) HashExists(id string) bool {
	return m.References > 0
}

type MockRelationRepo struct {
	HasBeenSaved   bool
	HasBeenLoaded  bool
	HasBeenDeleted bool
	HasBeenChecked bool
	SaveID         string
	LoadID         string
	DeleteID       string
	FileExists     bool
}

func (m *MockRelationRepo) Save(file Relation) (Relation, error) {
	m.HasBeenSaved = true
	return Relation{}, nil
}
func (m *MockRelationRepo) Load(hash string) (Relation, error) {
	m.HasBeenLoaded = true

	if m.FileExists == false {
		return Relation{}, errors.New("file doesn't exist")
	}

	return Relation{}, nil
}
func (m *MockRelationRepo) Delete(hash string) error {
	m.HasBeenDeleted = true
	return nil
}
func (m *MockRelationRepo) HashExists(hash string) bool {
	m.HasBeenChecked = true
	return m.FileExists
}

func TestStoreFile(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)
	gorage := NewGorage(s, r, m)
	gorage.Save("test", FileContent("foo"))
	if s.HasBeenWritten == false {
		t.Fatal("Should persist file")
	}
}

func TestStoreFileChecksIfFileExists(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	gorage := NewGorage(s, r, m)
	gorage.Save("test", FileContent("foo"))
	if r.HasBeenChecked != true {
		t.Fatal("Should check if file exists")
	}
}

func TestDontStoreExistingFiles(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	r.FileExists = true

	gorage := NewGorage(s, r, m)
	gorage.Save("test", FileContent("foo"))

	if s.HasBeenWritten == true {
		t.Fatal("Should store existing files")
	}
}

func TestSaveWholeFileContent(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	r.FileExists = true

	name := "foo"
	body := []byte("test")

	gorage := NewGorage(s, r, m)
	f, _ := gorage.Save(name, body)

	if f.Name != name {
		t.Fatal("name not set")
	}

	if !testEq(f.Content, body) {
		t.Fatal("content not set")
	}
}

func TestSaveFileToRelationsRepo(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	name := "foo"
	body := []byte("test")

	gorage := NewGorage(s, r, m)
	gorage.Save(name, body)

	if r.HasBeenSaved != true {
		t.Fatal("Should save to repo")
	}
}

func TestDontSaveExistingFiles(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	r.FileExists = true

	name := "foo"
	body := []byte("test")

	gorage := NewGorage(s, r, m)
	gorage.Save(name, body)

	if r.HasBeenSaved == true {
		t.Fatal("Should not save to repo")
	}
}

func TestLoadsFromStorageIfFileExists(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	r.FileExists = true

	name := "foo"
	body := []byte("test")

	gorage := NewGorage(s, r, m)
	gorage.Save(name, body)

	if r.HasBeenLoaded != true {
		t.Fatal("Should have been loaded")
	}
}

func TestCreatesMetaForNewFiles(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	name := "foo"
	body := []byte("test")

	gorage := NewGorage(s, r, m)
	gorage.Save(name, body)

	if m.HasBeenSaved != true {
		t.Fatal("Should have been saved to meta")
	}
}

func TestCreatesMetaForExistingFiles(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	r.FileExists = true

	name := "foo"
	body := []byte("test")

	gorage := NewGorage(s, r, m)
	gorage.Save(name, body)

	if m.HasBeenSaved != true {
		t.Fatal("Should have been saved to meta")
	}
}

func TestDeleteChecksIfFileExistsInMeta(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	m.MetaExists = true
	gorage := NewGorage(s, r, m)
	gorage.Delete("test")

	if m.HasBeenLoaded != true {
		t.Fatal("should load file from meta")
	}
}

func TestDeleteReturnsErrorIfFileDoesntExist(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	gorage := NewGorage(s, r, m)
	err := gorage.Delete("test")

	if err == nil {
		t.Fatal("should return error")
	}
}

func TestItLoadsFileFromRelations(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	m.MetaExists = true
	gorage := NewGorage(s, r, m)
	gorage.Delete("test")

	if r.HasBeenLoaded == false {
		t.Fatal("should have been loaded")
	}
}

func TestItShouldntLoadFromRelationIfNoMetaExists(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	m.MetaExists = false
	gorage := NewGorage(s, r, m)
	gorage.Delete("test")

	if r.HasBeenLoaded == true {
		t.Fatal("should not have been loaded")
	}
}

func TestItShouldReturnErrorIfRelationDoesntExist(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	m.MetaExists = true
	r.FileExists = false
	gorage := NewGorage(s, r, m)
	err := gorage.Delete("test")

	if err == nil {
		t.Fatal("should return error")
	}
}

func TestItShouldDeleteExistingFilesFromMeta(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	m.MetaExists = true
	r.FileExists = true
	gorage := NewGorage(s, r, m)
	gorage.Delete("test")

	if m.HasBeenDeleted == false || m.DeleteID != "test" {
		t.Fatal("should have been deleted")
	}
}

func TestItShouldDeleteRelationIfNotExists(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	m.MetaExists = true
	m.References = 0
	r.FileExists = true

	gorage := NewGorage(s, r, m)
	gorage.Delete("test")

	if r.HasBeenDeleted == false {
		t.Fatal("should have been deleted")
	}
}

func TestItShouldDeleteDeleteFileIfLastRefrence(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	m.MetaExists = true
	m.References = 0
	r.FileExists = true

	gorage := NewGorage(s, r, m)
	gorage.Delete("test")

	if s.HasBeenDeleted == false {
		t.Fatal("should have been deleted")
	}
}

func TestItShouldNotDeleteRelationIfMoreReferencesExist(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	m.MetaExists = true
	m.References = 2
	r.FileExists = true

	gorage := NewGorage(s, r, m)
	gorage.Delete("test")

	if r.HasBeenDeleted == true {
		t.Fatal("should have been deleted")
	}
}

func TestItShouldNotDeleteDeleteFileIfNotLastRefrence(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	m.MetaExists = true
	m.References = 2
	r.FileExists = true

	gorage := NewGorage(s, r, m)
	gorage.Delete("test")

	if s.HasBeenDeleted == true {
		t.Fatal("should not have been deleted")
	}
}

func TestItShouldReturnErrorIfMetaDoesNotExist(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	m.MetaExists = false
	r.FileExists = false

	gorage := NewGorage(s, r, m)
	_, err := gorage.Load("1234")

	if err == nil {
		t.Fatal("should return error")
	}
}

func TestItShouldNotReturnErrorIfMetaExist(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	m.MetaExists = true
	r.FileExists = true

	gorage := NewGorage(s, r, m)
	_, err := gorage.Load("1234")

	if err != nil {
		t.Fatal("should not return error")
	}
}

func TestItShouldReturnErrorWhenRelationNotExist(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	m.MetaExists = true
	r.FileExists = false

	gorage := NewGorage(s, r, m)
	_, err := gorage.Load("1234")

	if err == nil {
		t.Fatal("should return error")
	}
}

func TestShouldLoadFromStorageIfFileExists(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	m.MetaExists = true
	r.FileExists = true

	gorage := NewGorage(s, r, m)
	gorage.Load("1234")

	if s.HasBeenRead == false {
		t.Fatal("should read from storage")
	}
}

func TestShouldNotLoadFromStorageIfFileDoesNotExists(t *testing.T) {
	s := new(MockStorage)
	m := new(MockMetaRepo)
	r := new(MockRelationRepo)

	m.MetaExists = false
	r.FileExists = false

	gorage := NewGorage(s, r, m)
	gorage.Load("1234")

	if s.HasBeenRead == true {
		t.Fatal("should not read from storage")
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
