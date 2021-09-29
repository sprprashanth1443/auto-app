package models

import (
	"github.com/auto-app/backend/config"
	"gopkg.in/mgo.v2"
	"time"
)

// nolint:gochecknoglobals
var (
	session *mgo.Session
	info    = &mgo.DialInfo{
		Addrs: []string{
			config.DbAddr,
		},
		Database: config.DbName,
		Timeout:  15 * time.Second,
	}

	_ = func() (err error) {
		session, err = mgo.DialWithInfo(info)
		if err != nil {
			return err
		}

		return nil
	}()
)
