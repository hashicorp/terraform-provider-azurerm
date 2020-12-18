package sdk

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPluginSDKAndDecoder(t *testing.T) {
	os.Setenv("TF_ACC", "1")

	type NestedType struct {
		Key string `tfschema:"key"`
	}
	type MyType struct {
		Hello         string             `tfschema:"hello"`
		RandomNumber  int                `tfschema:"random_number"`
		Enabled       bool               `tfschema:"enabled"`
		ListOfStrings []string           `tfschema:"list_of_strings"`
		ListOfNumbers []int              `tfschema:"list_of_numbers"`
		ListOfBools   []bool             `tfschema:"list_of_bools"`
		ListOfFloats  []float64          `tfschema:"list_of_floats"`
		NestedObject  []NestedType       `tfschema:"nested_object"`
		MapOfStrings  map[string]string  `tfschema:"map_of_strings"`
		MapOfNumbers  map[string]int     `tfschema:"map_of_numbers"`
		MapOfBools    map[string]bool    `tfschema:"map_of_bools"`
		MapOfFloats   map[string]float64 `tfschema:"map_of_floats"`
		// Sets are handled in a separate test, since the orders can be different
	}

	expected := MyType{
		Hello:         "world",
		RandomNumber:  42,
		Enabled:       true,
		ListOfStrings: []string{"hello", "there"},
		ListOfNumbers: []int{1, 2, 4},
		ListOfBools:   []bool{true, false},
		ListOfFloats:  []float64{-1.234567894321, 2.3456789},
		NestedObject: []NestedType{
			{
				Key: "value",
			},
		},
		MapOfStrings: map[string]string{
			"bingo": "bango",
		},
		MapOfNumbers: map[string]int{
			"lucky": 21,
		},
		MapOfBools: map[string]bool{
			"friday": true,
		},
		MapOfFloats: map[string]float64{
			"pi": 3.14159,
		},
	}

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]terraform.ResourceProviderFactory{
			"validator": func() (terraform.ResourceProvider, error) {
				return &schema.Provider{
					DataSourcesMap: map[string]*schema.Resource{},
					ResourcesMap: map[string]*schema.Resource{
						"validator_decoder": {
							Schema: map[string]*schema.Schema{
								"hello": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"random_number": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"enabled": {
									Type:     schema.TypeBool,
									Computed: true,
								},
								"list_of_strings": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"list_of_numbers": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeInt,
									},
								},
								"list_of_bools": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeBool,
									},
								},
								"list_of_floats": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeFloat,
									},
								},
								"nested_object": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"key": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"map_of_strings": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"map_of_numbers": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeInt,
									},
								},
								"map_of_bools": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeBool,
									},
								},
								"map_of_floats": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeFloat,
									},
								},
							},
							Create: func(d *schema.ResourceData, i interface{}) error {
								d.SetId("some-id")
								d.Set("hello", "world")
								d.Set("random_number", 42)
								d.Set("enabled", true)
								d.Set("list_of_strings", []string{"hello", "there"})
								d.Set("list_of_numbers", []int{1, 2, 4})
								d.Set("list_of_bools", []bool{true, false})
								d.Set("list_of_floats", []float64{-1.234567894321, 2.3456789})
								d.Set("nested_object", []interface{}{
									map[string]interface{}{
										"key": "value",
									},
								})
								d.Set("map_of_strings", map[string]string{
									"bingo": "bango",
								})
								d.Set("map_of_numbers", map[string]int{
									"lucky": 21,
								})
								d.Set("map_of_bools", map[string]bool{
									"friday": true,
								})
								d.Set("map_of_floats", map[string]float64{
									"pi": 3.14159,
								})
								return nil
							},
							Read: func(d *schema.ResourceData, _ interface{}) error {
								wrapper := ResourceMetaData{
									ResourceData:             d,
									Logger:                   ConsoleLogger{},
									serializationDebugLogger: ConsoleLogger{},
								}

								var actual MyType
								if err := wrapper.Decode(&actual); err != nil {
									return fmt.Errorf("decoding: %+v", err)
								}

								if !reflect.DeepEqual(actual, expected) {
									return fmt.Errorf("Values did not match - Expected:\n%+v\n\nActual:\n%+v", expected, actual)
								}

								return nil
							},
							Delete: func(_ *schema.ResourceData, _ interface{}) error {
								return nil
							},
						},
					},
				}, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `resource "validator_decoder" "test" {}`,
			},
		},
	})
}

