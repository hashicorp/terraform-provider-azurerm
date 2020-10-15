package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = VpnSiteId{}

func TestVpnSiteIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewVpnSiteID("group1", "site1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/vpnSites/site1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestVpnSiteID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *VpnSiteId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Resource Groups Segment",
			Input: "/subscriptions/11111111-1111-1111-1111-111111111111",
			Error: true,
		},
		{
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/",
			Error: true,
		},
		{
			Name:  "Missing leading slash",
			Input: "subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/group1",
			Error: true,
		},
		{
			Name:  "Malformed segments",
			Input: "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/group1/foo/bar",
			Error: true,
		},
		{
			Name:  "No vpn site segment",
			Input: "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/group1/providers/Microsoft.Network",
			Error: true,
		},
		{
			Name:  "Correct",
			Input: "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/group1/providers/Microsoft.Network/vpnSites/site1",
			Expect: &VpnSiteId{
				ResourceGroup: "group1",
				Name:          "site1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := VpnSiteID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get")
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
	}
}
