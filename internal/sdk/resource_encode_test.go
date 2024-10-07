// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

type encodeTestData struct {
	Input       interface{}
	Expected    map[string]interface{}
	ExpectError bool
}

func TestResourceEncode_TopLevel(t *testing.T) {
	type SimpleType struct {
		String        string            `tfschema:"string"`
		Number        int64             `tfschema:"number"`
		Price         float64           `tfschema:"price"`
		Enabled       bool              `tfschema:"enabled"`
		ListOfFloats  []float64         `tfschema:"list_of_floats"`
		ListOfNumbers []int64           `tfschema:"list_of_numbers"`
		ListOfStrings []string          `tfschema:"list_of_strings"`
		MapOfBools    map[string]bool   `tfschema:"map_of_bools"`
		MapOfNumbers  map[string]int64  `tfschema:"map_of_numbers"`
		MapOfStrings  map[string]string `tfschema:"map_of_strings"`
	}

	encodeTestData{
		Input: &SimpleType{
			String:  "world",
			Number:  42,
			Price:   129.99,
			Enabled: true,
			ListOfFloats: []float64{
				1.0,
				2.0,
				3.0,
				1.234567890,
			},
			ListOfNumbers: []int64{1, 2, 3},
			ListOfStrings: []string{
				"have",
				"you",
				"heard",
			},
			MapOfBools: map[string]bool{
				"awesome_feature": true,
			},
			MapOfNumbers: map[string]int64{
				"hello": 1,
				"there": 3,
			},
			MapOfStrings: map[string]string{
				"hello":   "there",
				"salut":   "tout les monde",
				"guten":   "tag",
				"morning": "alvaro",
			},
		},
		Expected: map[string]interface{}{
			"number":  int64(42),
			"price":   129.99,
			"string":  "world",
			"enabled": true,
			"list_of_floats": []float64{
				1.0,
				2.0,
				3.0,
				1.234567890,
			},
			"list_of_numbers": []int64{1, 2, 3},
			"list_of_strings": []string{
				"have",
				"you",
				"heard",
			},
			"map_of_bools": map[string]interface{}{
				"awesome_feature": true,
			},
			"map_of_numbers": map[string]interface{}{
				"hello": int64(1),
				"there": int64(3),
			},
			"map_of_strings": map[string]interface{}{
				"hello":   "there",
				"salut":   "tout les monde",
				"guten":   "tag",
				"morning": "alvaro",
			},
		},
	}.test(t)
}

func TestResourceEncode_TopLevelAllTypesAndCombinations(t *testing.T) {
	encodeTestData{
		Input: &OneOfEverything{
			RequiredStr:           "foo",
			OptionalStr:           pointer.To("bar"),
			RequiredInt64:         101,
			OptionalInt64:         pointer.To(int64(20)),
			RequiredFloat:         3.14159,
			OptionalFloat:         pointer.To(1.41442),
			RequiredBoolean:       true,
			OptionalBoolean:       pointer.To(true),
			RequiredListOfFloat:   []float64{3.142, 1.414, 2.718},
			OptionalListOfFloat:   pointer.To([]float64{2.718, 1.414, 3.142}),
			RequiredListOfInt64:   []int64{10, 20, 30, 40, 50},
			OptionalListOfInt64:   pointer.To([]int64{100, 90, 80, 70, 60}),
			RequiredListOfStrings: []string{"foo", "bar"},
			OptionalListOfStrings: pointer.To([]string{"bar", "foo"}),
			RequiredMapOfBooleans: map[string]bool{"itsTrue": true, "itsFalse": false},
			OptionalMapOfBooleans: pointer.To(map[string]bool{"itsNotNotFalse": false, "itsMoreTrue": true}),
			RequiredMapOfFloat:    map[string]float64{"avogadro": 6.022},
			OptionalMapOfFloat:    pointer.To(map[string]float64{"pythagoras": 1.41421, "pi": 3.14159}),
			RequiredMapOfInt64:    map[string]int64{"alpha": 200, "beta": 300},
			OptionalMapOfInt64:    pointer.To(map[string]int64{"gamma": 400, "delta": 500}),
			RequiredMapOfStrings:  map[string]string{"epsilon": "zeta", "eta": "theta"},
			OptionalMapOfStrings:  pointer.To(map[string]string{"iota": "kappa", "lambda": "mu", "nu": "xi"}),
		},
		Expected: map[string]interface{}{
			"required_str":             "foo",
			"optional_str":             "bar",
			"required_int64":           int64(101),
			"optional_int64":           int64(20),
			"required_float":           3.14159,
			"optional_float":           1.41442,
			"required_boolean":         true,
			"optional_boolean":         true,
			"required_list_of_float":   []float64{3.142, 1.414, 2.718},
			"optional_list_of_float":   []float64{2.718, 1.414, 3.142},
			"required_list_of_int64":   []int64{10, 20, 30, 40, 50},
			"optional_list_of_int64":   []int64{100, 90, 80, 70, 60},
			"required_list_of_strings": []string{"foo", "bar"},
			"optional_list_of_strings": []string{"bar", "foo"},
			"required_map_of_booleans": map[string]interface{}{"itsTrue": true, "itsFalse": false},
			"optional_map_of_booleans": map[string]interface{}{"itsNotNotFalse": false, "itsMoreTrue": true},
			"required_map_of_float":    map[string]interface{}{"avogadro": 6.022},
			"optional_map_of_float":    map[string]interface{}{"pythagoras": 1.41421, "pi": 3.14159},
			"required_map_of_int64":    map[string]interface{}{"alpha": int64(200), "beta": int64(300)},
			"optional_map_of_int64":    map[string]interface{}{"gamma": int64(400), "delta": int64(500)},
			"required_map_of_strings":  map[string]interface{}{"epsilon": "zeta", "eta": "theta"},
			"optional_map_of_strings":  map[string]interface{}{"iota": "kappa", "lambda": "mu", "nu": "xi"},
		},
	}.test(t)
}

