package parse

import "testing"

func TestMySQLServerKeyID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *KeyId
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
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DBforMySQL/servers/",
			Expected: nil,
		},
		{
			Name:     "MySQL Server ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DBforMySQL/servers/test-mysql/",
			Expected: nil,
		},
		{
			Name:     "Missing key name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DBforMySQL/servers/test-mysql/keys/",
			Expected: nil,
		},
		{
			Name:  "Valid",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DBforMySQL/servers/test-mysql/keys/key1",
			Expected: &KeyId{
				Name:          "key1",
				ResourceGroup: "test-rg",
				ServerName:    "test-mysql",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := KeyID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ServerName != v.Expected.ServerName {
			t.Fatalf("Expected %q but got %q for ServerName", v.Expected.ServerName, actual.ServerName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
