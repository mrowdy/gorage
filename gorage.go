package gorage

import (
	"errors"
	"net/http"
	"time"
)

type Gorage struct {
	Storage      Storage
	MetaRepo     MetaRepo
	RelationRepo RelationRepo
}

func NewGorage(storage Storage, relationRepo RelationRepo, metaRepo MetaRepo) *Gorage {
	gorage := new(Gorage)
	gorage.Storage = storage
	gorage.MetaRepo = metaRepo
	gorage.RelationRepo = relationRepo
	return gorage
}

func (g Gorage) Save(name string, body []byte) (File, error) {
	content := FileContent(body)
	f := File{}
	f.Name = name
	f.Content = content
	f.Hash = content.CalculateHash()
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
	m, err := g.MetaRepo.Save(m)

	if err != nil {
		return File{}, err
	}

	f.ID = m.ID
	f.UploadedAt = m.CreatedAt
	f.DeletedAt = m.DeletedAt

	return f, nil

}

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
		DeletedAt:  m.DeletedAt,
		UploadedAt: m.CreatedAt,
	}

	return f, nil
}

func getMymeType(f FileContent) string {
	return http.DetectContentType(f)
}

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

type MetaRepo interface {
	Save(file Meta) (Meta, error)
	Load(id string) (Meta, error)
	Delete(id string) error
	HashExists(id string) bool
}

type Meta struct {
	ID        string
	Hash      string
	CreatedAt time.Time
	DeletedAt *time.Time
	Name      string
}

type RelationRepo interface {
	Save(file Relation) (Relation, error)
	Load(hash string) (Relation, error)
	Delete(hash string) error
	HashExists(hash string) bool
}

type Relation struct {
	CreatedAt time.Time
	DeletedAt *time.Time
	Hash      string
	MimeType  string
	Size      int
}
