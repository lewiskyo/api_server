package convert

import (
	"reflect"
	"strconv"
	"unsafe"
)

//InterfaceToString convert interface to string.
// @param v
// @return string
func InterfaceToString(v interface{}) string {
	return ToString(v)
}

//InterfaceToInt  convert interface to int.
// @param v
// @return int
func InterfaceToInt(v interface{}) int {
	return ToInt(v)
}

//InterfaceToInt64 convert interface to int64.
// @param v
// @return int64
func InterfaceToInt64(v interface{}) int64 {
	return ToInt64(v)
}

//InterfaceToFloat64 convert interface to float64.
// @param v
// @return float64
func InterfaceToFloat64(v interface{}) float64 {
	return ToFloat64(v)
}

//InterfaceToBool  convert interface to bool.
// @param v
// @return bool
func InterfaceToBool(v interface{}) bool {
	return ToBool(v)
}

// B2S converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// S2B converts string to a byte slice without memory allocation.
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func S2B(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

// Atoi string covert to int
func Atoi(s string, df int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return df
	}
	return i
}
