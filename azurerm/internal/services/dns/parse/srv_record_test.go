package parse

import (
	"testing"
)

func TestDnsSrvRecordId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *SrvRecordId
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
			Name:     "Missing SRV Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1/SRV/",
			Expected: nil,
		},
		{
			Name:  "DNS SRV Record ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1/SRV/myrecord1",
			Expected: &SrvRecordId{
				ResourceGroup: "resGroup1",
				DnszoneName:   "zone1",
				SRVName:       "myrecord1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1/srv/myrecord1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := SrvRecordID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.SRVName != v.Expected.SRVName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.SRVName, actual.SRVName)
		}
		if actual.DnszoneName != v.Expected.DnszoneName {
			t.Fatalf("Expected %q but got %q for ZoneName", v.Expected.DnszoneName, actual.DnszoneName)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
