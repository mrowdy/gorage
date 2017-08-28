package meta

import (
	"errors"
	"github.com/slemgrim/gorage"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Mongo struct {
	Collection *mgo.Collection
}

func (db Mongo) Save(m gorage.Meta) (gorage.Meta, error) {
	t := time.Now()

	m.ID = uuid.New().String()
	m.DeletedAt = nil
	m.CreatedAt = t
	err := db.Collection.Insert(&m)
	if err != nil {
		return m, err
	}

	return m, nil
}

func (db Mongo) Load(id string) (gorage.Meta, error) {
	meta := gorage.Meta{}

	err := db.Collection.Find(bson.M{"id": id, "deletedat": nil}).One(&meta)
	if err != nil {
		return meta, errors.New("Meta not found")
	}

	return meta, nil
}

func (db Mongo) Delete(id string) error {
	m := new(gorage.Meta)
	err := db.Collection.Find(bson.M{"id": id}).One(m)
	if err != nil {
		return err
	}

	t := time.Now()
	m.DeletedAt = &t

	err = db.Collection.Update(bson.M{"id": id, "deletedat": nil}, m)

	if err != nil {
		return err
	}

	return nil
}

func (db Mongo) HashExists(hash string) bool {

	r := new(gorage.Relation)
	err := db.Collection.Find(bson.M{"hash": hash, "deletedat": nil}).One(r)
	if err != nil {
		return false
	}

	return true
}
