// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

type decodeTestData struct {
	State       map[string]interface{}
	Input       interface{}
	Expected    interface{}
	ExpectError bool
}

type AllRequired struct {
	String        string             `tfschema:"string"`
	Int64         int64              `tfschema:"int64"`
	Float         float64            `tfschema:"float"`
	Enabled       bool               `tfschema:"enabled"`
	ListOfFloat   []float64          `tfschema:"list_of_float"`
	ListOfInt64   []int64            `tfschema:"list_of_int64"`
	ListOfStrings []string           `tfschema:"list_of_strings"`
	MapOfBooleans map[string]bool    `tfschema:"map_of_booleans"`
	MapOfFloat    map[string]float64 `tfschema:"map_of_float"`
	MapOfInt64    map[string]int64   `tfschema:"map_of_int64"`
	MapOfStrings  map[string]string  `tfschema:"map_of_strings"`
}

type OneRequiredRestOptional struct {
	Required      string              `tfschema:"required"`
	String        *string             `tfschema:"string"`
	Int64         *int64              `tfschema:"int64"`
	Float         *float64            `tfschema:"float"`
	EmptyBoolean  bool                `tfschema:"empty_boolean"`
	Enabled       *bool               `tfschema:"enabled"`
	ListOfFloat   *[]float64          `tfschema:"list_of_float"`
	ListOfInt64   *[]int64            `tfschema:"list_of_int64"`
	ListOfStrings *[]string           `tfschema:"list_of_strings"`
	MapOfBooleans *map[string]bool    `tfschema:"map_of_booleans"`
	MapOfFloat    *map[string]float64 `tfschema:"map_of_float"`
	MapOfInt64    *map[string]int64   `tfschema:"map_of_int64"`
	MapOfStrings  *map[string]string  `tfschema:"map_of_strings"`
}

type OneOfEverything struct {
	RequiredStr           string              `tfschema:"required_str"`
	OptionalStr           *string             `tfschema:"optional_str"`
	RequiredInt64         int64               `tfschema:"required_int64"`
	OptionalInt64         *int64              `tfschema:"optional_int64"`
	RequiredFloat         float64             `tfschema:"required_float"`
	OptionalFloat         *float64            `tfschema:"optional_float"`
	RequiredBoolean       bool                `tfschema:"required_boolean"`
	OptionalBoolean       *bool               `tfschema:"optional_boolean"`
	RequiredListOfFloat   []float64           `tfschema:"required_list_of_float"`
	OptionalListOfFloat   *[]float64          `tfschema:"optional_list_of_float"`
	RequiredListOfInt64   []int64             `tfschema:"required_list_of_int64"`
	OptionalListOfInt64   *[]int64            `tfschema:"optional_list_of_int64"`
	RequiredListOfStrings []string            `tfschema:"required_list_of_strings"`
	OptionalListOfStrings *[]string           `tfschema:"optional_list_of_strings"`
	RequiredMapOfBooleans map[string]bool     `tfschema:"required_map_of_booleans"`
	OptionalMapOfBooleans *map[string]bool    `tfschema:"optional_map_of_booleans"`
	RequiredMapOfFloat    map[string]float64  `tfschema:"required_map_of_float"`
	OptionalMapOfFloat    *map[string]float64 `tfschema:"optional_map_of_float"`
	RequiredMapOfInt64    map[string]int64    `tfschema:"required_map_of_int64"`
	OptionalMapOfInt64    *map[string]int64   `tfschema:"optional_map_of_int64"`
	RequiredMapOfStrings  map[string]string   `tfschema:"required_map_of_strings"`
	OptionalMapOfStrings  *map[string]string  `tfschema:"optional_map_of_strings"`
}

