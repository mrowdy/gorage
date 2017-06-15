package gorage

/*
* Inteface for handling files. Possible implementations would be Filesystem, Memcache, Redis, Database, etc
 */
type Storage interface {
	Write(file File) error
	Read(hash string) (FileContent, error)
	Delete(hash string) error
}
