package sort

import (
  "testing"
  "reflect"
)

func TestSortByValue(t *testing.T) {
  m := map[string]int{
    "foo": 5,
    "bar": 2,
    "baz": 13,
  }
  keys := ByValue(m, func(a, b reflect.Value) bool {
    return a.Int() > b.Int()
  }).Interface().([]string)
  if keys[0] != "baz" || keys[1] != "foo" || keys[2] != "bar" {
    t.Fail()
  }
}
