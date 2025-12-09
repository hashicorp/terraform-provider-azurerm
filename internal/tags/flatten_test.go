// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func TestFlatten(t *testing.T) {
	testData := []struct {
		Name     string
		Input    map[string]*string
		Expected map[string]interface{}
	}{
		{
			Name:     "Empty",
			Input:    map[string]*string{},
			Expected: map[string]interface{}{},
		},
		{
			Name: "One Item",
			Input: map[string]*string{
				"hello": pointer.To("there"),
			},
			Expected: map[string]interface{}{
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
			Expected: map[string]interface{}{
				"euros": "3",
				"hello": "there",
				"panda": "pops",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Test %q", v.Name)

		actual := Flatten(v.Input)
		if !reflect.DeepEqual(actual, v.Expected) {
			t.Fatalf("Expected %+v but got %+v", actual, v.Expected)
		}
	}
}
