package jsonutils

import (
	"encoding/json"
	"fmt"
	"testing"
)

var (
	str = []byte(`{"checked": true, "name": "standard", "created": "2018-04-09T23:00:00Z", "price": {"apple": 3.5, "banana": 7.99}, "fruit": ["apple", "banana", "orange"], "owner": null, "ref": 999}`)
)

func TestLoadsKey(t *testing.T) {
	var obj JsonObject
	if json.Unmarshal(str, &obj) != nil {
		t.Error("fail to load")
		return
	}
	checked, err := obj.GetBool("checked")
	fmt.Printf("checked:%v, err:%v\n", checked, err)
	fruit, err := obj.GetJsonArray("fruit")
	fmt.Printf("Fruit:%v, err:%v\n", fruit, err)
	fruitArray, err := fruit.ToStringArray()
	fmt.Printf("FruitArray:%v, err:%v\n", fruitArray, err)
	price, err := obj.GetJsonObject("price")
	fmt.Printf("price:%v, err:%v\n", price, err)
	applePrice, err := price.GetFloat64("apple")
	fmt.Printf("apple price:%v, err:%v\n", applePrice, err)
	ref, err := obj.GetFloat64("ref")
	fmt.Printf("ref:%v, err:%v\n", int32(ref), err)
}
