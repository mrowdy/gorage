package meta

import (
	"errors"
	"time"

	"encoding/json"

	"github.com/slemgrim/gorage"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var tableName = "gorage_relations"

type Sql struct {
	Db *gorm.DB
}

func NewSql(table string, db *gorm.DB) *Sql {
	tableName = table
	db.AutoMigrate(&metaTable{})
	Db := new(Sql)
	Db.Db = db
	return Db
}

func (db Sql) Save(m gorage.Meta) (gorage.Meta, error) {
	meta := new(metaTable)
	meta.Name = m.Name
	meta.Hash = m.Hash
	meta.Context, _ = json.Marshal(m.Context)

	result := db.Db.Create(meta)
	if result.Error != nil {
		return m, result.Error
	}

	m.ID = meta.ID
	m.CreatedAt = meta.CreatedAt
	m.DeletedAt = meta.DeletedAt

	return m, nil
}

func (db Sql) Load(id string) (gorage.Meta, error) {
	m := metaTable{}
	meta := gorage.Meta{}

	notFound := db.Db.Where(metaTable{
		ID:        id,
		DeletedAt: nil,
	}, id).First(&m).RecordNotFound()

	if !notFound {

		var context interface{}
		err := json.Unmarshal(m.Context, &context)

		if err != nil {
			panic(err)
			return meta, err
		}

		meta.ID = m.ID
		meta.Name = m.Name
		meta.Hash = m.Hash
		meta.CreatedAt = m.CreatedAt
		meta.DeletedAt = m.DeletedAt
		meta.Context = context

		return meta, nil
	}

	return meta, errors.New("Meta not found")
}

func (db Sql) Delete(id string) error {
	result := db.Db.Where(metaTable{
		ID:        id,
		DeletedAt: nil,
	}, id).Delete(&metaTable{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db Sql) HashExists(hash string) bool {
	count := 0
	db.Db.Model(&metaTable{}).Where(&metaTable{
		Hash:      hash,
		DeletedAt: nil,
	}).Count(&count)

	return count > 0
}

func (fc *metaTable) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New().String())
	return nil
}

func (metaTable) TableName() string {
	return tableName
}

type metaTable struct {
	ID        string `gorm:"primary_key"`
	Hash      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Name      string
	Context   []byte
}
