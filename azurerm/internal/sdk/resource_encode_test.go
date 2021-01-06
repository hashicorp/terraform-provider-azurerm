package sdk

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type encodeTestData struct {
	Input       interface{}
	Expected    map[string]interface{}
	ExpectError bool
}

func TestResourceEncode_TopLevel(t *testing.T) {
	type SimpleType struct {
		String        string            `tfschema:"string"`
		Number        int               `tfschema:"number"`
		Price         float64           `tfschema:"price"`
		Enabled       bool              `tfschema:"enabled"`
		ListOfFloats  []float64         `tfschema:"list_of_floats"`
		ListOfNumbers []int             `tfschema:"list_of_numbers"`
		ListOfStrings []string          `tfschema:"list_of_strings"`
		MapOfBools    map[string]bool   `tfschema:"map_of_bools"`
		MapOfNumbers  map[string]int    `tfschema:"map_of_numbers"`
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
				1.234567890},
			ListOfNumbers: []int{1, 2, 3},
			ListOfStrings: []string{
				"have",
				"you",
				"heard",
			},
			MapOfBools: map[string]bool{
				"awesome_feature": true,
			},
			MapOfNumbers: map[string]int{
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
			"list_of_numbers": []int{1, 2, 3},
			"list_of_strings": []string{
				"have",
				"you",
				"heard",
			},
			"map_of_bools": map[string]interface{}{
				"awesome_feature": true,
			},
			"map_of_numbers": map[string]interface{}{
				"hello": 1,
				"there": 3,
			},
			"map_of_strings": map[string]interface{}{
				"hello":   "there",
				"salut":   "tous les monde",
				"guten":   "tag",
				"morning": "alvaro",
			},
		},
	}.test(t)
}

func TestResourceEncode_TopLevelOmitted(t *testing.T) {
	type SimpleType struct {
		String        string            `tfschema:"string"`
		Number        int               `tfschema:"number"`
		Price         float64           `tfschema:"price"`
		Enabled       bool              `tfschema:"enabled"`
		ListOfFloats  []float64         `tfschema:"list_of_floats"`
		ListOfNumbers []int             `tfschema:"list_of_numbers"`
		ListOfStrings []string          `tfschema:"list_of_strings"`
		MapOfBools    map[string]bool   `tfschema:"map_of_bools"`
		MapOfNumbers  map[string]int    `tfschema:"map_of_numbers"`
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
			"list_of_numbers": []int{},
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
		ComputedNumber        int                `tfschema:"computed_number" computed:"true"`
		ComputedBool          bool               `tfschema:"computed_bool" computed:"true"`
		ComputedListOfNumbers []int              `tfschema:"computed_list_of_numbers" computed:"true"`
		ComputedListOfStrings []string           `tfschema:"computed_list_of_strings" computed:"true"`
		ComputedMapOfBools    map[string]bool    `tfschema:"computed_map_of_bools" computed:"true"`
		ComputedMapOfFloats   map[string]float64 `tfschema:"computed_map_of_floats" computed:"true"`
		ComputedMapOfInts     map[string]int     `tfschema:"computed_map_of_ints" computed:"true"`
		ComputedMapOfStrings  map[string]string  `tfschema:"computed_map_of_strings" computed:"true"`
	}
	encodeTestData{
		Input: &SimpleType{
			ComputedString:        "je suis computed",
			ComputedNumber:        732,
			ComputedBool:          true,
			ComputedListOfNumbers: []int{1, 2, 3},
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
			ComputedMapOfInts: map[string]int{
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
			"computed_list_of_numbers": []int{1, 2, 3},
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
				"first":  1,
				"second": 2,
				"third":  3,
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
		Number        int               `tfschema:"number"`
		Price         float64           `tfschema:"price"`
		Enabled       bool              `tfschema:"enabled"`
		ListOfFloats  []float64         `tfschema:"list_of_floats"`
		ListOfNumbers []int             `tfschema:"list_of_numbers"`
		ListOfStrings []string          `tfschema:"list_of_strings"`
		MapOfBools    map[string]bool   `tfschema:"map_of_bools"`
		MapOfNumbers  map[string]int    `tfschema:"map_of_numbers"`
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
						1.234567890},
					ListOfNumbers: []int{1, 2, 3},
					ListOfStrings: []string{
						"have",
						"you",
						"heard",
					},
					MapOfBools: map[string]bool{
						"awesome_feature": true,
					},
					MapOfNumbers: map[string]int{
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
					"list_of_numbers": []int{1, 2, 3},
					"list_of_strings": []string{
						"have",
						"you",
						"heard",
					},
					"map_of_bools": map[string]interface{}{
						"awesome_feature": true,
					},
					"map_of_numbers": map[string]interface{}{
						"hello": 1,
						"there": 3,
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
		Number               int                `tfschema:"number"`
		Price                float64            `tfschema:"price"`
		Enabled              bool               `tfschema:"enabled"`
		ListOfFloats         []float64          `tfschema:"list_of_floats"`
		ListOfNumbers        []int              `tfschema:"list_of_numbers"`
		ListOfStrings        []string           `tfschema:"list_of_strings"`
		MapOfBools           map[string]bool    `tfschema:"map_of_bools"`
		MapOfNumbers         map[string]int     `tfschema:"map_of_numbers"`
		MapOfStrings         map[string]string  `tfschema:"map_of_strings"`
		ComputedMapOfBools   map[string]bool    `tfschema:"computed_map_of_bools" computed:"true"`
		ComputedMapOfFloats  map[string]float64 `tfschema:"computed_map_of_floats" computed:"true"`
		ComputedMapOfInts    map[string]int     `tfschema:"computed_map_of_ints" computed:"true"`
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
					"list_of_numbers":         []int{},
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

func (testData encodeTestData) test(t *testing.T) {
	objType := reflect.TypeOf(testData.Input).Elem()
	objVal := reflect.ValueOf(testData.Input).Elem()
	fieldName := reflect.ValueOf(testData.Input).Elem().String()
	debugLogger := ConsoleLogger{}

	output, err := recurse(objType, objVal, fieldName, debugLogger)
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

	if !cmp.Equal(output, testData.Expected) {
		t.Fatalf("Output mismatch:\n\n Expected: %+v\n\n Received: %+v\n\n", testData.Expected, output)
	}
}