func TestDecode_TopLevelFieldsRequired(t *testing.T) {
	decodeTestData{
		State: map[string]interface{}{
			"int64":   42,
			"float":   129.99,
			"string":  "world",
			"enabled": true,
			"list_of_float": []interface{}{
				1.0,
				2.0,
				3.0,
				1.234567890,
			},
			"list_of_int64": []interface{}{1, 2, 3},
			"list_of_strings": []interface{}{
				"have",
				"you",
				"heard",
			},
			"map_of_booleans": map[string]interface{}{
				"awesome_feature": true,
			},
			"map_of_int": map[string]interface{}{
				"hello": 1,
				"there": 3,
			},
			"map_of_int64": map[string]interface{}{
				"ten":    10,
				"eleven": 11,
			},
			"map_of_float": map[string]interface{}{
				"pi":    3.12,
				"fifth": 0.2,
			},
			"map_of_strings": map[string]interface{}{
				"hello":   "there",
				"salut":   "tous les monde",
				"guten":   "tag",
				"morning": "alvaro",
			},
		},
		Input: &AllRequired{},
		Expected: &AllRequired{
			String:  "world",
			Float:   129.99,
			Enabled: true,
			Int64:   42,
			ListOfFloat: []float64{
				1.0,
				2.0,
				3.0,
				1.234567890,
			},
			ListOfInt64: []int64{1, 2, 3},
			ListOfStrings: []string{
				"have",
				"you",
				"heard",
			},
			MapOfBooleans: map[string]bool{
				"awesome_feature": true,
			},
			MapOfFloat: map[string]float64{
				"pi":    3.12,
				"fifth": 0.2,
			},
			MapOfInt64: map[string]int64{
				"ten":    10,
				"eleven": 11,
			},
			MapOfStrings: map[string]string{
				"hello":   "there",
				"salut":   "tous les monde",
				"guten":   "tag",
				"morning": "alvaro",
			},
		},
		ExpectError: false,
	}.test(t)
}

func TestDecode_TopLevelFieldsComputedNoValues(t *testing.T) {
	// NOTE: this scenario covers Create without any existing Computed values
	type SimpleType struct {
		ComputedMapOfBools   map[string]bool    `tfschema:"computed_map_of_bools"`
		ComputedMapOfFloats  map[string]float64 `tfschema:"computed_map_of_floats"`
		ComputedMapOfInts    map[string]int64   `tfschema:"computed_map_of_ints"`
		ComputedMapOfStrings map[string]string  `tfschema:"computed_map_of_strings"`
	}
	decodeTestData{
		State: map[string]interface{}{
			"computed_map_of_bools":   map[string]interface{}{},
			"computed_map_of_floats":  map[string]interface{}{},
			"computed_map_of_ints":    map[string]interface{}{},
			"computed_map_of_strings": map[string]interface{}{},
		},
		Input: &SimpleType{},
		Expected: &SimpleType{
			ComputedMapOfBools:   map[string]bool{},
			ComputedMapOfFloats:  map[string]float64{},
			ComputedMapOfInts:    map[string]int64{},
			ComputedMapOfStrings: map[string]string{},
		},
		ExpectError: false,
	}.test(t)
}

func TestDecode_TopLevelFieldsComputedWithValues(t *testing.T) {
	// NOTE: this scenario covers Update/Read with existing Computed values or Computed/Optional
	type SimpleType struct {
		ComputedMapOfBools   map[string]bool    `tfschema:"computed_map_of_bools"`
		ComputedMapOfFloats  map[string]float64 `tfschema:"computed_map_of_floats"`
		ComputedMapOfInts    map[string]int64   `tfschema:"computed_map_of_ints"`
		ComputedMapOfStrings map[string]string  `tfschema:"computed_map_of_strings"`
	}
	decodeTestData{
		State: map[string]interface{}{
			"computed_map_of_bools": map[string]interface{}{
				"bingo": true,
				"bango": false,
			},
			"computed_map_of_floats": map[string]interface{}{
				"bingo": -2.197234,
				"bango": 3.123456789,
			},
			"computed_map_of_ints": map[string]interface{}{
				"bingo": 2197234,
				"bango": 3123456789,
			},
			"computed_map_of_strings": map[string]interface{}{
				"matthew": "brisket",
				"tom":     "coffee",
			},
		},
		Input: &SimpleType{},
		Expected: &SimpleType{
			ComputedMapOfBools: map[string]bool{
				"bingo": true,
				"bango": false,
			},
			ComputedMapOfFloats: map[string]float64{
				"bingo": -2.197234,
				"bango": 3.123456789,
			},
			ComputedMapOfInts: map[string]int64{
				"bingo": 2197234,
				"bango": 3123456789,
			},
			ComputedMapOfStrings: map[string]string{
				"matthew": "brisket",
				"tom":     "coffee",
			},
		},
		ExpectError: false,
	}.test(t)
}

