package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const (
	tweet = "tweet"
)

func init() {
	s := session.Copy()
	defer s.Close()

	c := s.DB(info.Database).C(tweet)

	indexes := []mgo.Index{
		{Key: []string{"url", "address"}, Unique: true},
	}

	for i := 0; i < len(indexes); i++ {
		indexes[i].Background = true
		if err := c.EnsureIndex(indexes[i]); err != nil {
			log.Panicln(err)
		}
	}
}

type Tweet struct {
	Url     string `json:"url" bson:"url"`
	AddedBy string `json:"added_by" bson:"added_by"`
}

func (t Tweet) NewTweet(tweet, address string) Tweet {
	return Tweet{
		Url:     tweet,
		AddedBy: address,
	}
}

func (t *Tweet) FindOne(query bson.M) error {
	s := session.Copy()
	defer s.Close()

	c := s.DB(info.Database).C(tweet)
	if err := c.Find(query).One(t); err != nil {
		return err
	}

	return nil
}

func (t *Tweet) Save() error {
	s := session.Copy()
	defer s.Close()

	c := s.DB(info.Database).C(tweet)
	if err := c.Insert(t); err != nil {
		return err
	}

	return nil
}

func (t *Tweet) FindOneAndUpdate(query, update bson.M, upsert, remove, _new bool) error {
	s := session.Copy()
	defer s.Close()

	c := s.DB(info.Database).C(tweet)
	change := mgo.Change{
		Update:    update,
		Upsert:    upsert,
		Remove:    remove,
		ReturnNew: _new,
	}

	if _, err := c.Find(query).Apply(change, t); err != nil {
		return err
	}

	return nil
}

type Tweets []Tweet

func (t *Tweets) Find(query, selector bson.M, sort []string, skip, limit int) error {
	s := session.Copy()
	defer s.Close()

	c := s.DB(info.Database).C(tweet)

	err := c.Find(query).Select(selector).Sort(sort...).Skip(skip).Limit(limit).All(t)
	if err != nil {
		return err
	}

	return nil
}
