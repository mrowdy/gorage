package gorage

type Storage interface {
	Write(file *File) error
	Read(hash string) (File, error)
}
