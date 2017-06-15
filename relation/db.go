package relation

import (
	"time"

	"errors"

	"github.com/Slemgrim/gorage"
	"github.com/jinzhu/gorm"
)

type Db struct {
	Db *gorm.DB
}

var tableName = "gorage_mate"

func NewDb(table string, db *gorm.DB) *Db {
	tableName = table
	db.AutoMigrate(&relationTable{})
	Db := new(Db)
	Db.Db = db
	return Db
}

func (db Db) Save(r gorage.Relation) (gorage.Relation, error) {

	relation, found := db.getRelation(r.Hash)

	if !found {
		relation = new(relationTable)
		relation.Hash = r.Hash
		relation.MimeType = r.MimeType
		relation.Size = r.Size

		result := db.Db.Create(relation)
		if result.Error != nil {
			return gorage.Relation{}, result.Error
		}

		r.Hash = relation.Hash
		r.CreatedAt = relation.CreatedAt
	}

	return r, nil
}

func (db Db) Load(hash string) (gorage.Relation, error) {
	r := relationTable{}
	relation := gorage.Relation{}
	notFound := db.Db.Where(relationTable{
		Hash:      hash,
		DeletedAt: nil,
	}, hash).First(&r).RecordNotFound()

	if !notFound {
		relation.Hash = r.Hash
		relation.MimeType = r.MimeType
		relation.Size = r.Size
		relation.CreatedAt = r.CreatedAt
		relation.DeletedAt = r.DeletedAt
		return relation, nil
	}

	return gorage.Relation{}, errors.New("Relation not found")

}

func (db Db) Delete(hash string) error {
	result := db.Db.Where(relationTable{
		Hash:      hash,
		DeletedAt: nil,
	}, hash).Delete(&relationTable{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db Db) HashExists(hash string) bool {
	relation := relationTable{}
	result := db.Db.Where(&relationTable{Hash: hash}).First(&relation)

	if result.Error != nil {
		return false
	}

	return true
}

func (db Db) getRelation(hash string) (*relationTable, bool) {
	relation := new(relationTable)
	found := false

	result := db.Db.Where(&relationTable{Hash: hash}).First(&relation)
	if result.Error == nil {
		found = true
	}

	return relation, found
}

type relationTable struct {
	CreatedAt time.Time
	DeletedAt *time.Time
	Hash      string
	MimeType  string
	Size      int
}

func (relationTable) TableName() string {
	return tableName
}
