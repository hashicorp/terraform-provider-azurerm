// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package edgezones

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
)

func TestFlatten(t *testing.T) {
	testData := []struct {
		Name     string
		Input    *edgezones.Model
		Expected string
	}{
		{"Usual case", &edgezones.Model{Name: "East US"}, "eastus"},
		{"All upcase", &edgezones.Model{Name: "EASTUS"}, "eastus"},
		{"Camel case", &edgezones.Model{Name: "EastUS"}, "eastus"},
		{"Lower case with space", &edgezones.Model{Name: "east us"}, "eastus"},
		{"Upper case with space", &edgezones.Model{Name: "EAST US"}, "eastus"},
		{"West A b C 2", &edgezones.Model{Name: "West A b C 2"}, "westabc2"},
		{"Empty string", &edgezones.Model{Name: ""}, ""},
		{"Whitespace only", &edgezones.Model{Name: "  "}, ""},
		{"Nil input", nil, ""},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Test %q", v.Name)

		actual := Flatten(v.Input)
		if !reflect.DeepEqual(actual, v.Expected) {
			t.Fatalf("Expected %+v but got %+v", v.Expected, actual)
		}
	}
}
