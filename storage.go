package gorage

type Storage interface {
	Write(file File) error
	Read(hash string) (FileContent, error)
	Delete(hash string) error
}
