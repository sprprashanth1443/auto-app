package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const (
	users = "users"
)

// nolint:gochecknoinits
func init() {
	s := session.Copy()
	defer s.Close()

	c := s.DB(info.Database).C(users)

	indexes := []mgo.Index{
		{Key: []string{"address"}, Unique: true},
	}

	for i := 0; i < len(indexes); i++ {
		indexes[i].Background = true
		if err := c.EnsureIndex(indexes[i]); err != nil {
			log.Panicln(err)
		}
	}
}

type User struct {
	Address string `json:"address" bson:"address"`
	Handle  string `json:"handle" bson:"handle"`
}

func (u User) NewUser(address, handler string) User {
	return User{
		Address: address,
		Handle:  handler,
	}
}

func (u *User) FindOne(query bson.M) error {
	s := session.Copy()
	defer s.Close()

	c := s.DB(info.Database).C(users)
	if err := c.Find(query).One(u); err != nil {
		return err
	}

	return nil
}

func (u *User) Save() error {
	s := session.Copy()
	defer s.Close()

	c := s.DB(info.Database).C(users)
	if err := c.Insert(u); err != nil {
		return err
	}

	return nil
}

func (u *User) FindOneAndUpdate(query, update bson.M, upsert, remove, _new bool) error {
	s := session.Copy()
	defer s.Close()

	c := s.DB(info.Database).C(users)
	change := mgo.Change{
		Update:    update,
		Upsert:    upsert,
		Remove:    remove,
		ReturnNew: _new,
	}

	if _, err := c.Find(query).Apply(change, u); err != nil {
		return err
	}

	return nil
}

type Users []User

func (u *Users) Find(query, selector bson.M, sort []string, skip, limit int) error {
	s := session.Copy()
	defer s.Close()

	c := s.DB(info.Database).C(users)

	err := c.Find(query).Select(selector).Sort(sort...).Skip(skip).Limit(limit).All(u)
	if err != nil {
		return err
	}

	return nil
}