func TestResourceEncode_TopLevelEmptyPointers(t *testing.T) {
	type SimpleType struct {
		String           string             `tfschema:"string"`
		StringPtr        *string            `tfschema:"string_ptr"`
		Number           int64              `tfschema:"number"`
		NumberPtr        *int64             `tfschema:"number_ptr"`
		Price            float64            `tfschema:"price"`
		PricePtr         *float64           `tfschema:"price_ptr"`
		Enabled          bool               `tfschema:"enabled"`
		EnabledPtr       *bool              `tfschema:"enabled_ptr"`
		ListOfFloats     []float64          `tfschema:"list_of_floats"`
		ListOfNumbers    []int64            `tfschema:"list_of_numbers"`
		ListOfStrings    []string           `tfschema:"list_of_strings"`
		ListOfStringsPtr *[]string          `tfschema:"list_of_strings_ptr"`
		MapOfBools       map[string]bool    `tfschema:"map_of_bools"`
		MapOfNumbers     map[string]int64   `tfschema:"map_of_numbers"`
		MapOfStrings     map[string]string  `tfschema:"map_of_strings"`
		MapOfStringsPtr  *map[string]string `tfschema:"map_of_strings_ptr"`
	}

	encodeTestData{
		Input: &SimpleType{
			String:  "world",
			Number:  42,
			Price:   129.99,
			Enabled: true,
			ListOfFloats: []float64{
				1.0,
				2.0,
				3.0,
				1.234567890,
			},
			ListOfNumbers: []int64{1, 2, 3},
			ListOfStrings: []string{
				"have",
				"you",
				"heard",
			},
			MapOfBools: map[string]bool{
				"awesome_feature": true,
			},
			MapOfNumbers: map[string]int64{
				"hello": 1,
				"there": 3,
			},
			MapOfStrings: map[string]string{
				"hello":   "there",
				"salut":   "tous les monde",
				"guten":   "tag",
				"morning": "alvaro",
			},
		},
		Expected: map[string]interface{}{
			"number":      int64(42),
			"number_ptr":  nil,
			"price":       float64(129.99),
			"price_ptr":   nil,
			"string":      "world",
			"string_ptr":  nil,
			"enabled":     true,
			"enabled_ptr": nil,
			"list_of_floats": []float64{
				1.0,
				2.0,
				3.0,
				1.234567890,
			},
			"list_of_numbers": []int64{1, 2, 3},
			"list_of_strings": []string{
				"have",
				"you",
				"heard",
			},
			"list_of_strings_ptr": nil,
			"map_of_bools": map[string]interface{}{
				"awesome_feature": true,
			},
			"map_of_numbers": map[string]interface{}{
				"hello": int64(1),
				"there": int64(3),
			},
			"map_of_strings": map[string]interface{}{
				"hello":   "there",
				"salut":   "tous les monde",
				"guten":   "tag",
				"morning": "alvaro",
			},
			"map_of_strings_ptr": nil,
		},
	}.test(t)
}

