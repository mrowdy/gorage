package gorage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/http"
)

/*
* Gorage is an abstraction for file IO which can save a lot of disk space by preventing the storage of same files.
* It achieves this by splitting the file content and meta data.
 */
type Gorage struct {
	Storage      Storage
	MetaRepo     MetaRepo
	RelationRepo RelationRepo
}

/*
* Create a new instance of Gorage with given drivers
 */
func NewGorage(storage Storage, relationRepo RelationRepo, metaRepo MetaRepo) *Gorage {
	gorage := new(Gorage)
	gorage.Storage = storage
	gorage.MetaRepo = metaRepo
	gorage.RelationRepo = relationRepo
	return gorage
}

/*
* Save a file with its body and content. It returns a File object with all the relevant meta data and an error
* in case the file could not be saved
 */
func (g Gorage) Save(name string, body []byte, context interface{}) (File, error) {
	content := FileContent(body)
	f := File{}
	f.Name = name
	f.Content = content
	f.Hash = calculateHash(content)
	f.MimeType = getMymeType(content)
	f.Size = len(content)

	fileExists := g.RelationRepo.HashExists(f.Hash)
	r := Relation{}

	if !fileExists {
		err := g.Storage.Write(f)

		if err != nil {
			return File{}, err
		}

		r.Hash = f.Hash
		r.MimeType = f.MimeType
		r.Size = f.Size
		savedRelation, err := g.RelationRepo.Save(r)

		if err != nil {
			return File{}, err
		}

		r = savedRelation
	} else {
		re, err := g.RelationRepo.Load(f.Hash)

		if err != nil {
			return File{}, err
		}

		r = re
	}

	m := Meta{}
	m.Name = f.Name
	m.Hash = r.Hash
	m.Context = context
	m, err := g.MetaRepo.Save(m)

	if err != nil {
		return File{}, err
	}

	f.ID = m.ID
	f.UploadedAt = m.CreatedAt
	f.DeletedAt = m.DeletedAt

	return f, nil

}

/*
* Load a file by an ID. It returns a File object with all the relevant meta data and an error
* in case the file could not be saved
 */
func (g Gorage) Load(id string) (File, error) {

	m, mErr := g.MetaRepo.Load(id)

	if mErr != nil {
		return File{}, errors.New("File not found")
	}

	r, rErr := g.RelationRepo.Load(m.Hash)

	if rErr != nil {
		return File{}, errors.New("No related file found")
	}

	content, sErr := g.Storage.Read(r.Hash)

	if sErr != nil {
		return File{}, sErr
	}

	f := File{
		ID:         m.ID,
		Content:    content,
		Name:       m.Name,
		MimeType:   r.MimeType,
		Size:       r.Size,
		Context:    m.Context,
		DeletedAt:  m.DeletedAt,
		UploadedAt: m.CreatedAt,
	}

	return f, nil
}

/*
* Delete a file by its ID. Returns an error if file can't be deleted or was already deleted
 */
func (g Gorage) Delete(id string) error {
	m, mErr := g.MetaRepo.Load(id)
	if mErr != nil {
		return errors.New("File not Found")
	}

	r, rErr := g.RelationRepo.Load(m.Hash)
	if rErr != nil {
		return errors.New("File not Found")
	}

	g.MetaRepo.Delete(id)

	if !g.MetaRepo.HashExists(r.Hash) {
		err := g.RelationRepo.Delete(r.Hash)
		if err != nil {
			return err
		}
		err = g.Storage.Delete(r.Hash)
		if err != nil {
			return err
		}
	}

	return nil
}

func getMymeType(f FileContent) string {
	return http.DetectContentType(f)
}

func calculateHash(f FileContent) string {
	h := sha1.New()
	h.Write(f)
	return fmt.Sprintf("%x", h.Sum(nil))
}
