package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = DedicatedHostId{}

func TestDedicatedHostIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	hostGroupId := NewDedicatedHostGroupId("group1", "hostGroup1")
	actual := NewDedicatedHostId(hostGroupId, "host1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Compute/hostGroups/hostGroup1/hosts/host1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestDedicatedHostID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *DedicatedHostId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Resource Groups Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},
		{
			Name:  "Resource Group ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Error: true,
		},
		{
			Name:  "Missing Host Group Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/",
			Error: true,
		},
		{
			Name:  "Host Group ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/group1/",
			Error: true,
		},
		{
			Name:  "Missing Host Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/group1/hosts/",
			Error: true,
		},
		{
			Name:  "Host ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/group1/hosts/host1",
			Error: false,
			Expect: &DedicatedHostId{
				ResourceGroup: "resGroup1",
				HostGroup:     "group1",
				Name:          "host1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/group1/Hosts/host1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := DedicatedHostID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}

		if actual.HostGroup != v.Expect.HostGroup {
			t.Fatalf("Expected %q but got %q for HostGroup", v.Expect.HostGroup, actual.HostGroup)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}
	}
}