func TestResourceEncode_TopLevelNonEmptyPointers(t *testing.T) {
	type SimpleType struct {
		String           string             `tfschema:"string"`
		StringPtr        *string            `tfschema:"string_ptr"`
		Number           int64              `tfschema:"number"`
		NumberPtr        *int64             `tfschema:"number_ptr"`
		Price            float64            `tfschema:"price"`
		PricePtr         *float64           `tfschema:"price_ptr"`
		Enabled          bool               `tfschema:"enabled"`
		EnabledPtr       *bool              `tfschema:"enabled_ptr"`
		ListOfFloats     []float64          `tfschema:"list_of_floats"`
		ListOfNumbers    []int64            `tfschema:"list_of_numbers"`
		ListOfStrings    []string           `tfschema:"list_of_strings"`
		ListOfStringsPtr *[]string          `tfschema:"list_of_strings_ptr"`
		MapOfBools       map[string]bool    `tfschema:"map_of_bools"`
		MapOfNumbers     map[string]int64   `tfschema:"map_of_numbers"`
		MapOfStrings     map[string]string  `tfschema:"map_of_strings"`
		MapOfStringsPtr  *map[string]string `tfschema:"map_of_strings_ptr"`
	}

	encodeTestData{
		Input: &SimpleType{
			String:     "world",
			StringPtr:  pointer.To("foo"),
			Number:     42,
			NumberPtr:  pointer.To(int64(22)),
			Price:      129.99,
			PricePtr:   pointer.To(3.50),
			Enabled:    true,
			EnabledPtr: pointer.To(true),
			ListOfFloats: []float64{
				1.0,
				2.0,
				3.0,
				1.234567890,
			},
			ListOfNumbers: []int64{1, 2, 3},
			ListOfStrings: []string{
				"have",
				"you",
				"heard",
			},
			ListOfStringsPtr: pointer.To([]string{
				"about",
				"the",
				"bird",
			}),
			MapOfBools: map[string]bool{
				"awesome_feature": true,
			},
			MapOfNumbers: map[string]int64{
				"hello": 1,
				"there": 3,
			},
			MapOfStrings: map[string]string{
				"hello":   "there",
				"salut":   "tous les monde",
				"guten":   "tag",
				"morning": "alvaro",
			},
			MapOfStringsPtr: pointer.To(map[string]string{
				"foo": "bar",
			}),
		},
		Expected: map[string]interface{}{
			"number":      int64(42),
			"number_ptr":  int64(22),
			"price":       float64(129.99),
			"price_ptr":   float64(3.50),
			"string":      "world",
			"string_ptr":  "foo",
			"enabled":     true,
			"enabled_ptr": true,
			"list_of_floats": []float64{
				1.0,
				2.0,
				3.0,
				1.234567890,
			},
			"list_of_numbers": []int64{1, 2, 3},
			"list_of_strings": []string{
				"have",
				"you",
				"heard",
			},
			"list_of_strings_ptr": []string{
				"about",
				"the",
				"bird",
			},
			"map_of_bools": map[string]interface{}{
				"awesome_feature": true,
			},
			"map_of_numbers": map[string]interface{}{
				"hello": int64(1),
				"there": int64(3),
			},
			"map_of_strings": map[string]interface{}{
				"hello":   "there",
				"salut":   "tous les monde",
				"guten":   "tag",
				"morning": "alvaro",
			},
			"map_of_strings_ptr": map[string]interface{}{
				"foo": "bar",
			},
		},
	}.test(t)
}

func TestResourceEncode_TopLevelOmitted(t *testing.T) {
	type SimpleType struct {
		String        string            `tfschema:"string"`
		Number        int64             `tfschema:"number"`
		Price         float64           `tfschema:"price"`
		Enabled       bool              `tfschema:"enabled"`
		ListOfFloats  []float64         `tfschema:"list_of_floats"`
		ListOfNumbers []int64           `tfschema:"list_of_numbers"`
		ListOfStrings []string          `tfschema:"list_of_strings"`
		MapOfBools    map[string]bool   `tfschema:"map_of_bools"`
		MapOfNumbers  map[string]int64  `tfschema:"map_of_numbers"`
		MapOfStrings  map[string]string `tfschema:"map_of_strings"`
	}
	encodeTestData{
		Input: &SimpleType{},
		Expected: map[string]interface{}{
			"number":          int64(0),
			"price":           float64(0),
			"string":          "",
			"enabled":         false,
			"list_of_floats":  []float64{},
			"list_of_numbers": []int64{},
			"list_of_strings": []string{},
			"map_of_bools":    map[string]interface{}{},
			"map_of_numbers":  map[string]interface{}{},
			"map_of_strings":  map[string]interface{}{},
		},
	}.test(t)
}

