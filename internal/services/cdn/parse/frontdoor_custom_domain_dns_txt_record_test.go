package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = FrontdoorCustomDomainDnsTxtRecordId{}

func TestFrontdoorCustomDomainDnsTxtRecordIDFormatter(t *testing.T) {
	actual := NewFrontdoorCustomDomainDnsTxtRecordID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "profile1", "customDomain1", "txt1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/customDomain1/dnsTxt/txt1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestFrontdoorCustomDomainDnsTxtRecordID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FrontdoorCustomDomainDnsTxtRecordId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing ProfileName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/",
			Error: true,
		},

		{
			// missing value for ProfileName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/",
			Error: true,
		},

		{
			// missing CustomDomainName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/",
			Error: true,
		},

		{
			// missing value for CustomDomainName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/",
			Error: true,
		},

		{
			// missing DnsTxtName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/customDomain1/",
			Error: true,
		},

		{
			// missing value for DnsTxtName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/customDomain1/dnsTxt/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/customDomain1/dnsTxt/txt1",
			Expected: &FrontdoorCustomDomainDnsTxtRecordId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "resourceGroup1",
				ProfileName:      "profile1",
				CustomDomainName: "customDomain1",
				DnsTxtName:       "txt1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.CDN/PROFILES/PROFILE1/CUSTOMDOMAINS/CUSTOMDOMAIN1/DNSTXT/TXT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := FrontdoorCustomDomainDnsTxtRecordID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.ProfileName != v.Expected.ProfileName {
			t.Fatalf("Expected %q but got %q for ProfileName", v.Expected.ProfileName, actual.ProfileName)
		}
		if actual.CustomDomainName != v.Expected.CustomDomainName {
			t.Fatalf("Expected %q but got %q for CustomDomainName", v.Expected.CustomDomainName, actual.CustomDomainName)
		}
		if actual.DnsTxtName != v.Expected.DnsTxtName {
			t.Fatalf("Expected %q but got %q for DnsTxtName", v.Expected.DnsTxtName, actual.DnsTxtName)
		}
	}
}

func TestFrontdoorCustomDomainDnsTxtRecordIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FrontdoorCustomDomainDnsTxtRecordId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing ProfileName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/",
			Error: true,
		},

		{
			// missing value for ProfileName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/",
			Error: true,
		},

		{
			// missing CustomDomainName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/",
			Error: true,
		},

		{
			// missing value for CustomDomainName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/",
			Error: true,
		},

		{
			// missing DnsTxtName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/customDomain1/",
			Error: true,
		},

		{
			// missing value for DnsTxtName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/customDomain1/dnsTxt/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/customDomain1/dnsTxt/txt1",
			Expected: &FrontdoorCustomDomainDnsTxtRecordId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "resourceGroup1",
				ProfileName:      "profile1",
				CustomDomainName: "customDomain1",
				DnsTxtName:       "txt1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customdomains/customDomain1/dnstxt/txt1",
			Expected: &FrontdoorCustomDomainDnsTxtRecordId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "resourceGroup1",
				ProfileName:      "profile1",
				CustomDomainName: "customDomain1",
				DnsTxtName:       "txt1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/PROFILES/profile1/CUSTOMDOMAINS/customDomain1/DNSTXT/txt1",
			Expected: &FrontdoorCustomDomainDnsTxtRecordId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "resourceGroup1",
				ProfileName:      "profile1",
				CustomDomainName: "customDomain1",
				DnsTxtName:       "txt1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/PrOfIlEs/profile1/CuStOmDoMaInS/customDomain1/DnStXt/txt1",
			Expected: &FrontdoorCustomDomainDnsTxtRecordId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "resourceGroup1",
				ProfileName:      "profile1",
				CustomDomainName: "customDomain1",
				DnsTxtName:       "txt1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := FrontdoorCustomDomainDnsTxtRecordIDInsensitively(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.ProfileName != v.Expected.ProfileName {
			t.Fatalf("Expected %q but got %q for ProfileName", v.Expected.ProfileName, actual.ProfileName)
		}
		if actual.CustomDomainName != v.Expected.CustomDomainName {
			t.Fatalf("Expected %q but got %q for CustomDomainName", v.Expected.CustomDomainName, actual.CustomDomainName)
		}
		if actual.DnsTxtName != v.Expected.DnsTxtName {
			t.Fatalf("Expected %q but got %q for DnsTxtName", v.Expected.DnsTxtName, actual.DnsTxtName)
		}
	}
}
