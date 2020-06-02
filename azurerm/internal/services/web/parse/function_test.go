package parse

import (
	"testing"
)

func TestFunctionId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *FunctionID
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Function",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Web/sites/site1",
			Expected: nil,
		},
		{
			Name:  "App Service Function ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Web/sites/site1/functions/function1",
			Expected: &FunctionID{
				Name:            "function1",
				FunctionAppName: "site1",
				ResourceGroup:   "group1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ParseFunctionID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.FunctionAppName != v.Expected.FunctionAppName {
			t.Fatalf("Expected %q but got %q for Function App", v.Expected.FunctionAppName, actual.FunctionAppName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