func TestResourceEncode_TopLevelComputed(t *testing.T) {
	type SimpleType struct {
		ComputedString        string             `tfschema:"computed_string" computed:"true"`
		ComputedNumber        int64              `tfschema:"computed_number" computed:"true"`
		ComputedBool          bool               `tfschema:"computed_bool" computed:"true"`
		ComputedListOfNumbers []int64            `tfschema:"computed_list_of_numbers" computed:"true"`
		ComputedListOfStrings []string           `tfschema:"computed_list_of_strings" computed:"true"`
		ComputedMapOfBools    map[string]bool    `tfschema:"computed_map_of_bools" computed:"true"`
		ComputedMapOfFloats   map[string]float64 `tfschema:"computed_map_of_floats" computed:"true"`
		ComputedMapOfInts     map[string]int64   `tfschema:"computed_map_of_ints" computed:"true"`
		ComputedMapOfStrings  map[string]string  `tfschema:"computed_map_of_strings" computed:"true"`
	}
	encodeTestData{
		Input: &SimpleType{
			ComputedString:        "je suis computed",
			ComputedNumber:        732,
			ComputedBool:          true,
			ComputedListOfNumbers: []int64{1, 2, 3},
			ComputedListOfStrings: []string{
				"have",
				"you",
				"heard",
			},
			ComputedMapOfBools: map[string]bool{
				"hello": true,
				"world": false,
			},
			ComputedMapOfFloats: map[string]float64{
				"hello": 1.8965345678,
				"world": 2.0,
			},
			ComputedMapOfInts: map[string]int64{
				"first":  1,
				"second": 2,
				"third":  3,
			},
			ComputedMapOfStrings: map[string]string{
				"hello": "world",
				"bingo": "bango",
			},
		},
		Expected: map[string]interface{}{
			"computed_string":          "je suis computed",
			"computed_number":          int64(732),
			"computed_bool":            true,
			"computed_list_of_numbers": []int64{1, 2, 3},
			"computed_list_of_strings": []string{
				"have",
				"you",
				"heard",
			},
			"computed_map_of_bools": map[string]interface{}{
				"hello": true,
				"world": false,
			},
			"computed_map_of_floats": map[string]interface{}{
				"hello": 1.8965345678,
				"world": 2.0,
			},
			"computed_map_of_ints": map[string]interface{}{
				"first":  int64(1),
				"second": int64(2),
				"third":  int64(3),
			},
			"computed_map_of_strings": map[string]interface{}{
				"hello": "world",
				"bingo": "bango",
			},
		},
	}.test(t)
}

func TestResourceEncode_NestedOneLevelDeepEmpty(t *testing.T) {
	type Inner struct {
		Value string `tfschema:"value"`
	}
	type Type struct {
		NestedObject []Inner `tfschema:"inner"`
	}
	encodeTestData{
		Input: &Type{
			NestedObject: []Inner{},
		},
		Expected: map[string]interface{}{
			"inner": []interface{}{},
		},
	}.test(t)
}