func TestAccPluginSDKAndDecoderOptionalComputed(t *testing.T) {
	os.Setenv("TF_ACC", "1")

	type MyType struct {
		Hello   string `tfschema:"hello"`
		Number  int    `tfschema:"number"`
		Enabled bool   `tfschema:"enabled"`
		// TODO: do we need other field types, or is this sufficient?
	}

	var commonSchema = map[string]*schema.Schema{
		"hello": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"number": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
	}
	var readFunc = func(expected MyType) func(*schema.ResourceData, interface{}) error {
		return func(d *schema.ResourceData, _ interface{}) error {
			wrapper := ResourceMetaData{
				ResourceData:             d,
				Logger:                   ConsoleLogger{},
				serializationDebugLogger: ConsoleLogger{},
			}

			var actual MyType
			if err := wrapper.Decode(&actual); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if !reflect.DeepEqual(actual, expected) {
				return fmt.Errorf("Values did not match - Expected:\n%+v\n\nActual:\n%+v", expected, actual)
			}

			return nil
		}
	}

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]terraform.ResourceProviderFactory{
			"validator": func() (terraform.ResourceProvider, error) {
				return &schema.Provider{
					DataSourcesMap: map[string]*schema.Resource{},
					ResourcesMap: map[string]*schema.Resource{
						"validator_decoder_specified": {
							Schema: commonSchema,
							Create: func(d *schema.ResourceData, i interface{}) error {
								d.SetId("some-id")
								return nil
							},
							Read: readFunc(MyType{ // expected
								Hello:   "value-from-config",
								Number:  21,
								Enabled: true,
							}),
							Delete: func(_ *schema.ResourceData, _ interface{}) error {
								return nil
							},
						},

						"validator_decoder_unspecified": {
							Schema: commonSchema,
							Create: func(d *schema.ResourceData, i interface{}) error {
								d.SetId("some-id")
								d.Set("hello", "value-from-create")
								d.Set("number", 42)
								d.Set("enabled", false)
								return nil
							},
							Read: readFunc(MyType{ // expected
								Hello:   "value-from-create",
								Number:  42,
								Enabled: false,
							}),
							Delete: func(_ *schema.ResourceData, _ interface{}) error {
								return nil
							},
						},
					},
				}, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
resource "validator_decoder_specified" "test" {
  hello   = "value-from-config"
  number  = 21
  enabled = true
}
resource "validator_decoder_unspecified" "test" {}
`,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceStateMatches("validator_decoder_specified.test", map[string]interface{}{
						"id":      "some-id",
						"enabled": "true",
						"hello":   "value-from-config",
						"number":  "21",
					}),
					testCheckResourceStateMatches("validator_decoder_unspecified.test", map[string]interface{}{
						"id":      "some-id",
						"enabled": "false",
						"hello":   "value-from-create",
						"number":  "42",
					}),
				),
			},
		},
	})
}

func TestAccPluginSDKAndDecoderOptionalComputedOverride(t *testing.T) {
	os.Setenv("TF_ACC", "1")

	type MyType struct {
		Hello   string `tfschema:"hello"`
		Number  int    `tfschema:"number"`
		Enabled bool   `tfschema:"enabled"`
		// TODO: do we need other field types, or is this sufficient?
	}

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]terraform.ResourceProviderFactory{
			"validator": func() (terraform.ResourceProvider, error) {
				return &schema.Provider{
					DataSourcesMap: map[string]*schema.Resource{},
					ResourcesMap: map[string]*schema.Resource{
						"validator_decoder_override": {
							Schema: map[string]*schema.Schema{
								"hello": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"number": {
									Type:     schema.TypeInt,
									Optional: true,
									Computed: true,
								},
								"enabled": {
									Type:     schema.TypeBool,
									Optional: true,
									Computed: true,
								},
							},
							Create: func(d *schema.ResourceData, i interface{}) error {
								d.SetId("some-id")
								d.Set("hello", "value-from-create")
								d.Set("number", 42)
								d.Set("enabled", false)
								return nil
							},
							Read: func(d *schema.ResourceData, _ interface{}) error {
								wrapper := ResourceMetaData{
									ResourceData:             d,
									Logger:                   ConsoleLogger{},
									serializationDebugLogger: ConsoleLogger{},
								}

								var actual MyType
								if err := wrapper.Decode(&actual); err != nil {
									return fmt.Errorf("decoding: %+v", err)
								}

								expected := MyType{
									Hello:   "value-from-create",
									Number:  42,
									Enabled: false,
								}

								if !reflect.DeepEqual(actual, expected) {
									return fmt.Errorf("Values did not match - Expected:\n%+v\n\nActual:\n%+v", expected, actual)
								}

								return nil
							},
							Delete: func(_ *schema.ResourceData, _ interface{}) error {
								return nil
							},
						},
					},
				}, nil
			},
		},
		Steps: []resource.TestStep{
			{
				// Apply it
				Config: `
resource "validator_decoder_override" "test" {
  hello = ""
}
`,
			},
			{
				// Then run a plan, to detect that the default value of an empty string is picked up during the Decode
				Config: `
resource "validator_decoder_override" "test" {
  hello = ""
}
`,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceStateMatches("validator_decoder_override.test", map[string]interface{}{
						"id":      "some-id",
						"enabled": "false",
						"hello":   "",
						"number":  "42",
					}),
				),
				PlanOnly: true,
			},
		},
	})
}

func TestAccPluginSDKAndDecoderSets(t *testing.T) {
	os.Setenv("TF_ACC", "1")

	type MyType struct {
		SetOfStrings []string  `tfschema:"set_of_strings"`
		SetOfNumbers []int     `tfschema:"set_of_numbers"`
		SetOfBools   []bool    `tfschema:"set_of_bools"`
		SetOfFloats  []float64 `tfschema:"set_of_floats"`
		// we could arguably extend this with nested Sets, but they're tested in the Decode function
		// so we should be covered via this test alone
	}

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]terraform.ResourceProviderFactory{
			"validator": func() (terraform.ResourceProvider, error) {
				return &schema.Provider{
					DataSourcesMap: map[string]*schema.Resource{},
					ResourcesMap: map[string]*schema.Resource{
						"validator_decoder": {
							Schema: map[string]*schema.Schema{
								"set_of_strings": {
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"set_of_numbers": {
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeInt,
									},
								},
								"set_of_bools": {
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeBool,
									},
								},
								"set_of_floats": {
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeFloat,
									},
								},
							},
							Create: func(d *schema.ResourceData, i interface{}) error {
								d.SetId("some-id")
								d.Set("set_of_strings", []string{
									"some",
									"value",
								})
								d.Set("set_of_numbers", []int{
									1,
									2,
								})
								d.Set("set_of_bools", []bool{
									true,
									false,
								})
								d.Set("set_of_floats", []float64{
									1.1,
									2.2,
								})
								return nil
							},
							Read: func(d *schema.ResourceData, _ interface{}) error {
								wrapper := ResourceMetaData{
									ResourceData:             d,
									Logger:                   ConsoleLogger{},
									serializationDebugLogger: ConsoleLogger{},
								}

								var actual MyType
								if err := wrapper.Decode(&actual); err != nil {
									return fmt.Errorf("decoding: %+v", err)
								}

								expectedStrings := []string{
									"some",
									"value",
								}
								if len(actual.SetOfStrings) != len(expectedStrings) {
									return fmt.Errorf("expected %d strings but got %d", len(expectedStrings), len(actual.SetOfStrings))
								}
								for _, v := range expectedStrings {
									exists := false
									for _, a := range actual.SetOfStrings {
										if v == a {
											exists = true
											break
										}
									}
									if !exists {
										return fmt.Errorf("expected the string %q to exist but it didn't", v)
									}
								}

								expectedNumbers := []int{
									1,
									2,
								}
								if len(actual.SetOfNumbers) != len(expectedNumbers) {
									return fmt.Errorf("expected %d ints but got %d", len(expectedNumbers), len(actual.SetOfNumbers))
								}
								for _, v := range expectedNumbers {
									exists := false
									for _, a := range actual.SetOfNumbers {
										if v == a {
											exists = true
											break
										}
									}
									if !exists {
										return fmt.Errorf("expected the number %d to exist but it didn't", v)
									}
								}

								expectedBools := []bool{
									true,
									false,
								}
								if len(actual.SetOfBools) != len(expectedBools) {
									return fmt.Errorf("expected %d bools but got %d", len(expectedBools), len(actual.SetOfBools))
								}
								for _, v := range expectedBools {
									exists := false
									for _, a := range actual.SetOfBools {
										if v == a {
											exists = true
											break
										}
									}
									if !exists {
										return fmt.Errorf("expected the bool %t to exist but it didn't", v)
									}
								}

								expectedFloats := []float64{
									1.1,
									2.2,
								}
								if len(actual.SetOfFloats) != len(expectedFloats) {
									return fmt.Errorf("expected %d floats but got %d", len(expectedFloats), len(actual.SetOfFloats))
								}
								for _, v := range expectedFloats {
									exists := false
									for _, a := range actual.SetOfFloats {
										if v == a {
											exists = true
											break
										}
									}
									if !exists {
										return fmt.Errorf("expected the float %f to exist but it didn't", v)
									}
								}

								return nil
							},
							Delete: func(_ *schema.ResourceData, _ interface{}) error {
								return nil
							},
						},
					},
				}, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `resource "validator_decoder" "test" {}`,
			},
		},
	})
}

func TestAccPluginSDKAndEncoder(t *testing.T) {
	os.Setenv("TF_ACC", "1")

	type NestedType struct {
		Key string `tfschema:"key"`
	}
	type MyType struct {
		Hello         string             `tfschema:"hello"`
		RandomNumber  int                `tfschema:"random_number"`
		Enabled       bool               `tfschema:"enabled"`
		ListOfStrings []string           `tfschema:"list_of_strings"`
		ListOfNumbers []int              `tfschema:"list_of_numbers"`
		ListOfBools   []bool             `tfschema:"list_of_bools"`
		ListOfFloats  []float64          `tfschema:"list_of_floats"`
		NestedObject  []NestedType       `tfschema:"nested_object"`
		MapOfStrings  map[string]string  `tfschema:"map_of_strings"`
		MapOfNumbers  map[string]int     `tfschema:"map_of_numbers"`
		MapOfBools    map[string]bool    `tfschema:"map_of_bools"`
		MapOfFloats   map[string]float64 `tfschema:"map_of_floats"`
		SetOfStrings  []string           `tfschema:"set_of_strings"`
		SetOfNumbers  []int              `tfschema:"set_of_numbers"`
		SetOfBools    []bool             `tfschema:"set_of_bools"`
		SetOfFloats   []float64          `tfschema:"set_of_floats"`
	}

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]terraform.ResourceProviderFactory{
			"validator": func() (terraform.ResourceProvider, error) {
				return &schema.Provider{
					DataSourcesMap: map[string]*schema.Resource{},
					ResourcesMap: map[string]*schema.Resource{
						"validator_encoder": {
							Schema: map[string]*schema.Schema{
								"hello": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"random_number": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"enabled": {
									Type:     schema.TypeBool,
									Computed: true,
								},
								"list_of_strings": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"list_of_numbers": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeInt,
									},
								},
								"list_of_bools": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeBool,
									},
								},
								"list_of_floats": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeFloat,
									},
								},
								"nested_object": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"key": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"map_of_strings": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"map_of_numbers": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeInt,
									},
								},
								"map_of_bools": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeBool,
									},
								},
								"map_of_floats": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeFloat,
									},
								},
								"set_of_strings": {
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"set_of_numbers": {
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeInt,
									},
								},
								"set_of_bools": {
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeBool,
									},
								},
								"set_of_floats": {
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeFloat,
									},
								},
							},
							Create: func(d *schema.ResourceData, i interface{}) error {
								wrapper := ResourceMetaData{
									ResourceData:             d,
									Logger:                   ConsoleLogger{},
									serializationDebugLogger: ConsoleLogger{},
								}

								input := MyType{
									Hello:         "world",
									RandomNumber:  42,
									Enabled:       true,
									ListOfStrings: []string{"hello", "there"},
									ListOfNumbers: []int{1, 2, 4},
									ListOfBools:   []bool{true, false},
									ListOfFloats:  []float64{-1.234567894321, 2.3456789},
									NestedObject: []NestedType{
										{
											Key: "value",
										},
									},
									MapOfStrings: map[string]string{
										"bingo": "bango",
									},
									MapOfNumbers: map[string]int{
										"lucky": 21,
									},
									MapOfBools: map[string]bool{
										"friday": true,
									},
									MapOfFloats: map[string]float64{
										"pi": 3.14159,
									},
								}

								d.SetId("some-id")
								if err := wrapper.Encode(&input); err != nil {
									return fmt.Errorf("encoding: %+v", err)
								}
								return nil
							},
							Read: func(d *schema.ResourceData, _ interface{}) error {
								return nil
							},
							Delete: func(_ *schema.ResourceData, _ interface{}) error {
								return nil
							},
						},
					},
				}, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `resource "validator_encoder" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceStateMatches("validator_encoder.test", map[string]interface{}{
						"id":                   "some-id",
						"hello":                "world",
						"random_number":        "42",
						"enabled":              "true",
						"list_of_strings.#":    "2",
						"list_of_strings.0":    "hello",
						"list_of_strings.1":    "there",
						"list_of_numbers.#":    "3",
						"list_of_numbers.0":    "1",
						"list_of_numbers.1":    "2",
						"list_of_numbers.2":    "4",
						"list_of_bools.#":      "2",
						"list_of_bools.0":      "true",
						"list_of_bools.1":      "false",
						"list_of_floats.#":     "2",
						"list_of_floats.0":     "-1.234567894321",
						"list_of_floats.1":     "2.3456789",
						"nested_object.#":      "1",
						"nested_object.0.key":  "value",
						"map_of_strings.%":     "1",
						"map_of_strings.bingo": "bango",
						"map_of_numbers.%":     "1",
						"map_of_numbers.lucky": "21",
						"map_of_bools.%":       "1",
						"map_of_bools.friday":  "true",
						"map_of_floats.%":      "1",
						"map_of_floats.pi":     "3.14159",
						"set_of_bools.#":       "0",
						"set_of_floats.#":      "0",
						"set_of_numbers.#":     "0",
						"set_of_strings.#":     "0",
					}),
				),
			},
		},
	})
}

