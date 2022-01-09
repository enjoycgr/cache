package app

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("app hit key1=1234 failed")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("app miss key2 failed")
	}
}

func TestOnEnvicted(t *testing.T) {
	keys := make([]string, 0)
	lru := New(
		int64(10),
		func(key string, value Value) {
			keys = append(keys, key)
		},
	)
	lru.Add("key1", String("123456"))
	lru.Add("k1", String("k1"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))

	expect := []string{"key1", "k1"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
