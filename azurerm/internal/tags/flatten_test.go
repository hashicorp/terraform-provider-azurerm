package tags

import (
	"reflect"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
				"hello": utils.String("there"),
			},
			Expected: map[string]interface{}{
				"hello": "there",
			},
		},
		{
			Name: "Multiple Items",
			Input: map[string]*string{
				"euros": utils.String("3"),
				"hello": utils.String("there"),
				"panda": utils.String("pops"),
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
