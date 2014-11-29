package context

import (
	"image/color"
	"log"
	"reflect"
)

// #include "sass_context.h"
import "C"

type SassValue interface{}

// make is needed to create types for use by test
func makevalue(v interface{}) *C.union_Sass_Value {
	f := reflect.ValueOf(v)
	switch f.Kind() {
	default:
		log.Printf("% #v\n", v)
		log.Fatalf("Fail: %v", f.Kind())
		return C.sass_make_null()
	case reflect.String:
		return C.sass_make_string(C.CString(v.(string)))
	case reflect.Slice:
		// Initialize the list
		l := C.sass_make_list(C.size_t(f.Len()), C.SASS_COMMA)

		for i := 0; i < f.Len(); i++ {
			t := makevalue(f.Index(i).Interface())
			C.sass_list_set_value(l, C.size_t(i), t)
		}

		return l
	}
}

func unmarshal(arg *C.union_Sass_Value, v interface{}) {
	f := reflect.ValueOf(v).Elem()
	switch {
	case bool(C.sass_value_is_null(arg)):
		//return nil
	case bool(C.sass_value_is_number(arg)):
		// return int(C.sass_number_get_value(arg))
	case bool(C.sass_value_is_string(arg)):
		c := C.sass_string_get_value(arg)
		gc := C.GoString(c)
		if !f.CanSet() {
			return
		}

		switch f.Kind() {
		default:
			log.Fatal("unknown type")
		case reflect.String:
			f.SetString(gc)
		case reflect.Interface:
			f.Set(reflect.ValueOf(gc))
		}
	case bool(C.sass_value_is_boolean(arg)):
		// return bool(C.sass_boolean_get_value(arg))
	case bool(C.sass_value_is_color(arg)):
		col := color.RGBA{
			R: uint8(C.sass_color_get_r(arg)),
			G: uint8(C.sass_color_get_g(arg)),
			B: uint8(C.sass_color_get_b(arg)),
			A: uint8(C.sass_color_get_a(arg)),
		}
		_ = col
		// return col
	case bool(C.sass_value_is_list(arg)):
		l := make([]SassValue, C.sass_list_get_length(arg))
		for i := range l {
			unmarshal(C.sass_list_get_value(arg, C.size_t(i)), &l[i])
		}
		fl := reflect.ValueOf(l)
		f.Set(fl)
	case bool(C.sass_value_is_map(arg)):
		len := int(C.sass_map_get_length(arg))
		m := make(map[SassValue]SassValue, len)
		for i := 0; i < len; i++ {
			//m[Decode(C.sass_map_get_key(arg, C.size_t(i)))] =
			//	Decode(C.sass_map_get_value(arg, C.size_t(i)))
		}
		_ = m
		// return m
	case bool(C.sass_value_is_error(arg)):
		// return C.GoString(C.sass_error_get_message(arg))
	}
	return
}

// Decode converts Sass Value to Go compatible data types.
func Unmarshal(arg *C.union_Sass_Value, v interface{}) {
	unmarshal(arg, v)
}
