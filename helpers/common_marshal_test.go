package helpers

import (
	"reflect"
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
	input := []string{"a", "b"}
	expected := []interface{}{"a", "b"}
	actual := FlattenStringSlice(&input)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}

	actualNil := FlattenStringSlice(nil)
	if !reflect.DeepEqual(actualNil, []interface{}{}) {
		t.Fatalf("expected empty slice for nil, got: %v", actualNil)
	}
}

func TestFlattenFloatSlice(t *testing.T) {
	input := []float64{1.1, 2.2}
	expected := []interface{}{1.1, 2.2}
	actual := FlattenFloatSlice(&input)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}

	actualNil := FlattenFloatSlice(nil)
	if !reflect.DeepEqual(actualNil, []interface{}{}) {
		t.Fatalf("expected empty slice for nil, got: %v", actualNil)
	}
}

func TestFlattenFloatRangeSlice(t *testing.T) {
	input := [][]float64{{1.1, 2.2}, {3.3, 4.4}}
	expected := [][]interface{}{{1.1, 2.2}, {3.3, 4.4}}
	actual := FlattenFloatRangeSlice(&input)

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
	input := []int32{1, 2}
	expected := []interface{}{int32(1), int32(2)}
	actual := FlattenInt32Slice(&input)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}

	actualNil := FlattenInt32Slice(nil)
	if !reflect.DeepEqual(actualNil, []interface{}{}) {
		t.Fatalf("expected empty slice for nil, got: %v", actualNil)
	}
}

func TestFlattenInt64Slice(t *testing.T) {
	input := []int64{1, 2}
	expected := []interface{}{int64(1), int64(2)}
	actual := FlattenInt64Slice(&input)

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
	input := "a,b,,c"
	expected := []interface{}{"a", "b", "", "c"}
	actual := FlattenStringSliceWithDelimiter(&input, ",")

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}

	actualNil := FlattenStringSliceWithDelimiter(nil, ",")
	if !reflect.DeepEqual(actualNil, []interface{}{}) {
		t.Fatalf("expected empty slice for nil, got: %v", actualNil)
	}
}
