package gorage

import "time"

type MetaRepo interface {
	Save(file Meta) (Meta, error)
	Load(id string) (Meta, error)
	Delete(id string) error
	HashExists(id string) bool
}

/*
 * Meta keeps all the relevant data of individual files. When two equal files are saved, both will have their own meta.
 */
type Meta struct {
	ID        string
	Hash      string
	CreatedAt time.Time
	DeletedAt *time.Time
	Name      string
	Context   interface{}
}
