package parse

import (
	"testing"
)

func TestDnsPtrRecordId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *PtrRecordId
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
			Name:     "Missing DNS Zones Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/",
			Expected: nil,
		},
		{
			Name:     "DNS Zone ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1",
			Expected: nil,
		},
		{
			Name:     "Missing PTR Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1/PTR/",
			Expected: nil,
		},
		{
			Name:  "DNS PTR Record ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1/PTR/myrecord1",
			Expected: &PtrRecordId{
				ResourceGroup: "resGroup1",
				DnszoneName:   "zone1",
				PTRName:       "myrecord1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1/ptr/myrecord1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := PtrRecordID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.PTRName != v.Expected.PTRName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.PTRName, actual.PTRName)
		}
		if actual.DnszoneName != v.Expected.DnszoneName {
			t.Fatalf("Expected %q but got %q for ZoneName", v.Expected.DnszoneName, actual.DnszoneName)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
