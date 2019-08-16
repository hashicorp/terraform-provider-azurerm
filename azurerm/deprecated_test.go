package azurerm

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidateMaximumNumberOfARMTags(t *testing.T) {
	tagsMap := make(map[string]interface{})
	for i := 0; i < 51; i++ {
		tagsMap[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
	}

	_, es := validateAzureRMTags(tagsMap, "tags")

	if len(es) != 1 {
		t.Fatal("Expected one validation error for too many tags")
	}

	if !strings.Contains(es[0].Error(), "a maximum of 50 tags") {
		t.Fatal("Wrong validation error message for too many tags")
	}
}

func TestValidateARMTagMaxKeyLength(t *testing.T) {
	tooLongKey := strings.Repeat("long", 128) + "a"
	tagsMap := make(map[string]interface{})
	tagsMap[tooLongKey] = "value"

	_, es := validateAzureRMTags(tagsMap, "tags")
	if len(es) != 1 {
		t.Fatal("Expected one validation error for a key which is > 512 chars")
	}

	if !strings.Contains(es[0].Error(), "maximum length for a tag key") {
		t.Fatal("Wrong validation error message maximum tag key length")
	}

	if !strings.Contains(es[0].Error(), tooLongKey) {
		t.Fatal("Expected validated error to contain the key name")
	}

	if !strings.Contains(es[0].Error(), "513") {
		t.Fatal("Expected the length in the validation error for tag key")
	}
}

func TestValidateARMTagMaxValueLength(t *testing.T) {
	tagsMap := make(map[string]interface{})
	tagsMap["toolong"] = strings.Repeat("long", 64) + "a"

	_, es := validateAzureRMTags(tagsMap, "tags")
	if len(es) != 1 {
		t.Fatal("Expected one validation error for a value which is > 256 chars")
	}

	if !strings.Contains(es[0].Error(), "maximum length for a tag value") {
		t.Fatal("Wrong validation error message for maximum tag value length")
	}

	if !strings.Contains(es[0].Error(), "toolong") {
		t.Fatal("Expected validated error to contain the key name")
	}

	if !strings.Contains(es[0].Error(), "257") {
		t.Fatal("Expected the length in the validation error for value")
	}
}

func TestExpandARMTags(t *testing.T) {
	testData := make(map[string]interface{})
	testData["key1"] = "value1"
	testData["key2"] = 21
	testData["key3"] = "value3"

	expanded := expandTags(testData)

	if len(expanded) != 3 {
		t.Fatalf("Expected 3 results in expanded tag map, got %d", len(expanded))
	}

	for k, v := range testData {
		var strVal string
		switch v := v.(type) {
		case string:
			strVal = v
		case int:
			strVal = fmt.Sprintf("%d", v)
		}

		if *expanded[k] != strVal {
			t.Fatalf("Expanded value %q incorrect: expected %q, got %q", k, strVal, *expanded[k])
		}
	}
}

func TestFilterARMTags(t *testing.T) {
	testData := make(map[string]*string)
	valueData := [3]string{"value1", "value2", "value3"}

	testData["key1"] = &valueData[0]
	testData["key2"] = &valueData[1]
	testData["key3"] = &valueData[2]

	filtered := filterTags(testData, "key1", "key3", "")

	if len(filtered) != 1 {
		t.Fatalf("Expected 1 result in filtered tag map, got %d", len(filtered))
	}

	if filtered["key2"] != &valueData[1] {
		t.Fatalf("Expected %v in filtered tag map, got %v", valueData[1], *filtered["key2"])
	}
}

func TestValidateRFC3339Date(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "",
			ErrCount: 1,
		},
		{
			Value:    "Random",
			ErrCount: 1,
		},
		{
			Value:    "2017-01-01",
			ErrCount: 1,
		},
		{
			Value:    "2017-01-01T01:23:45",
			ErrCount: 1,
		},
		{
			Value:    "2017-01-01T01:23:45+00:00",
			ErrCount: 0,
		},
		{
			Value:    "2017-01-01T01:23:45Z",
			ErrCount: 0,
		},
	}

	for _, tc := range cases {
		_, errors := validateRFC3339Date(tc.Value, "example")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected validateRFC3339Date to trigger '%d' errors for '%s' - got '%d'", tc.ErrCount, tc.Value, len(errors))
		}
	}
}

func TestValidateIso8601Duration(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			// Date components only
			Value:  "P1Y2M3D",
			Errors: 0,
		},
		{
			// Time components only
			Value:  "PT7H42M3S",
			Errors: 0,
		},
		{
			// Date and time components
			Value:  "P1Y2M3DT7H42M3S",
			Errors: 0,
		},
		{
			// Invalid prefix
			Value:  "1Y2M3DT7H42M3S",
			Errors: 1,
		},
		{
			// Wrong order of components, i.e. invalid format
			Value:  "PT7H42M3S1Y2M3D",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateIso8601Duration()(tc.Value, "example")

		if len(errors) != tc.Errors {
			t.Fatalf("Expected validateIso8601Duration to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}

func TestValidateCollation(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "en-US",
			Errors: 1,
		},
		{
			Value:  "en_US",
			Errors: 0,
		},
		{
			Value:  "en US",
			Errors: 0,
		},
		{
			Value:  "English_United States.1252",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := validateCollation()(tc.Value, "collation")
		if len(errors) != tc.Errors {
			t.Fatalf("Expected validateCollation to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}

func TestValidateAzureVirtualMachineTimeZone(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "",
			Errors: 0,
		},
		{
			Value:  "UTC",
			Errors: 0,
		},
		{
			Value:  "China Standard Time",
			Errors: 0,
		},
		{
			// Valid UTC time zone
			Value:  "utc-11",
			Errors: 0,
		},
		{
			// Invalid UTC time zone
			Value:  "UTC-30",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateAzureVirtualMachineTimeZone()(tc.Value, "unittest")

		if len(errors) != tc.Errors {
			t.Fatalf("Expected validateAzureVMTimeZone to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}
