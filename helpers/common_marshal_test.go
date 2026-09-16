package helpers

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func TestExpandStringSlice(t *testing.T) {
	input := []interface{}{"a", "b", nil, "c"}
	expected := []string{"a", "b", "", "c"}
	actual := ExpandStringSlice(input)

	if actual == nil || !reflect.DeepEqual(*actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestExpandFloatSlice(t *testing.T) {
	input := []interface{}{1.1, nil, 2.2}
	expected := []float64{1.1, 2.2}
	actual := ExpandFloatSlice(input)

	if actual == nil || !reflect.DeepEqual(*actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestExpandFloatRangeSlice(t *testing.T) {
	input := []interface{}{
		[]interface{}{1.1, 2.2},
		nil,
		[]interface{}{3.3, 4.4},
	}
	expected := [][]float64{{1.1, 2.2}, {3.3, 4.4}}
	actual := ExpandFloatRangeSlice(input)

	if actual == nil || !reflect.DeepEqual(*actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestExpandPtrMapStringString(t *testing.T) {
	input := map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
	}
	expected := map[string]string{
		"key1": "val1",
		"key2": "val2",
	}
	actual := ExpandPtrMapStringString(input)

	if actual == nil || !reflect.DeepEqual(*actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestExpandMapStringPtrString(t *testing.T) {
	input := map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
	}
	expected := map[string]*string{
		"key1": pointer.To("val1"),
		"key2": pointer.To("val2"),
	}
	actual := ExpandMapStringPtrString(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestExpandInt32Slice(t *testing.T) {
	input := []interface{}{int(1), int(2)}
	expected := []int32{1, 2}
	actual := ExpandInt32Slice(input)

	if actual == nil || !reflect.DeepEqual(*actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestExpandInt64Slice(t *testing.T) {
	input := []interface{}{int(1), int(2)}
	expected := []int64{1, 2}
	actual := ExpandInt64Slice(input)

	if actual == nil || !reflect.DeepEqual(*actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestFlattenStringSlice(t *testing.T) {
	expected := []interface{}{"a", "b"}
	actual := FlattenStringSlice(pointer.To([]string{"a", "b"}))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}

	actualNil := FlattenStringSlice(nil)
	if !reflect.DeepEqual(actualNil, []interface{}{}) {
		t.Fatalf("expected empty slice for nil, got: %v", actualNil)
	}
}

func TestFlattenFloatSlice(t *testing.T) {
	expected := []interface{}{1.1, 2.2}
	actual := FlattenFloatSlice(pointer.To([]float64{1.1, 2.2}))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}

	actualNil := FlattenFloatSlice(nil)
	if !reflect.DeepEqual(actualNil, []interface{}{}) {
		t.Fatalf("expected empty slice for nil, got: %v", actualNil)
	}
}

func TestFlattenFloatRangeSlice(t *testing.T) {
	expected := [][]interface{}{{1.1, 2.2}, {3.3, 4.4}}
	actual := FlattenFloatRangeSlice(pointer.To([][]float64{{1.1, 2.2}, {3.3, 4.4}}))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}

	actualNil := FlattenFloatRangeSlice(nil)
	if !reflect.DeepEqual(actualNil, [][]interface{}{}) {
		t.Fatalf("expected empty slice for nil, got: %v", actualNil)
	}
}

func TestFlattenMapStringPtrString(t *testing.T) {
	input := map[string]*string{
		"key1": pointer.To("val1"),
		"key2": nil,
	}
	expected := map[string]interface{}{
		"key1": "val1",
		"key2": "",
	}
	actual := FlattenMapStringPtrString(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestFlattenPtrMapStringString(t *testing.T) {
	input := map[string]string{
		"key1": "val1",
		"key2": "val2",
	}
	expected := map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
	}
	actual := FlattenPtrMapStringString(&input)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}

	actualNil := FlattenPtrMapStringString(nil)
	if !reflect.DeepEqual(actualNil, map[string]interface{}{}) {
		t.Fatalf("expected empty map for nil, got: %v", actualNil)
	}
}

func TestFlattenInt32Slice(t *testing.T) {
	expected := []interface{}{int32(1), int32(2)}
	actual := FlattenInt32Slice(pointer.To([]int32{1, 2}))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}

	actualNil := FlattenInt32Slice(nil)
	if !reflect.DeepEqual(actualNil, []interface{}{}) {
		t.Fatalf("expected empty slice for nil, got: %v", actualNil)
	}
}

func TestFlattenInt64Slice(t *testing.T) {
	expected := []interface{}{int64(1), int64(2)}
	actual := FlattenInt64Slice(pointer.To([]int64{1, 2}))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}

	actualNil := FlattenInt64Slice(nil)
	if !reflect.DeepEqual(actualNil, []interface{}{}) {
		t.Fatalf("expected empty slice for nil, got: %v", actualNil)
	}
}

func TestExpandStringSliceWithDelimiter(t *testing.T) {
	input := []interface{}{"a", "b", nil, "c"}
	expected := "a,b,,c"
	actual := ExpandStringSliceWithDelimiter(input, ",")

	if actual == nil || *actual != expected {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestExpandIntSliceWithDelimiter(t *testing.T) {
	input := []interface{}{int(1), int(2), nil, int(3)}
	expected := "1,2,,3"
	actual := ExpandIntSliceWithDelimiter(input, ",")

	if actual == nil || *actual != expected {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestFlattenStringSliceWithDelimiter(t *testing.T) {
	expected := []interface{}{"a", "b", "", "c"}
	actual := FlattenStringSliceWithDelimiter(pointer.To("a,b,,c"), ",")

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}

	actualNil := FlattenStringSliceWithDelimiter(nil, ",")
	if !reflect.DeepEqual(actualNil, []interface{}{}) {
		t.Fatalf("expected empty slice for nil, got: %v", actualNil)
	}
}

// --- tests for the exported generics ---

func TestGenericExpandSlice(t *testing.T) {
	t.Run("booleans", func(t *testing.T) {
		input := []interface{}{true, false, nil, true}
		expected := []bool{true, false, false, true}
		actual := ExpandSlice(input, func(i bool) bool { return i }, true)
		if actual == nil || !reflect.DeepEqual(*actual, expected) {
			t.Fatalf("expected: %v, got: %v", expected, actual)
		}
	})

	t.Run("struct pointers", func(t *testing.T) {
		type MyStruct struct{ Name string }
		s1 := &MyStruct{Name: "s1"}
		input := []interface{}{s1, nil}
		expected := []*MyStruct{s1, nil}
		actual := ExpandSlice(input, func(i *MyStruct) *MyStruct { return i }, true)
		if actual == nil || !reflect.DeepEqual(*actual, expected) {
			t.Fatalf("expected: %v, got: %v", expected, actual)
		}
	})

	t.Run("skip nils", func(t *testing.T) {
		input := []interface{}{1, nil, 3}
		expected := []int{1, 3}
		actual := ExpandSlice(input, func(i int) int { return i }, false)
		if actual == nil || !reflect.DeepEqual(*actual, expected) {
			t.Fatalf("expected: %v, got: %v", expected, actual)
		}
	})
}

func TestGenericFlattenSlice(t *testing.T) {
	t.Run("booleans", func(t *testing.T) {
		expected := []interface{}{true, false}
		actual := FlattenSlice(pointer.To([]bool{true, false}))
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("expected: %v, got: %v", expected, actual)
		}
	})

	t.Run("structs", func(t *testing.T) {
		type MyStruct struct{ Name string }
		expected := []interface{}{MyStruct{Name: "s1"}, MyStruct{Name: "s2"}}
		actual := FlattenSlice(pointer.To([]MyStruct{{Name: "s1"}, {Name: "s2"}}))
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("expected: %v, got: %v", expected, actual)
		}
	})
}

func TestGenericExpandMap(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		input := map[string]interface{}{
			"a": 1,
			"b": 2,
		}
		expected := map[string]int{
			"a": 1,
			"b": 2,
		}
		actual := ExpandMap(input, func(i int) int { return i })
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("expected: %v, got: %v", expected, actual)
		}
	})

	t.Run("booleans", func(t *testing.T) {
		input := map[string]interface{}{
			"a": true,
			"b": false,
		}
		expected := map[string]bool{
			"a": true,
			"b": false,
		}
		actual := ExpandMap(input, func(i bool) bool { return i })
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("expected: %v, got: %v", expected, actual)
		}
	})
}

func TestGenericExpandSliceWithDelimiter(t *testing.T) {
	t.Run("floats", func(t *testing.T) {
		input := []interface{}{1.1, 2.2, nil, 3.3}
		expected := "1.1|2.2||3.3"
		actual := ExpandSliceWithDelimiter(input, func(i float64) string {
			return strconv.FormatFloat(i, 'f', 1, 64)
		}, "|")

		if actual == nil || *actual != expected {
			t.Fatalf("expected: %v, got: %v", expected, actual)
		}
	})
}
