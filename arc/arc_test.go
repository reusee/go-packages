package arc

import (
  "testing"
  "fmt"
  "runtime"
)

type foo struct {}

func TestARC(t *testing.T) {
  arc := New(func(c int32) {
    fmt.Printf("%d\n", c)
  })

  for i := 0; i < 16; i++ {
    s := new(foo)
    arc.Clone(&s, func(**foo){})
    if i % 4 == 0 {
      runtime.GC()
    }
  }
}
