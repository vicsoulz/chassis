package mysql

import (
	"github.com/jmoiron/sqlx"

	"github.com/vicsoulz/chassis/db"
)

var (
	Db         *sqlx.DB
	sourcePool = make(map[string]*sqlx.DB)
)

func Init(conf *db.Config) (err error) {
	sourcePool[conf.Name], err = db.New(conf, "mysql")
	if err != nil {
		return
	}

	if conf.Name == db.DefaultName {
		Db = sourcePool[conf.Name]
	}
	return
}

func Get(name string) *sqlx.DB {
	return sourcePool[name]
}
