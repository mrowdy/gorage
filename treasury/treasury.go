package treasury

import (
	"github.com/jinzhu/gorm"
)

type Treasury struct {
	Storage Storage
	Db      *gorm.DB
}

func NewTreasury(storage Storage, db *gorm.DB) *Treasury {
	treasury := new(Treasury)
	treasury.Storage = storage
	treasury.Db = db
	return treasury
}

func (t *Treasury) Save(file string) (*File, error) {
	f := NewFile("foo", []byte("Test"))
	err := t.Storage.Write(f)
	if err != nil {
		return f, err
	}

	dbFile := new(DbFile)
	dbFile.Hash = f.Hash()

	t.Db.Save(dbFile)

	return f, nil
}

func (t *Treasury) Load(file string) (*File, error) {
	loadedFile, err := t.Storage.Read(file)
	return &loadedFile, err
}

type DbFile struct {
	gorm.Model
	Hash string `gorm:"not null;unique"`
}

func (DbFile) TableName() string {
	return "file"
}
