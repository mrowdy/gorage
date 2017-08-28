package relation

import (
	"errors"
	"github.com/slemgrim/gorage"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Mongo struct {
	Collection *mgo.Collection
}

func (db Mongo) Save(r gorage.Relation) (gorage.Relation, error) {
	relation, found := db.getRelation(r.Hash)

	if !found {

		t := time.Now()
		r.DeletedAt = nil
		r.CreatedAt = t
		err := db.Collection.Insert(&r)

		if err != nil {
			return r, err
		}
	} else {
		r = *relation
	}

	return r, nil
}

func (db Mongo) Load(hash string) (gorage.Relation, error) {
	r := gorage.Relation{}

	err := db.Collection.Find(bson.M{"hash": hash, "deletedat": nil}).One(&r)
	if err != nil {
		return r, errors.New("Relation not found")
	}

	return r, nil
}

func (db Mongo) Delete(hash string) error {
	r := new(gorage.Relation)
	err := db.Collection.Find(bson.M{"hash": hash}).One(r)
	if err != nil {
		return err
	}

	t := time.Now()
	r.DeletedAt = &t

	err = db.Collection.Update(bson.M{"hash": hash, "deletedat": nil}, r)
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

func (db Mongo) getRelation(hash string) (*gorage.Relation, bool) {
	r := new(gorage.Relation)

	err := db.Collection.Find(bson.M{"hash": hash, "deletedat": nil}).One(r)
	if err != nil {
		return r, false
	}

	return r, true
}
