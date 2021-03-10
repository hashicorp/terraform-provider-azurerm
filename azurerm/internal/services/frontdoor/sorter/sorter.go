package sorter

import (
	"fmt"
	"reflect"
	"sort"
)

type sorter struct {
	data  interface{}
	order Comparators
}

func FDSorter() *sorter {
	return &sorter{}
}

func (s *sorter) ByKeys(order []string) *sorter {
	comparators := make(Comparators, 0)
	for _, key := range order {
		comparators = append(comparators, Comparator{key, Ascending})
	}
	return s.ByKeyComparators(comparators)
}

type Comparator struct {
	Name string
	Comp func(interface{}, interface{}) CompareResult
}

type Comparators []Comparator

func (s *sorter) ByKeyComparators(comparators Comparators) *sorter {
	s.order = comparators
	return s
}

func (s *sorter) DoSort(data interface{}) {
	s.data = data
	sort.Sort(s)
}

func (s *sorter) Len() int {
	return reflect.ValueOf(s.data).Len()
}

func (s *sorter) Swap(i, j int) {
	if i > j {
		i, j = j, i
	}
	arr := reflect.ValueOf(s.data)

	tmp := arr.Index(i).Interface()
	arr.Index(i).Set(arr.Index(j))
	arr.Index(j).Set(reflect.ValueOf(tmp))
}

func (s *sorter) Less(i, j int) bool {
	arr := reflect.ValueOf(s.data)
	a := reflect.ValueOf(arr.Index(i).Interface())
	b := reflect.ValueOf(arr.Index(j).Interface())
	if a.Kind() != reflect.Map {
		iface := a.Interface()
		panic(fmt.Sprintf("[A] Kind: %s, Expected a map, but got a %T for %v", a.Kind(), iface, iface))
	}
	if b.Kind() != reflect.Map {
		iface := b.Interface()
		panic(fmt.Sprintf("[B] Kind: %s, Expected a map, but got a %T for %v", b.Kind(), iface, iface))
	}

	for i := 0; i < len(s.order); i += 1 {
		keyComp := s.order[i]
		af := a.MapIndex(reflect.ValueOf(keyComp.Name)).Interface()
		bf := b.MapIndex(reflect.ValueOf(keyComp.Name)).Interface()

		switch keyComp.Comp(af, bf) {
		case LESSER:
			return true
		case GREATER:
			return false
		}
	}
	return true
}

type CompareResult int8

const (
	LESSER CompareResult = -1 + iota
	EQUAL
	GREATER
)

func Ascending(a, b interface{}) CompareResult {
	switch Descending(a, b) {
	case LESSER:
		return GREATER
	case GREATER:
		return LESSER
	default:
		return EQUAL
	}
}

func Descending(a, b interface{}) CompareResult {
	if a == b {
		return EQUAL
	}
	switch a.(type) {
	case string:
		return lg(a.(string) > b.(string))
	case int:
		return lg(a.(int) > reflect.ValueOf(b).Convert(reflect.TypeOf(a)).Interface().(int))
	case int8:
		return lg(a.(int8) > reflect.ValueOf(b).Convert(reflect.TypeOf(a)).Interface().(int8))
	case int16:
		return lg(a.(int16) > reflect.ValueOf(b).Convert(reflect.TypeOf(a)).Interface().(int16))
	case int32:
		return lg(a.(int32) > reflect.ValueOf(b).Convert(reflect.TypeOf(a)).Interface().(int32))
	case int64:
		return lg(a.(int64) > reflect.ValueOf(b).Convert(reflect.TypeOf(a)).Interface().(int64))
	case uint:
		return lg(a.(uint) > reflect.ValueOf(b).Convert(reflect.TypeOf(a)).Interface().(uint))
	case uint8:
		return lg(a.(uint8) > reflect.ValueOf(b).Convert(reflect.TypeOf(a)).Interface().(uint8))
	case uint16:
		return lg(a.(uint16) > reflect.ValueOf(b).Convert(reflect.TypeOf(a)).Interface().(uint16))
	case uint32:
		return lg(a.(uint32) > reflect.ValueOf(b).Convert(reflect.TypeOf(a)).Interface().(uint32))
	case uint64:
		return lg(a.(uint64) > reflect.ValueOf(b).Convert(reflect.TypeOf(a)).Interface().(uint64))
	case float32:
		return lg(a.(float32) > reflect.ValueOf(b).Convert(reflect.TypeOf(a)).Interface().(float32))
	case float64:
		return lg(a.(float64) > reflect.ValueOf(b).Convert(reflect.TypeOf(a)).Interface().(float64))
	default:
		panic(fmt.Sprintf("dont know how to compare: %T", a))
	}
}

func lg(b bool) CompareResult {
	if b {
		return LESSER
	}
	return GREATER
}