func TestDecode_TopLevelFieldsOptional(t *testing.T) {
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
	decodeTestData{
		State: map[string]interface{}{
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
		Input: &SimpleType{},
		Expected: &SimpleType{
			MapOfBools:   map[string]bool{},
			MapOfNumbers: map[string]int64{},
			MapOfStrings: map[string]string{},
		},
		ExpectError: false,
	}.test(t)
}

func TestDecode_TopLevelFieldsOptionalNullValues(t *testing.T) {
	decodeTestData{
		State: map[string]interface{}{
			"required": "name",
		},
		Input: &OneRequiredRestOptional{},
		Expected: &OneRequiredRestOptional{
			Required:     "name",
			EmptyBoolean: false,
		},
		ExpectError: false,
	}.test(t)
}

func TestDecode_TopLevelFieldsOptionalMixedValues(t *testing.T) {
	decodeTestData{
		State: map[string]interface{}{
			"required": "name",
			"float":    3.5,
			"enabled":  false,
			"list_of_strings": []interface{}{
				"have",
				"you",
				"heard",
			},
			"map_of_int64": &map[string]interface{}{
				"ten":       10, // TODO - Should this be needed?
				"twentyone": 21,
			},
		},
		Input: &OneRequiredRestOptional{},
		Expected: &OneRequiredRestOptional{
			Required: "name",
			Float:    pointer.To(3.5),
			Enabled:  pointer.To(false),
			ListOfStrings: pointer.To([]string{
				"have",
				"you",
				"heard",
			}),
			MapOfInt64: pointer.To(map[string]int64{
				"ten":       10,
				"twentyone": 21,
			}),
		},
		ExpectError: false,
	}.test(t)
}

func TestDecode_TopLevelFieldsOptionalComplete(t *testing.T) {
	decodeTestData{
		State: map[string]interface{}{
			"required":   "name",
			"int64":      1984,
			"float":      3.5,
			"string":     "do you know where your towel is",
			"enabled":    true,
			"empty_bool": nil,
			"list_of_float": []interface{}{
				1.2,
				3.142,
			},
			"list_of_int64": []interface{}{
				100,
				200,
			},
			"list_of_strings": []interface{}{
				"have",
				"you",
				"heard",
			},
			"map_of_booleans": map[string]interface{}{
				"first":  true,
				"second": false,
			},
			"map_of_int64": map[string]interface{}{
				"Orwell":     1984,
				"fahrenheit": 451,
			},
			"map_of_strings": map[string]interface{}{
				"foo":  "bar",
				"ford": "prefect",
			},
			"map_of_float": map[string]interface{}{
				"pi":    3.12,
				"fifth": 0.2,
			},
		},
		Input: &OneRequiredRestOptional{},
		Expected: &OneRequiredRestOptional{
			Required:     "name",
			String:       pointer.To("do you know where your towel is"),
			Float:        pointer.To(3.5),
			Int64:        pointer.To(int64(1984)),
			Enabled:      pointer.To(true),
			EmptyBoolean: false,
			ListOfFloat: pointer.To([]float64{
				1.2,
				3.142,
			}),
			ListOfInt64: pointer.To([]int64{
				100,
				200,
			}),
			ListOfStrings: pointer.To([]string{
				"have",
				"you",
				"heard",
			}),
			MapOfInt64: pointer.To(map[string]int64{
				"Orwell":     1984,
				"fahrenheit": 451,
			}),
			MapOfBooleans: pointer.To(map[string]bool{
				"first":  true,
				"second": false,
			}),
			MapOfStrings: pointer.To(map[string]string{
				"foo":  "bar",
				"ford": "prefect",
			}),
			MapOfFloat: pointer.To(map[string]float64{
				"pi":    3.12,
				"fifth": 0.2,
			}),
		},
		ExpectError: false,
	}.test(t)
}

func TestDecode_TopLevelFieldsComputed(t *testing.T) {
	type SimpleType struct {
		ComputedString        string   `tfschema:"computed_string"`
		ComputedNumber        int64    `tfschema:"computed_number"`
		ComputedBool          bool     `tfschema:"computed_bool"`
		ComputedListOfNumbers []int64  `tfschema:"computed_list_of_numbers"`
		ComputedListOfStrings []string `tfschema:"computed_list_of_strings"`
	}
	decodeTestData{
		State: map[string]interface{}{
			"computed_string":          "je suis computed",
			"computed_number":          732,
			"computed_bool":            true,
			"computed_list_of_numbers": []interface{}{1, 2, 3},
			"computed_list_of_strings": []interface{}{
				"have",
				"you",
				"heard",
			},
		},
		Input: &SimpleType{},
		Expected: &SimpleType{
			ComputedString:        "je suis computed",
			ComputedNumber:        732,
			ComputedBool:          true,
			ComputedListOfNumbers: []int64{1, 2, 3},
			ComputedListOfStrings: []string{
				"have",
				"you",
				"heard",
			},
		},
		ExpectError: false,
	}.test(t)
}

func TestResourceDecode_NestedOneLevelDeepEmpty(t *testing.T) {
	type Inner struct {
		Value string `tfschema:"value"`
	}
	type Type struct {
		NestedObject []Inner `tfschema:"inner"`
	}
	decodeTestData{
		State: map[string]interface{}{
			"inner": []interface{}{},
		},
		Input: &Type{},
		Expected: &Type{
			NestedObject: []Inner{},
		},
	}.test(t)
}

func TestResourceDecode_NestedOneLevelDeepSingle(t *testing.T) {
	type Inner struct {
		String        string            `tfschema:"string"`
		Number        int64             `tfschema:"number"`
		Price         float64           `tfschema:"price"`
		Enabled       bool              `tfschema:"enabled"`
		ListOfFloats  []float64         `tfschema:"list_of_floats"`
		ListOfInt64   []int64           `tfschema:"list_of_int64"`
		ListOfStrings []string          `tfschema:"list_of_strings"`
		MapOfBools    map[string]bool   `tfschema:"map_of_bools"`
		MapOfInt64    map[string]int64  `tfschema:"map_of_int64"`
		MapOfStrings  map[string]string `tfschema:"map_of_strings"`
	}
	type Type struct {
		NestedObject []Inner `tfschema:"inner"`
	}
	decodeTestData{
		State: map[string]interface{}{
			"inner": []interface{}{
				map[string]interface{}{
					"number":  42,
					"price":   129.99,
					"string":  "world",
					"enabled": true,
					"list_of_floats": []interface{}{
						1.0,
						2.0,
						3.0,
						1.234567890,
					},
					"list_of_int64": []interface{}{2, 4, 6},
					"list_of_strings": []interface{}{
						"have",
						"you",
						"heard",
					},
					"map_of_bools": map[string]interface{}{
						"awesome_feature": true,
					},
					"map_of_int64": map[string]interface{}{
						"hello": 2,
						"there": 6,
					},
					"map_of_strings": map[string]interface{}{
						"hello":   "there",
						"salut":   "tout les monde",
						"guten":   "tag",
						"morning": "alvaro",
					},
				},
			},
		},
		Input: &Type{},
		Expected: &Type{
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
					ListOfInt64: []int64{2, 4, 6},
					ListOfStrings: []string{
						"have",
						"you",
						"heard",
					},
					MapOfBools: map[string]bool{
						"awesome_feature": true,
					},
					MapOfInt64: map[string]int64{
						"hello": 2,
						"there": 6,
					},
					MapOfStrings: map[string]string{
						"hello":   "there",
						"salut":   "tout les monde",
						"guten":   "tag",
						"morning": "alvaro",
					},
				},
			},
		},
	}.test(t)
}

