package gorage

type Persistance interface {
	Save(file *File) error
	Load(id string) (File, error)
	HashExists(hash string) bool
}
