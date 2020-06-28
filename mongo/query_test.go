package mongo

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/globalsign/mgo/bson"
)

var (
	localCollect     string
	relationCollect  string
	relationCollect1 string
)

func testInit() {
	if err := InitDefault(); err != nil {
		panic(err)
	}

	initCollectName()
	initData()
}

func initCollectName() {
	if localCollect == "" {
		randNum := math.Round(100)
		localCollect = fmt.Sprintf("test_relationquery_local_%d", int(randNum))
		relationCollect = fmt.Sprintf("test_relationquery_relation_%d", int(randNum))
		relationCollect1 = fmt.Sprintf("test_relationquery_relation1_%d", int(randNum))
	}
}

func initData() {
	_, err := C(localCollect).RemoveAll(bson.M{})
	if err != nil {
		panic(err)
	}
	_, err = C(relationCollect).RemoveAll(bson.M{})
	if err != nil {
		panic(err)
	}

	_, err = C(relationCollect1).RemoveAll(bson.M{})
	if err != nil {
		panic(err)
	}

	b := C(localCollect).Bulk()
	b.Insert(map[string]interface{}{
		"order_id": 1,
		"order_sn": "test_sn1",
	})
	b.Insert(map[string]interface{}{
		"order_id": 2,
		"order_sn": "test_sn2",
	})
	b.Insert(map[string]interface{}{
		"order_id": 3,
		"order_sn": "test_sn3",
	})

	_, err = b.Run()
	if err != nil {
		panic(err)
	}

	b = C(relationCollect).Bulk()
	b.Insert(map[string]interface{}{
		"order_id":      1,
		"order_item_id": 1,
		"order_item_sn": "test_item_sn1",
	})

	b.Insert(map[string]interface{}{
		"order_id":      1,
		"order_item_id": 2,
		"order_item_sn": "test_item_sn2",
	})

	b.Insert(map[string]interface{}{
		"order_id":      2,
		"order_item_id": 3,
		"order_item_sn": "test_item_sn3",
	})

	_, err = b.Run()
	if err != nil {
		panic(err)
	}

	b = C(relationCollect1).Bulk()
	b.Insert(map[string]interface{}{
		"order_item_id":   1,
		"package_item_id": 1,
	})

	b.Insert(map[string]interface{}{
		"order_item_id":   2,
		"package_item_id": 2,
	})

	b.Insert(map[string]interface{}{
		"order_item_id":   3,
		"package_item_id": 3,
	})

	_, err = b.Run()
	if err != nil {
		panic(err)
	}
}

func TestRelationQuery_Query(t *testing.T) {
	testInit()
	rq := RelationQuery{
		Nodes: []*RelationQueryNode{
			{
				Collection: localCollect,
				LocalField: "order_id",
				Project:    bson.M{"_id": 0},
			},
			{
				Collection:   relationCollect,
				LocalField:   "order_item_id",
				ForeignField: "order_id",
				Project:      bson.M{"_id": 0},
				As:           relationCollect,
			},
			{
				Collection:   relationCollect1,
				ForeignField: "order_item_id",
				Project:      bson.M{"_id": 0},
				As:           relationCollect1,
			},
		},
	}
	var res []map[string]interface{}
	err := rq.Query(&res)
	if err != nil {
		t.Error(err.Error())
	}

	compare := []map[string]interface{}{
		map[string]interface{}{
			"order_id": 1,
			"order_sn": "test_sn1",
			relationCollect: map[string]interface{}{
				"order_id":      1,
				"order_item_id": 1,
				"order_item_sn": "test_item_sn1",
				relationCollect1: map[string]interface{}{
					"order_item_id":   1,
					"package_item_id": 1,
				},
			},
		},
		map[string]interface{}{
			"order_id": 1,
			"order_sn": "test_sn1",
			relationCollect: map[string]interface{}{
				"order_id":      1,
				"order_item_id": 2,
				"order_item_sn": "test_item_sn2",
				relationCollect1: map[string]interface{}{
					"order_item_id":   2,
					"package_item_id": 2,
				},
			},
		},
		map[string]interface{}{
			"order_id": 2,
			"order_sn": "test_sn2",
			relationCollect: map[string]interface{}{
				"order_id":      2,
				"order_item_id": 3,
				"order_item_sn": "test_item_sn3",
				relationCollect1: map[string]interface{}{
					"order_item_id":   3,
					"package_item_id": 3,
				},
			},
		},
	}

	if !reflect.DeepEqual(compare, res) {
		fmt.Println(res)
		t.Error("compare res no equal")
		return
	}

	rq = RelationQuery{
		Nodes: []*RelationQueryNode{
			{
				Collection: localCollect,
				LocalField: "order_id",
				Project:    bson.M{"_id": 0},
			},
		},
	}
	err = rq.Query(&res)
	if err != nil {
		t.Error(err.Error())
	}

	compare = []map[string]interface{}{
		map[string]interface{}{
			"order_id": 1,
			"order_sn": "test_sn1",
		},
		map[string]interface{}{
			"order_id": 2,
			"order_sn": "test_sn2",
		},
		map[string]interface{}{
			"order_id": 3,
			"order_sn": "test_sn3",
		},
	}

	if !reflect.DeepEqual(compare, res) {
		fmt.Println(res)
		t.Error("compare res no equal")
		return
	}

}

func TestRelationQueryNode_Query(t *testing.T) {
	testInit()
	n := RelationQueryNode{
		Collection: localCollect,
		Project: bson.M{
			"order_sn": 1,
			"_id":      0,
		},
		Condition: bson.M{
			"order_id": 1,
		},
	}

	err := n.Query()
	if err != nil {
		t.Error(err)
		return
	}

	compare := []map[string]interface{}{
		map[string]interface{}{
			"order_sn": "test_sn1",
		},
	}

	if !reflect.DeepEqual(compare, n.Data) {
		fmt.Println(n.Data)
		t.Error("compare res no equal")
	}
}
