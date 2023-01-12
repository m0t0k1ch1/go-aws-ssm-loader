package ssmloader

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

func Load(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("v must be a pointer of a struct")
	}

	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Struct {
		return errors.New("v must be a pointer of a struct")
	}

	rt := rv.Type()

	keys := make([]string, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		keys[i] = rt.Field(i).Tag.Get("ssm")
	}

	// WIP
	fmt.Println(keys)

	return nil
}
