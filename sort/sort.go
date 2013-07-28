package sort

import (
  "reflect"
  "sort"
)

type sortByValue struct {
  m interface{}
  key reflect.Value
  fun func(a, b reflect.Value) bool
}

func (self *sortByValue) Len() int {
  return reflect.ValueOf(self.m).Len()
}

func (self *sortByValue) Less(i, j int) bool {
  v := reflect.ValueOf(self.m)
  return self.fun(v.MapIndex(self.key.Index(i)),
    v.MapIndex(self.key.Index(j)))
}

func (self *sortByValue) Swap(i, j int) {
  tmp := reflect.ValueOf(self.key.Index(i).Interface())
  self.key.Index(i).Set(self.key.Index(j))
  self.key.Index(j).Set(tmp)
}

func ByValue(m interface{}, fun func(a, b reflect.Value) bool) reflect.Value {
  sm := new(sortByValue)
  sm.m = m
  sm.fun = fun
  t := reflect.TypeOf(m)
  v := reflect.ValueOf(m)
  sm.key = reflect.MakeSlice(reflect.SliceOf(t.Key()), 0, v.Len())
  for _, key := range v.MapKeys() {
    sm.key = reflect.Append(sm.key, key)
  }
  sort.Sort(sm)
  return sm.key
}