func TestResourceEncode_NestedOneLevelDeepSingle(t *testing.T) {
	type Inner struct {
		String        string            `tfschema:"string"`
		Number        int64             `tfschema:"number"`
		Price         float64           `tfschema:"price"`
		Enabled       bool              `tfschema:"enabled"`
		ListOfFloats  []float64         `tfschema:"list_of_floats"`
		ListOfNumbers []int64           `tfschema:"list_of_numbers"`
		ListOfStrings []string          `tfschema:"list_of_strings"`
		MapOfBools    map[string]bool   `tfschema:"map_of_bools"`
		MapOfNumbers  map[string]int64  `tfschema:"map_of_numbers"`
		MapOfStrings  map[string]string `tfschema:"map_of_strings"`
	}
	type Type struct {
		NestedObject []Inner `tfschema:"inner"`
	}
	encodeTestData{
		Input: &Type{
			NestedObject: []Inner{
				{
					String:  "world",
					Number:  42,
					Price:   129.99,
					Enabled: true,
					ListOfFloats: []float64{
						1.0,
						2.0,
						3.0,
						1.234567890,
					},
					ListOfNumbers: []int64{1, 2, 3},
					ListOfStrings: []string{
						"have",
						"you",
						"heard",
					},
					MapOfBools: map[string]bool{
						"awesome_feature": true,
					},
					MapOfNumbers: map[string]int64{
						"hello": 1,
						"there": 3,
					},
					MapOfStrings: map[string]string{
						"hello":   "there",
						"salut":   "tous les monde",
						"guten":   "tag",
						"morning": "alvaro",
					},
				},
			},
		},
		Expected: map[string]interface{}{
			"inner": []interface{}{
				map[string]interface{}{
					"number":  int64(42),
					"price":   float64(129.99),
					"string":  "world",
					"enabled": true,
					"list_of_floats": []float64{
						1.0,
						2.0,
						3.0,
						1.234567890,
					},
					"list_of_numbers": []int64{1, 2, 3},
					"list_of_strings": []string{
						"have",
						"you",
						"heard",
					},
					"map_of_bools": map[string]interface{}{
						"awesome_feature": true,
					},
					"map_of_numbers": map[string]interface{}{
						"hello": int64(1),
						"there": int64(3),
					},
					"map_of_strings": map[string]interface{}{
						"hello":   "there",
						"salut":   "tous les monde",
						"guten":   "tag",
						"morning": "alvaro",
					},
				},
			},
		},
	}.test(t)
}

func TestResourceEncode_NestedOneLevelDeepSingleOmittedValues(t *testing.T) {
	type Inner struct {
		String               string             `tfschema:"string"`
		Number               int64              `tfschema:"number"`
		Price                float64            `tfschema:"price"`
		Enabled              bool               `tfschema:"enabled"`
		ListOfFloats         []float64          `tfschema:"list_of_floats"`
		ListOfNumbers        []int64            `tfschema:"list_of_numbers"`
		ListOfStrings        []string           `tfschema:"list_of_strings"`
		MapOfBools           map[string]bool    `tfschema:"map_of_bools"`
		MapOfNumbers         map[string]int64   `tfschema:"map_of_numbers"`
		MapOfStrings         map[string]string  `tfschema:"map_of_strings"`
		ComputedMapOfBools   map[string]bool    `tfschema:"computed_map_of_bools" computed:"true"`
		ComputedMapOfFloats  map[string]float64 `tfschema:"computed_map_of_floats" computed:"true"`
		ComputedMapOfInts    map[string]int64   `tfschema:"computed_map_of_ints" computed:"true"`
		ComputedMapOfStrings map[string]string  `tfschema:"computed_map_of_strings" computed:"true"`
	}
	type Type struct {
		NestedObject []Inner `tfschema:"inner"`
	}
	encodeTestData{
		Input: &Type{
			NestedObject: []Inner{
				{},
			},
		},
		Expected: map[string]interface{}{
			"inner": []interface{}{
				map[string]interface{}{
					"number":                  int64(0),
					"price":                   float64(0),
					"string":                  "",
					"enabled":                 false,
					"list_of_floats":          []float64{},
					"list_of_numbers":         []int64{},
					"list_of_strings":         []string{},
					"map_of_bools":            map[string]interface{}{},
					"map_of_numbers":          map[string]interface{}{},
					"map_of_strings":          map[string]interface{}{},
					"computed_map_of_bools":   map[string]interface{}{},
					"computed_map_of_floats":  map[string]interface{}{},
					"computed_map_of_ints":    map[string]interface{}{},
					"computed_map_of_strings": map[string]interface{}{},
				},
			},
		},
	}.test(t)
}

func TestResourceEncode_NestedOneLevelDeepSingleMultiple(t *testing.T) {
	type Inner struct {
		Value string `tfschema:"value"`
	}
	type Type struct {
		NestedObject []Inner `tfschema:"inner"`
	}
	encodeTestData{
		Input: &Type{
			NestedObject: []Inner{
				{
					Value: "first",
				},
				{
					Value: "second",
				},
				{
					Value: "third",
				},
			},
		},
		Expected: map[string]interface{}{
			"inner": []interface{}{
				map[string]interface{}{
					"value": "first",
				},
				map[string]interface{}{
					"value": "second",
				},
				map[string]interface{}{
					"value": "third",
				},
			},
		},
	}.test(t)
}

