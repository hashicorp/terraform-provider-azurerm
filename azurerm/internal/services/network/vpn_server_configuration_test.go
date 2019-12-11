package network

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func TestParseVpnServerConfiguration(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *VpnServerConfigurationResourceID
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No VPN Server Configurations Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Expected: nil,
		},
		{
			Name:     "No VPN Server Configurations Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/vpnServerConfigurations/",
			Expected: nil,
		},
		{
			Name:  "Completed",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/vpnServerConfigurations/example",
			Expected: &VpnServerConfigurationResourceID{
				Name: "example",
				Base: azure.ResourceID{
					ResourceGroup: "foo",
				},
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ParseVpnServerConfigurationID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.Base.ResourceGroup != v.Expected.Base.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.Base.ResourceGroup, actual.Base.ResourceGroup)
		}
	}
}

func TestValidateVpnServerConfiguration(t *testing.T) {
	testData := []struct {
		Name  string
		Input string
		Valid bool
	}{
		{
			Name:  "Empty",
			Input: "",
			Valid: false,
		},
		{
			Name:  "No VPN Server Configurations Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Valid: false,
		},
		{
			Name:  "No VPN Server Configurations Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/vpnServerConfigurations/",
			Valid: false,
		},
		{
			Name:  "Completed",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/vpnServerConfigurations/example",
			Valid: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		_, errors := ValidateVpnServerConfigurationID(v.Input, "vpn_server_configuration_id")
		isValid := len(errors) == 0
		if v.Valid != isValid {
			t.Fatalf("Expected %t but got %t", v.Valid, isValid)
		}
	}
}
