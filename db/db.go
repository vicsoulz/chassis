package db

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Name            string
	Uri             string
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

const (
	DefaultMaxIdleConns = 10
	DefaultMaxOpenConns = 20
	DefaultName         = "default"
)

func New(conf *Config, driverName string) (*sqlx.DB, error) {
	err := checkConfig(conf)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open(driverName, conf.Uri)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	if conf.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(conf.ConnMaxLifetime)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func checkConfig(c *Config) error {
	if c.Name == "" {
		c.Name = DefaultName
	}

	if c.MaxIdleConns == 0 {
		c.ConnMaxLifetime = DefaultMaxIdleConns
	}

	if c.MaxIdleConns == 0 {
		c.MaxOpenConns = DefaultMaxOpenConns
	}

	return nil
}

func SelectToMap(db *sqlx.DB, sql string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)
	for {
		if !rows.Next() {
			break
		}
		tmp := make(map[string]interface{})

		err = rows.MapScan(tmp)
		if err != nil {
			return nil, err
		}

		result = append(result, tmp)
	}

	return result, nil
}
