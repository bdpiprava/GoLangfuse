package types

import (
	"reflect"
)

// deepCopyAny performs a deep copy of any value
func deepCopyAny(src any) any {
	if src == nil {
		return nil
	}

	val := reflect.ValueOf(src)
	copied := deepCopyValue(val)
	if !copied.IsValid() {
		return nil
	}
	return copied.Interface()
}

// deepCopyValue performs deep copy using reflection
func deepCopyValue(src reflect.Value) reflect.Value {
	if !src.IsValid() {
		return reflect.Value{}
	}

	switch src.Kind() {
	case reflect.Ptr:
		if src.IsNil() {
			return src
		}
		dst := reflect.New(src.Type().Elem())
		dst.Elem().Set(deepCopyValue(src.Elem()))
		return dst

	case reflect.Slice:
		if src.IsNil() {
			return src
		}
		dst := reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
		for i := range src.Len() {
			dst.Index(i).Set(deepCopyValue(src.Index(i)))
		}
		return dst

	case reflect.Map:
		if src.IsNil() {
			return src
		}
		dst := reflect.MakeMap(src.Type())
		for _, key := range src.MapKeys() {
			// Deep copy both key and value
			copiedKey := deepCopyValue(key)
			copiedValue := deepCopyValue(src.MapIndex(key))
			dst.SetMapIndex(copiedKey, copiedValue)
		}
		return dst

	case reflect.Struct:
		dst := reflect.New(src.Type()).Elem()
		for i := range src.NumField() {
			srcField := src.Field(i)
			if srcField.CanInterface() && dst.Field(i).CanSet() {
				dst.Field(i).Set(deepCopyValue(srcField))
			}
		}
		return dst

	case reflect.Array:
		dst := reflect.New(src.Type()).Elem()
		for i := range src.Len() {
			dst.Index(i).Set(deepCopyValue(src.Index(i)))
		}
		return dst

	case reflect.Interface:
		if src.IsNil() {
			return src
		}
		// Get the concrete value and deep copy it
		concrete := src.Elem()
		if !concrete.IsValid() {
			return src
		}
		copied := deepCopyValue(concrete)
		if !copied.IsValid() {
			return src
		}
		// Create a new interface value with the copied concrete value
		dst := reflect.New(src.Type()).Elem()
		dst.Set(copied)
		return dst

	default:
		// For basic types (int, string, bool, etc.), just return the value
		return src
	}
}
