package gorage

import (
	"io/ioutil"
	"net/http"
	"path"
)

type Gorage struct {
	Storage     Storage
	Persistance Persistance
}

func NewGorage(storage Storage, persistance Persistance) *Gorage {
	gorage := new(Gorage)
	gorage.Storage = storage
	gorage.Persistance = persistance
	return gorage
}

func (t Gorage) Save(file string) (*File, error) {
	content, _ := ioutil.ReadFile(file)
	f := new(File)
	f.Name = path.Base(file)
	f.Content = content
	f.Hash = f.CalculateHash()
	f.MimeType = getMymeType(content)
	f.Size = len(content)

	if !t.Persistance.HashExists(f.Hash) {
		err := t.Storage.Write(f)
		if err != nil {
			return f, err
		}
	}

	t.Persistance.Save(f)
	return f, nil
}

func (t Gorage) Load(id string) (*File, error) {

	file, err := t.Persistance.Load(id)
	if err != nil {
		return new(File), err
	}

	content, err := t.Storage.Read(file.Hash)
	file.Content = content
	return &file, err
}

func getMymeType(f FileContent) string {
	return http.DetectContentType(f)
}
