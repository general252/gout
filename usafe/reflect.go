package usafe

import (
	"errors"
	"fmt"
	"reflect"
)

func ReflectNew(object interface{}) (v interface{}, err error) {
	if object == nil {
		return nil, errors.New("nil param")
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	var (
		typeObject  reflect.Type
		valueObject reflect.Value
		target      interface{}
	)

	typeObject = reflect.TypeOf(object)
	if typeObject.Kind() == reflect.Ptr {
		typeObject = typeObject.Elem()
	}

	// slice: reflect.MakeSlice(reflect.SliceOf(typeObject), 10, 10).Interface()
	valueObject = reflect.New(typeObject)

	target = valueObject.Elem().Interface()

	return target, nil
}
