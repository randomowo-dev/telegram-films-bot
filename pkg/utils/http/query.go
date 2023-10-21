package http

import (
	"fmt"
	netUrl "net/url"
	"reflect"
	"strconv"
	"strings"
)

func ParseToQuery(v reflect.Value, key string, query *netUrl.Values) error {
	switch v.Kind() {
	case reflect.Struct:
		for _, sf := range reflect.VisibleFields(v.Type()) {
			if err := ParseToQuery(v.FieldByIndex(sf.Index), sf.Tag.Get("params"), query); err != nil {
				return err
			}
		}
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			var value string
			mv := iter.Value()
			switch mv.Kind() {
			case reflect.String:
				value = mv.String()
			case reflect.Int:
				value = strconv.Itoa(mv.Interface().(int))
			default:
				return fmt.Errorf("unsupported map param type")
			}
			query.Add(iter.Key().String(), value)
		}
	case reflect.Array, reflect.Slice:
		values := make([]string, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			av := v.Index(i)
			switch av.Kind() {
			case reflect.String:
				values = append(values, av.String())
			case reflect.Int:
				values = append(values, strconv.Itoa(av.Interface().(int)))
			default:
				return fmt.Errorf("unsupported map param type")
			}
		}
		if v.Len() > 0 {
			query.Add(key, strings.Join(values, ","))
		}
	case reflect.Int:
		query.Add(key, strconv.Itoa(v.Interface().(int)))
	case reflect.String:
		query.Add(key, v.String())
	}
	return nil
}