func TestResourceEncode_NestedThreeLevelsDeepEmpty(t *testing.T) {
	type ThirdInner struct {
		Value string `tfschema:"value"`
	}
	type SecondInner struct {
		Third []ThirdInner `tfschema:"third"`
	}
	type FirstInner struct {
		Second []SecondInner `tfschema:"second"`
	}
	type Type struct {
		First []FirstInner `tfschema:"first"`
	}

	t.Log("Top Level Empty")
	encodeTestData{
		Input: &Type{
			First: []FirstInner{},
		},
		Expected: map[string]interface{}{
			"first": []interface{}{},
		},
	}.test(t)

	t.Log("Second Level Empty")
	encodeTestData{
		Input: &Type{
			First: []FirstInner{
				{
					Second: []SecondInner{},
				},
			},
		},
		Expected: map[string]interface{}{
			"first": []interface{}{
				map[string]interface{}{
					"second": []interface{}{},
				},
			},
		},
	}.test(t)

	t.Log("Third Level Empty")
	encodeTestData{
		Input: &Type{
			First: []FirstInner{
				{
					Second: []SecondInner{
						{
							Third: []ThirdInner{},
						},
					},
				},
			},
		},
		Expected: map[string]interface{}{
			"first": []interface{}{
				map[string]interface{}{
					"second": []interface{}{
						map[string]interface{}{
							"third": []interface{}{},
						},
					},
				},
			},
		},
	}.test(t)
}

func TestResourceEncode_NestedThreeLevelsDeepSingleItem(t *testing.T) {
	type ThirdInner struct {
		Value string `tfschema:"value"`
	}
	type SecondInner struct {
		Third []ThirdInner `tfschema:"third"`
	}
	type FirstInner struct {
		Second []SecondInner `tfschema:"second"`
	}
	type Type struct {
		First []FirstInner `tfschema:"first"`
	}

	encodeTestData{
		Input: &Type{
			First: []FirstInner{
				{
					Second: []SecondInner{
						{
							Third: []ThirdInner{
								{
									Value: "salut",
								},
							},
						},
					},
				},
			},
		},
		Expected: map[string]interface{}{
			"first": []interface{}{
				map[string]interface{}{
					"second": []interface{}{
						map[string]interface{}{
							"third": []interface{}{
								map[string]interface{}{
									"value": "salut",
								},
							},
						},
					},
				},
			},
		},
	}.test(t)
}

func TestResourceEncode_NestedThreeLevelsDeepMultipleItems(t *testing.T) {
	type ThirdInner struct {
		Value string `tfschema:"value"`
	}
	type SecondInner struct {
		Value string       `tfschema:"value"`
		Third []ThirdInner `tfschema:"third"`
	}
	type FirstInner struct {
		Value  string        `tfschema:"value"`
		Second []SecondInner `tfschema:"second"`
	}
	type Type struct {
		First []FirstInner `tfschema:"first"`
	}

	encodeTestData{
		Input: &Type{
			First: []FirstInner{
				{
					Value: "first - 1",
					Second: []SecondInner{
						{
							Value: "second - 1",
							Third: []ThirdInner{
								{
									Value: "third - 1",
								},
								{
									Value: "third - 2",
								},
								{
									Value: "third - 3",
								},
							},
						},
						{
							Value: "second - 2",
							Third: []ThirdInner{
								{
									Value: "third - 4",
								},
								{
									Value: "third - 5",
								},
								{
									Value: "third - 6",
								},
							},
						},
					},
				},
				{
					Value: "first - 2",
					Second: []SecondInner{
						{
							Value: "second - 3",
							Third: []ThirdInner{
								{
									Value: "third - 7",
								},
								{
									Value: "third - 8",
								},
							},
						},
						{
							Value: "second - 4",
							Third: []ThirdInner{
								{
									Value: "third - 9",
								},
							},
						},
					},
				},
			},
		},
		Expected: map[string]interface{}{
			"first": []interface{}{
				map[string]interface{}{
					"value": "first - 1",
					"second": []interface{}{
						map[string]interface{}{
							"value": "second - 1",
							"third": []interface{}{
								map[string]interface{}{
									"value": "third - 1",
								},
								map[string]interface{}{
									"value": "third - 2",
								},
								map[string]interface{}{
									"value": "third - 3",
								},
							},
						},
						map[string]interface{}{
							"value": "second - 2",
							"third": []interface{}{
								map[string]interface{}{
									"value": "third - 4",
								},
								map[string]interface{}{
									"value": "third - 5",
								},
								map[string]interface{}{
									"value": "third - 6",
								},
							},
						},
					},
				},
				map[string]interface{}{
					"value": "first - 2",
					"second": []interface{}{
						map[string]interface{}{
							"value": "second - 3",
							"third": []interface{}{
								map[string]interface{}{
									"value": "third - 7",
								},
								map[string]interface{}{
									"value": "third - 8",
								},
							},
						},
						map[string]interface{}{
							"value": "second - 4",
							"third": []interface{}{
								map[string]interface{}{
									"value": "third - 9",
								},
							},
						},
					},
				},
			},
		},
	}.test(t)
}

