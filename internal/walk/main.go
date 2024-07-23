package walk

import "reflect"

type Iter func(reflect.Value, reflect.StructField) bool

func Into(t any, action Iter) {
	value := reflect.ValueOf(t)
	if v := value.Kind(); v == reflect.Pointer || v == reflect.Interface {
		value = value.Elem()
	}

	rtype := reflect.TypeOf(t)
	if rtype.Kind() == reflect.Pointer {
		rtype = rtype.Elem()
	}
	Struct(value, rtype, action)
}

func Struct(
	value reflect.Value,
	str reflect.Type,
	action Iter,
) {
	for i := range str.NumField() {
		v := value.Field(i)
		t := str.Field(i)

		if !t.IsExported() {
			continue
		}

		if t.Type.Kind() == reflect.Struct {
			Struct(v, t.Type, action)
			continue
		}
		if !action(v, t) {
			break
		}
	}
}
