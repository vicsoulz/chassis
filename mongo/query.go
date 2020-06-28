package mongo

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type RelationQuery struct {
	DB    *mgo.Database
	Limit int
	Nodes []*RelationQueryNode
}

type RelationQueryNode struct {
	DB             *mgo.Database
	Data           []map[string]interface{}
	Collection     string
	Condition      bson.M
	MergeCondition bson.M
	Project        bson.M
	ForeignField   string
	LocalField     string
	Limit          int
	Skip           int
	Sort           string
	As             string
}

func (q *RelationQuery) Query(result interface{}) (err error) {
	err = q.checkParam(result)
	if err != nil {
		return
	}

	q.AsyncNodesParam()

	for i := 0; i < len(q.Nodes); i++ {
		err = q.Nodes[i].Query()
		if err != nil {
			err = errors.New(fmt.Sprintf("collection: %s,condition: %v, query err: %s", q.Nodes[i].Collection, q.Nodes[i].Condition, err.Error()))
			return
		}

		if len(q.Nodes)-1 > i {
			var relValues []interface{}
			relValues, err = q.Nodes[i].getRelationValue()
			if err != nil {
				err = errors.New(fmt.Sprintf("node init foreign values: %s", err.Error()))
				return
			}

			if q.Nodes[i+1].Condition == nil {
				q.Nodes[i+1].Condition = bson.M{}
			}
			q.Nodes[i+1].Condition[q.Nodes[i+1].ForeignField] = bson.M{"$in": relValues}
		}
	}

	var data []map[string]interface{}
	data, err = q.mergeNodeData()
	if err != nil {
		return
	}

	return q.assignResult(result, data)
}

func (q *RelationQuery) AsyncNodesParam() {
	if q.DB == nil {
		q.DB = DB
	}

	for i := 0; i < len(q.Nodes); i++ {
		if q.Nodes[i].DB == nil {
			q.Nodes[i].DB = q.DB
		}

		if q.Limit > 0 && q.Nodes[i].Limit <= 0 {
			q.Nodes[i].Limit = q.Limit
		}
	}
}

func (q RelationQuery) mergeNodeData() (result []map[string]interface{}, err error) {
	if q.Nodes == nil || len(q.Nodes) == 0 {
		err = errors.New("node length is nil")
		return
	}
	if len(q.Nodes) == 1 {
		result = q.Nodes[0].Data
		return
	}

	result = q.Nodes[len(q.Nodes)-1].Data
	for i := len(q.Nodes) - 1; i > 0; i-- {
		result, err = q.merge(q.Nodes[i-1].Data, result, q.Nodes[i-1].LocalField, q.Nodes[i].ForeignField, q.Nodes[i].As)
		if err != nil {
			return
		}
	}
	return
}

func (q RelationQuery) checkParam(result interface{}) error {
	if q.Nodes == nil || len(q.Nodes) == 0 {
		return errors.New("query node is nil")
	}

	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr {
		return errors.New("result argument must be a slice address")
	}

	slicev := resultv.Elem()
	if slicev.Kind() != reflect.Slice {
		return errors.New("result argument must be a slice address")
	}

	return nil
}

func (q RelationQuery) assignResult(result interface{}, data []map[string]interface{}) error {
	resultv := reflect.ValueOf(result)
	slicev := resultv.Elem()
	slicev = slicev.Slice(0, slicev.Cap())

	elem := slicev.Type().Elem()
	i := 0
	for _, v := range data {
		byteData, err := bson.Marshal(v)
		if err != nil {
			return err
		}

		if slicev.Len() == i {
			elemp := reflect.New(elem)

			err = bson.Unmarshal(byteData, elemp.Interface())
			if err != nil {
				return err
			}

			slicev = reflect.Append(slicev, elemp.Elem())
			slicev = slicev.Slice(0, slicev.Cap())
		} else {
			err = bson.Unmarshal(byteData, slicev.Index(i).Addr().Interface())
			if err != nil {
				return err
			}
		}
		i++
	}

	resultv.Elem().Set(slicev.Slice(0, i))
	return nil
}

// 获取要关联字段的值list
func (n RelationQueryNode) getRelationValue() (result []interface{}, err error) {
	if n.LocalField == "" {
		err = errors.New("field is nil")
		return
	}

	result = make([]interface{}, 0, len(n.Data))
	for i := 0; i < len(n.Data); i++ {
		if _, ok := n.Data[i][n.LocalField]; !ok {
			err = errors.New(fmt.Sprintf("field: %v is not exist", n.LocalField))
			return
		}

		result = append(result, n.Data[i][n.LocalField])
	}

	return
}

// 关联到主数据上
func (q *RelationQuery) merge(left, right []map[string]interface{}, local string, foreign string, as string) (result []map[string]interface{}, err error) {
	if left == nil || len(left) == 0 || right == nil || len(right) == 0 {
		err = errors.New("merge data is nil")
		return
	}

	if local == "" || foreign == "" {
		err = errors.New("merge relation key is nil")
		return
	}

	// 创建关联数据和关联字段的mapping
	mapping := make(map[interface{}][]map[string]interface{}, len(right))
	for i := 0; i < len(right); i++ {
		var val interface{}
		var ok bool
		if val, ok = right[i][foreign]; !ok {
			err = errors.New(fmt.Sprintf("merge data is not exist field: %s", foreign))
			return
		}

		if mapping[val] == nil {
			mapping[val] = make([]map[string]interface{}, 0)
		}
		mapping[val] = append(mapping[val], right[i])
	}

	result = make([]map[string]interface{}, 0, len(right))
	// 关联
	for i := 0; i < len(left); i++ {
		var relValue interface{}
		var ok bool
		if relValue, ok = left[i][local]; !ok {
			err = errors.New(fmt.Sprintf("origin data is not exist field: %s", local))
			return
		}
		if _, ok := mapping[relValue]; ok {
			// 关联数据和原始数据按as字段合并
			for _, mVal := range mapping[relValue] {
				t := make(map[string]interface{})
				if as == "" {
					t = mVal
				} else {
					t[as] = mVal
				}
				for k, v := range left[i] {
					t[k] = v
				}

				if err != nil {
					return
				}
				result = append(result, t)
			}
		}
	}

	return
}

func (n *RelationQueryNode) Query() (err error) {
	if n.Collection == "" {
		err = errors.New("collection is nil")
		return
	}
	if n.Condition == nil {
		n.Condition = bson.M{}
	}

	qs := n.DB.C(n.Collection).Find(n.Condition)
	if n.Project != nil {
		qs = qs.Select(n.Project)
	}
	if n.Sort != "" {
		qs = qs.Sort(n.Sort)
	}
	if n.Limit != 0 {
		qs = qs.Limit(n.Limit)
	}
	if n.Skip != 0 {
		qs = qs.Skip(n.Skip)
	}

	err = qs.All(&n.Data)
	if err != nil {
		return
	}
	if len(n.Data) == 0 {
		err = mgo.ErrNotFound
		return
	}

	return
}
