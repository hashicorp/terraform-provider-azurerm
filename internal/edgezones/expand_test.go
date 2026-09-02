// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package edgezones

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
)

func TestExpandEdgeZone(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *edgezones.Model
	}{
		{"Usual case", "East US", &edgezones.Model{Name: "eastus"}},
		{"All upcase", "EASTUS", &edgezones.Model{Name: "eastus"}},
		{"Camel case", "EastUS", &edgezones.Model{Name: "eastus"}},
		{"Lower case with space", "east us", &edgezones.Model{Name: "eastus"}},
		{"Upper case with space", "EAST US", &edgezones.Model{Name: "eastus"}},
		{"West A b C 2", "West A b C 2", &edgezones.Model{Name: "westabc2"}},
		{"Empty string", "", nil},
		{"Whitespace only", "  ", nil},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Test %q", v.Name)

		actual := ExpandEdgeZone(v.Input)
		if !reflect.DeepEqual(actual, v.Expected) {
			t.Fatalf("Expected %+v but got %+v", v.Expected, actual)
		}
	}
}
