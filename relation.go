package gorage

import "time"

type RelationRepo interface {
	Save(file Relation) (Relation, error)
	Load(hash string) (Relation, error)
	Delete(hash string) error
	HashExists(hash string) bool
}

/*
* A relations is a one to on mapping with an actual file (storage) and keeps all the relevant file information like size, mime-type and hash.
* The hash is identifier and representation of an file
 */
type Relation struct {
	CreatedAt time.Time
	DeletedAt *time.Time
	Hash      string
	MimeType  string
	Size      int
}
