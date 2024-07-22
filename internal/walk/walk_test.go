package walk_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Tom5521/gtk4tools/internal/walk"
)

type A struct {
	F1 string
	F2 bool
	F3 int
	B  struct {
		F1 string
		C  struct {
			F1 string
		}
	}
}

func TestUnmarshall(_ *testing.T) {
	str := A{
		F1: "Meow in A",
		F2: true,
		F3: 20,
	}
	str.B.F1 = "Meow in B"
	str.B.C.F1 = "Meow in C"

	walk.Into(
		&str,
		func(v reflect.Value, _ reflect.StructField) bool {
			if v.String() == "Meow in C" {
				v.Set(reflect.ValueOf(v.String() + " | Edited"))
			}
			fmt.Println(v)
			return true
		},
	)
}