func TestAccPluginSDKReturnsComputedFields(t *testing.T) {
	os.Setenv("TF_ACC", "1")

	resourceName := "validator_computed.test"
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]terraform.ResourceProviderFactory{
			"validator": func() (terraform.ResourceProvider, error) {
				return &schema.Provider{
					DataSourcesMap: map[string]*schema.Resource{},
					ResourcesMap: map[string]*schema.Resource{
						"validator_computed": computedFieldsResource(),
					},
				}, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `resource "validator_computed" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceStateMatches(resourceName, map[string]interface{}{
						"id":                  "does-not-matter",
						"hello":               "world",
						"random_number":       "42",
						"enabled":             "true",
						"list_of_strings.#":   "2",
						"list_of_strings.0":   "hello",
						"list_of_strings.1":   "there",
						"list_of_numbers.#":   "3",
						"list_of_numbers.0":   "1",
						"list_of_numbers.1":   "2",
						"list_of_numbers.2":   "4",
						"list_of_bools.#":     "2",
						"list_of_bools.0":     "true",
						"list_of_bools.1":     "false",
						"list_of_floats.#":    "2",
						"list_of_floats.0":    "-1.234567894321",
						"list_of_floats.1":    "2.3456789",
						"nested_object.#":     "1",
						"nested_object.0.key": "value",
						// Sets can't really be computed, so this isn't that big a deal
					}),
				),
			},
		},
	})
}

func computedFieldsResource() *schema.Resource {
	var readFunc = func(d *schema.ResourceData, _ interface{}) error {
		d.Set("hello", "world")
		d.Set("random_number", 42)
		d.Set("enabled", true)
		d.Set("list_of_strings", []string{"hello", "there"})
		d.Set("list_of_numbers", []int{1, 2, 4})
		d.Set("list_of_bools", []bool{true, false})
		d.Set("list_of_floats", []float64{-1.234567894321, 2.3456789})
		d.Set("nested_object", []interface{}{
			map[string]interface{}{
				"key": "value",
			},
		})
		return nil
	}
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hello": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"random_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"list_of_strings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"list_of_numbers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"list_of_bools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
			},
			"list_of_floats": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeFloat,
				},
			},
			"nested_object": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Create: func(d *schema.ResourceData, meta interface{}) error {
			d.SetId("does-not-matter")
			return readFunc(d, meta)
		},
		Read: readFunc,
		Delete: func(_ *schema.ResourceData, _ interface{}) error {
			return nil
		},
	}
}

func testCheckResourceStateMatches(resourceName string, values map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resources, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %q was not found in the state", resourceName)
		}

		state := resources.Primary
		if len(state.Attributes) != len(values) {
			return fmt.Errorf("expected %d values but got %d.\n\nExpected: %+v\n\nActual: %+v", len(values), len(state.Attributes), values, state.Attributes)
		}

		for key, expectedValue := range values {
			actualValue, exists := state.Attributes[key]
			if !exists {
				return fmt.Errorf("key %q was not found", key)
			}

			if !reflect.DeepEqual(expectedValue, actualValue) {
				return fmt.Errorf("values didn't match for %q.\n\nExpected: %+v\n\nActual: %+v", key, expectedValue, actualValue)
			}
		}

		return nil
	}
}
