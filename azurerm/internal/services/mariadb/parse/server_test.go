package parse

import (
	"testing"
)

func TestValidateMariaDbServerServerID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *ServerId
	}{
		{
			Name:     "Empty resource ID",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No resourceGroups segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No resource group name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource group",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/",
			Expected: nil,
		},
		{
			Name:     "Missing server name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DBforMariaDB/servers/",
			Expected: nil,
		},
		{
			Name:  "Valid",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DBforMariaDB/servers/test-mariadb",
			Expected: &ServerId{
				Name:          "test-mariadb",
				ResourceGroup: "test-rg",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ServerID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