func TestResourceDecode_NestedOneLevelDeepSingleOmittedValues(t *testing.T) {
	type Inner struct {
		String        string            `tfschema:"string"`
		Number        int64             `tfschema:"number"`
		Price         float64           `tfschema:"price"`
		Enabled       bool              `tfschema:"enabled"`
		ListOfFloats  []float64         `tfschema:"list_of_floats"`
		ListOfInt64   []int64           `tfschema:"list_of_int64"`
		ListOfStrings []string          `tfschema:"list_of_strings"`
		MapOfBools    map[string]bool   `tfschema:"map_of_bools"`
		MapOfInt64    map[string]int64  `tfschema:"map_of_int64"`
		MapOfStrings  map[string]string `tfschema:"map_of_strings"`
	}
	type Type struct {
		NestedObject []Inner `tfschema:"inner"`
	}
	decodeTestData{
		State: map[string]interface{}{
			"inner": []interface{}{
				map[string]interface{}{
					"number":          0,
					"price":           float64(0),
					"string":          "",
					"enabled":         false,
					"list_of_floats":  []float64{},
					"list_of_int64":   []int64{},
					"list_of_strings": []string{},
					"map_of_bools":    map[string]interface{}{},
					"map_of_numbers":  map[string]interface{}{},
					"map_of_int64":    map[string]interface{}{},
					"map_of_strings":  map[string]interface{}{},
				},
			},
		},
		Input: &Type{},
		Expected: &Type{
			NestedObject: []Inner{
				{
					MapOfBools:   map[string]bool{},
					MapOfStrings: map[string]string{},
					MapOfInt64:   map[string]int64{},
				},
			},
		},
	}.test(t)
}

