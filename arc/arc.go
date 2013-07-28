package arc

import (
  "sync/atomic"
  "reflect"
  "runtime"
)

type ARC struct {
  count int32
  cb Callback
}

type Callback func(c int32)

func New(cb Callback) *ARC {
  return &ARC{
    cb: cb,
  }
}

func (self *ARC) Clone(p interface{}, nilfunc interface{}) {
  atomic.AddInt32(&self.count, 1)
  funcType := reflect.TypeOf(nilfunc)
  funcImpl := func(args []reflect.Value) []reflect.Value {
    c := atomic.AddInt32(&self.count, -1)
    self.cb(c)
    return nil
  }
  runtime.SetFinalizer(p, reflect.MakeFunc(funcType, funcImpl).Interface())
}
