// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func TestFromTypedObject(t *testing.T) {
	testData := map[string]string{
		"key1": "value1",
		"key2": "21",
		"key3": "value3",
	}

	expanded := FromTypedObject(testData)

	if len(expanded) != 3 {
		t.Fatalf("Expected 3 results in expanded tag map, got %d", len(expanded))
	}

	for k, v := range testData {
		if *expanded[k] != v {
			t.Fatalf("Expanded value %q incorrect: expected %q, got %q", k, v, *expanded[k])
		}
	}
}

func TestToTypedObject(t *testing.T) {
	testData := []struct {
		Name     string
		Input    map[string]*string
		Expected map[string]string
	}{
		{
			Name:     "Empty",
			Input:    map[string]*string{},
			Expected: map[string]string{},
		},
		{
			Name: "One Item",
			Input: map[string]*string{
				"hello": pointer.To("there"),
			},
			Expected: map[string]string{
				"hello": "there",
			},
		},
		{
			Name: "Multiple Items",
			Input: map[string]*string{
				"euros": pointer.To("3"),
				"hello": pointer.To("there"),
				"panda": pointer.To("pops"),
			},
			Expected: map[string]string{
				"euros": "3",
				"hello": "there",
				"panda": "pops",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Test %q", v.Name)

		actual := ToTypedObject(v.Input)
		if !reflect.DeepEqual(actual, v.Expected) {
			t.Fatalf("Expected %+v but got %+v", actual, v.Expected)
		}
	}
}
