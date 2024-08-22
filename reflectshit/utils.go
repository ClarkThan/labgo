package reflectshit

import (
	"fmt"
	"reflect"
	"strings"
)

func structTo(rv reflect.Value) map[string]any {
	rt := rv.Type()
	ret := make(map[string]any, rt.NumField())

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.PkgPath != "" { // unexported field
			continue
		}
		v := rv.Field(i)
		key := f.Name
		jsonTag := rt.Field(i).Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		tagName, tagOpts := parseStructTag(jsonTag)
		if tagName != "" {
			key = tagName
		}

		if tagOpts.Has("omitempty") && v.IsZero() {
			continue
		}
		// fmt.Printf("[%s](%t): %s -> %v [%s]\n", v.Kind(), v.IsZero(), f.Name, v, jsonTag)
		// ret[key] = v.Interface()

		if tagOpts.Has("omitnested") {
			ret[key] = v.Interface()
		} else {
			ret[key] = nestedConvert(v)
		}
	}
	return ret
}

func struct2Map(st any) map[string]any {
	rv := reflect.ValueOf(st)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	if tk := rv.Kind(); tk != reflect.Struct {
		panic(fmt.Sprintf("only support struct, got %s", tk))
	}

	return structTo(rv)
}

func nestedConvert(rv reflect.Value) any {
	// fmt.Println("nested check --->", name)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		// []**string -> string
		rt := rv.Type().Elem()
		for rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
		}

		switch rt.Kind() {
		case reflect.Uint8:
			ret := make([]uint8, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i)).(uint8)
			}
			return ret
		case reflect.Int8:
			ret := make([]int8, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i)).(int8)
			}
			return ret
		case reflect.Uint16:
			ret := make([]uint16, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i)).(uint16)
			}
			return ret
		case reflect.Int16:
			ret := make([]int16, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i)).(int16)
			}
			return ret
		case reflect.Uint32:
			ret := make([]uint32, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i)).(uint32)
			}
			return ret
		case reflect.Int32:
			ret := make([]int32, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i)).(int32)
			}
			return ret
		case reflect.Uint64:
			ret := make([]uint64, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i)).(uint64)
			}
			return ret
		case reflect.Int64:
			ret := make([]int64, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i)).(int64)
			}
			return ret
		case reflect.Uint:
			ret := make([]uint, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i)).(uint)
			}
			return ret
		case reflect.Int:
			ret := make([]int, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i)).(int)
			}
			return ret
		case reflect.Float32:
			ret := make([]float32, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i)).(float32)
			}
			return ret
		case reflect.Float64:
			ret := make([]float64, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i)).(float64)
			}
			return ret
		case reflect.Bool:
			ret := make([]bool, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i)).(bool)
			}
			return ret
		case reflect.String:
			ret := make([]string, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i)).(string)
			}
			return ret
		default:
			ret := make([]any, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ret[i] = nestedConvert(rv.Index(i))
			}
			return ret
		}

		// if rv.IsZero() {
		// 	return reflect.Zero(rv.Type()).Interface()
		// 	// return []any{}
		// }
		// if rv.Type().Elem().Kind() != reflect.Struct &&
		// 	!(rv.Type().Elem().Kind() == reflect.Ptr &&
		// 		rv.Type().Elem().Elem().Kind() == reflect.Struct) {
		// 	fmt.Println("shit ---------")
		// 	return rv.Interface()
		// }

		// ret := make([]interface{}, rv.Len())
		// for i := 0; i < rv.Len(); i++ {
		// 	ret[i] = nestedConvert(rv.Index(i))
		// }
		// return ret

	// case reflect.Chan:
	case reflect.Map:
		rt := rv.Type()
		kt, vt := rt.Key(), rt.Elem()
		if kt.Kind() != reflect.String {
			panic("map key is not string")
		}

		for vt.Kind() == reflect.Pointer {
			vt = vt.Elem()
		}

		switch vt.Kind() {
		case reflect.Uint8:
			keys := rv.MapKeys()
			ret := make(map[string]uint8, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k)).(uint8)
			}
			return ret
		case reflect.Int8:
			keys := rv.MapKeys()
			ret := make(map[string]int8, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k)).(int8)
			}
			return ret
		case reflect.Uint16:
			keys := rv.MapKeys()
			ret := make(map[string]uint16, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k)).(uint16)
			}
			return ret
		case reflect.Int16:
			keys := rv.MapKeys()
			ret := make(map[string]int16, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k)).(int16)
			}
			return ret
		case reflect.Uint32:
			keys := rv.MapKeys()
			ret := make(map[string]uint32, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k)).(uint32)
			}
			return ret
		case reflect.Int32:
			keys := rv.MapKeys()
			ret := make(map[string]int32, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k)).(int32)
			}
			return ret
		case reflect.Uint64:
			keys := rv.MapKeys()
			ret := make(map[string]uint64, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k)).(uint64)
			}
			return ret
		case reflect.Int64:
			keys := rv.MapKeys()
			ret := make(map[string]int64, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k)).(int64)
			}
			return ret
		case reflect.Uint:
			keys := rv.MapKeys()
			ret := make(map[string]uint, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k)).(uint)
			}
			return ret
		case reflect.Int:
			keys := rv.MapKeys()
			ret := make(map[string]int, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k)).(int)
			}
			return ret
		case reflect.Float32:
			keys := rv.MapKeys()
			ret := make(map[string]float32, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k)).(float32)
			}
			return ret
		case reflect.Float64:
			keys := rv.MapKeys()
			ret := make(map[string]float64, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k)).(float64)
			}
			return ret
		case reflect.Bool:
			keys := rv.MapKeys()
			ret := make(map[string]bool, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k)).(bool)
			}
			return ret
		case reflect.String:
			keys := rv.MapKeys()
			ret := make(map[string]string, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k)).(string)
			}
			return ret
		default:
			// ret := make(map[string]any)
			// iter := rv.MapRange()
			// for iter.Next() {
			// 	ret[iter.Key().String()] = nestedConvert(iter.Value())
			// }
			keys := rv.MapKeys()
			ret := make(map[string]any, len(keys))
			for _, k := range keys {
				ret[k.String()] = nestedConvert(rv.MapIndex(k))
			}
			return ret
		}

		// if rv.IsNil() {
		// 	return reflect.Zero(rv.Type()).Interface()
		// 	// return map[string]any{}
		// }

		// // get the element type of the map
		// mapElem := rv.Type()
		// switch mapElem.Kind() {
		// case reflect.Pointer, reflect.Array, reflect.Map, reflect.Slice, reflect.Chan:
		// 	mapElem = mapElem.Elem()
		// 	for mapElem.Kind() == reflect.Pointer {
		// 		mapElem = mapElem.Elem()
		// 	}
		// 	// if mapElem.Kind() == reflect.Pointer {
		// 	// 	mapElem = mapElem.Elem()
		// 	// }
		// }

		// if mapElem.Kind() == reflect.Struct || (mapElem.Kind() == reflect.Slice && mapElem.Elem().Kind() == reflect.Struct) {
		// 	keys := rv.MapKeys()
		// 	ret := make(map[string]any, len(keys))
		// 	for _, k := range keys {
		// 		ret[k.String()] = nestedConvert(rv.MapIndex(k))
		// 	}

		// 	return ret
		// }

		// return rv.Interface()
	case reflect.Pointer:
		if rv.IsNil() {
			return rv.Interface()
			// return nil
		}
		return nestedConvert(rv.Elem())
	case reflect.Struct:
		return structTo(rv)
	default:
		return rv.Interface()
	}
}

type tagOpt []string

func (o tagOpt) Has(t string) bool {
	for _, opt := range []string(o) {
		if opt == t {
			return true
		}
	}

	return false
}

func parseStructTag(tag string) (string, tagOpt) {
	ret := strings.Split(tag, ",")
	if len(ret) >= 2 {
		return ret[0], ret[1:]
	}

	return ret[0], nil
}
