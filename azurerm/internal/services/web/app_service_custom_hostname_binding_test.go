package web

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
)

func TestParseAppServiceCustomHostnameBinding(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *parse.AppServiceCustomHostnameBindingId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Missing Sites Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/",
			Expected: nil,
		},
		{
			Name:     "App Service Resource ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1",
			Expected: nil,
		},
		{
			Name:     "Missing Host Name Bindings Valud",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1/hostNameBindings/",
			Expected: nil,
		},
		{
			Name:  "Valid Resource ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1/hostNameBindings/binding1",
			Expected: &parse.AppServiceCustomHostnameBindingId{
				Name:           "binding1",
				AppServiceName: "site1",
				ResourceGroup:  "mygroup1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/Sites/site1/HostNameBindings/binding1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := parse.AppServiceCustomHostnameBindingID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.AppServiceName != v.Expected.AppServiceName {
			t.Fatalf("Expected %q but got %q for AppServiceName", v.Expected.AppServiceName, actual.AppServiceName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