func TestResourceDecode_NestedOneLevelDeepSingleMultiple(t *testing.T) {
	type Inner struct {
		Value string `tfschema:"value"`
	}
	type Type struct {
		NestedObject []Inner `tfschema:"inner"`
	}
	decodeTestData{
		State: map[string]interface{}{
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
		Input: &Type{},
		Expected: &Type{
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
	}.test(t)
}

func TestResourceDecode_NestedThreeLevelsDeepEmpty(t *testing.T) {
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
	decodeTestData{
		State: map[string]interface{}{
			"first": []interface{}{},
		},
		Input: &Type{},
		Expected: &Type{
			First: []FirstInner{},
		},
	}.test(t)

	t.Log("Second Level Empty")
	decodeTestData{
		State: map[string]interface{}{
			"first": []interface{}{
				map[string]interface{}{
					"second": []interface{}{},
				},
			},
		},
		Input: &Type{},
		Expected: &Type{
			First: []FirstInner{
				{
					Second: []SecondInner{},
				},
			},
		},
	}.test(t)

	t.Log("Third Level Empty")
	decodeTestData{
		State: map[string]interface{}{
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
		Input: &Type{},
		Expected: &Type{
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
	}.test(t)
}

func TestResourceDecode_NestedThreeLevelsDeepSingleItem(t *testing.T) {
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

	decodeTestData{
		State: map[string]interface{}{
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
		Input: &Type{},
		Expected: &Type{
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
	}.test(t)
}

func TestResourceDecode_NestedThreeLevelsDeepMultipleItems(t *testing.T) {
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

	decodeTestData{
		State: map[string]interface{}{
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
		Input: &Type{},
		Expected: &Type{
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
	}.test(t)
}

func TestResourceDecode_NestedThreeLevelsDeepMultipleOptionalItems(t *testing.T) {
	type ThirdInner struct {
		Value string `tfschema:"value"`
	}
	type SecondInner struct {
		Value  *string       `tfschema:"value"`
		Number *int64        `tfschema:"number"`
		Third  *[]ThirdInner `tfschema:"third"`
	}
	type FirstInner struct {
		Value  *string        `tfschema:"value"`
		Second *[]SecondInner `tfschema:"second"`
	}
	type Type struct {
		First *[]FirstInner `tfschema:"first"`
	}

	decodeTestData{
		State: map[string]interface{}{
			"first": []interface{}{
				map[string]interface{}{
					"value": "first - 1",
					"second": []interface{}{
						map[string]interface{}{
							"value":  "second - 1",
							"number": 2,
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
		Input: &Type{},
		Expected: &Type{
			First: &[]FirstInner{
				{
					Value: pointer.To("first - 1"),
					Second: &[]SecondInner{
						{
							Value:  pointer.To("second - 1"),
							Number: pointer.To(int64(2)),
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
	}.test(t)
}

func (testData decodeTestData) test(t *testing.T) {
	debugLogger := ConsoleLogger{}
	state := testData.stateWrapper()
	if err := decodeReflectedType(testData.Input, state, debugLogger); err != nil {
		if testData.ExpectError {
			// we're good
			return
		}

		t.Fatalf("unexpected error: %+v", err)
	}
	if testData.ExpectError {
		t.Fatalf("expected an error but didn't get one!")
	}

	if diff := cmp.Diff(testData.Input, testData.Expected); diff != "" {
		t.Fatalf("Output mismatch, diff:\n\n %s", diff)
	}
}

func (testData decodeTestData) stateWrapper() testDataGetter {
	return testDataGetter{
		values: testData.State,
	}
}

type testDataGetter struct {
	values map[string]interface{}
}

func (td testDataGetter) Get(key string) interface{} {
	return td.values[key]
}

func (td testDataGetter) GetOk(key string) (interface{}, bool) {
	val, ok := td.values[key]
	return val, ok
}

func (td testDataGetter) GetOkExists(key string) (interface{}, bool) {
	// for the purposes of this test this should be sufficient, maybe?
	val, ok := td.values[key]
	return val, ok
}
