package parse

import (
	"testing"
)

func TestCdnEndpointCustomDomainID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *CdnEndpointCustomDomainId
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
			Name:  "Missing leading slash",
			Input: "subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1",
			Error: true,
		},
		{
			Name:  "Malformed segments",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/foo/bar",
			Error: true,
		},
		{
			Name:  "No Profile Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cdn/profiles",
			Error: true,
		},
		{
			Name:  "No Endpoint Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cdn/profiles/profile1/endpoints",
			Error: true,
		},
		{
			Name:  "No Custom Domain Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cdn/profiles/profile1/endpoints/endpoint1/customdomains",
			Error: true,
		},
		{
			Name:  "Malformed Custom Domain Name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cdn/profiles/profile1/endpoints/endpoint1/customdomains/domain1/foo",
			Error: true,
		},
		{
			Name:  "Correct Case",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cdn/profiles/profile1/endpoints/endpoint1/customdomains/domain1",
			Expect: &CdnEndpointCustomDomainId{
				ResourceGroup: "group1",
				ProfileName:   "profile1",
				EndpointName:  "endpoint1",
				Name:          "domain1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := CdnEndpointCustomDomainID(v.Input)
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