func TestResourceEncode_NestedThreeLevelsDeepMultipleOptionalItems(t *testing.T) {
	type ThirdInner struct {
		Value string `tfschema:"value"`
	}
	type SecondInner struct {
		Value *string       `tfschema:"value"`
		Third *[]ThirdInner `tfschema:"third"`
	}
	type FirstInner struct {
		Value  *string        `tfschema:"value"`
		Second *[]SecondInner `tfschema:"second"`
	}
	type Type struct {
		First *[]FirstInner `tfschema:"first"`
	}

	encodeTestData{
		Input: &Type{
			First: &[]FirstInner{
				{
					Value: pointer.To("first - 1"),
					Second: &[]SecondInner{
						{
							Value: pointer.To("second - 1"),
							Third: &[]ThirdInner{
								{
									Value: "third - 1",
								},
								{
									Value: "third - 2",
								},
								{
									Value: "third - 3",
								},
							},
						},
						{
							Value: pointer.To("second - 2"),
							Third: &[]ThirdInner{
								{
									Value: "third - 4",
								},
								{
									Value: "third - 5",
								},
								{
									Value: "third - 6",
								},
							},
						},
					},
				},
				{
					Value: pointer.To("first - 2"),
					Second: &[]SecondInner{
						{
							Value: pointer.To("second - 3"),
							Third: &[]ThirdInner{
								{
									Value: "third - 7",
								},
								{
									Value: "third - 8",
								},
							},
						},
						{
							Value: pointer.To("second - 4"),
							Third: &[]ThirdInner{
								{
									Value: "third - 9",
								},
							},
						},
					},
				},
			},
		},
		Expected: map[string]interface{}{
			"first": []interface{}{
				map[string]interface{}{
					"value": "first - 1",
					"second": []interface{}{
						map[string]interface{}{
							"value": "second - 1",
							"third": []interface{}{
								map[string]interface{}{
									"value": "third - 1",
								},
								map[string]interface{}{
									"value": "third - 2",
								},
								map[string]interface{}{
									"value": "third - 3",
								},
							},
						},
						map[string]interface{}{
							"value": "second - 2",
							"third": []interface{}{
								map[string]interface{}{
									"value": "third - 4",
								},
								map[string]interface{}{
									"value": "third - 5",
								},
								map[string]interface{}{
									"value": "third - 6",
								},
							},
						},
					},
				},
				map[string]interface{}{
					"value": "first - 2",
					"second": []interface{}{
						map[string]interface{}{
							"value": "second - 3",
							"third": []interface{}{
								map[string]interface{}{
									"value": "third - 7",
								},
								map[string]interface{}{
									"value": "third - 8",
								},
							},
						},
						map[string]interface{}{
							"value": "second - 4",
							"third": []interface{}{
								map[string]interface{}{
									"value": "third - 9",
								},
							},
						},
					},
				},
			},
		},
	}.test(t)
}

func (testData encodeTestData) test(t *testing.T) {
	objType := reflect.TypeOf(testData.Input).Elem()
	objVal := reflect.ValueOf(testData.Input).Elem()
	debugLogger := ConsoleLogger{}

	output, err := recurse(objType, objVal, debugLogger)
	if err != nil {
		if testData.ExpectError {
			// we're good
			return
		}

		t.Fatalf("encoding error: %+v", err)
	}
	if testData.ExpectError {
		t.Fatalf("expected an error but didn't get one!")
	}

	if diff := cmp.Diff(output, testData.Expected); diff != "" {
		t.Fatalf("Output mismatch, diff:\n\n %s", diff)
	}
}
