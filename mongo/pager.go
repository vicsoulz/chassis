package mongo

import "github.com/globalsign/mgo"

type Pager struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

func (p Pager) Query(query *mgo.Query) *mgo.Query {
	if p.Page > 0 {
		p.Page = p.Page - 1
	}

	if p.Limit <= 10 {
		p.Limit = 10
	}
	return query.Limit(p.Limit).Skip(p.Page * p.Limit)
}
