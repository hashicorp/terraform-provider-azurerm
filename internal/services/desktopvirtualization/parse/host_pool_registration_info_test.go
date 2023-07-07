// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = HostPoolRegistrationInfoId{}

func TestHostPoolRegistrationInfoIDFormatter(t *testing.T) {
	actual := NewHostPoolRegistrationInfoID("12345678-1234-9876-4563-123456789012", "resGroup1", "pool1", "default").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/hostPools/pool1/registrationInfo/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestHostPoolRegistrationInfoID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *HostPoolRegistrationInfoId
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
			// missing HostPoolName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/",
			Error: true,
		},

		{
			// missing value for HostPoolName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/hostPools/",
			Error: true,
		},

		{
			// missing RegistrationInfoName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/hostPools/pool1/",
			Error: true,
		},

		{
			// missing value for RegistrationInfoName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/hostPools/pool1/registrationInfo/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/hostPools/pool1/registrationInfo/default",
			Expected: &HostPoolRegistrationInfoId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resGroup1",
				HostPoolName:         "pool1",
				RegistrationInfoName: "default",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.DESKTOPVIRTUALIZATION/HOSTPOOLS/POOL1/REGISTRATIONINFO/DEFAULT",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := HostPoolRegistrationInfoID(v.Input)
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
		if actual.HostPoolName != v.Expected.HostPoolName {
			t.Fatalf("Expected %q but got %q for HostPoolName", v.Expected.HostPoolName, actual.HostPoolName)
		}
		if actual.RegistrationInfoName != v.Expected.RegistrationInfoName {
			t.Fatalf("Expected %q but got %q for RegistrationInfoName", v.Expected.RegistrationInfoName, actual.RegistrationInfoName)
		}
	}
}

func TestHostPoolRegistrationInfoIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *HostPoolRegistrationInfoId
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
			// missing HostPoolName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/",
			Error: true,
		},

		{
			// missing value for HostPoolName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/hostPools/",
			Error: true,
		},

		{
			// missing RegistrationInfoName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/hostPools/pool1/",
			Error: true,
		},

		{
			// missing value for RegistrationInfoName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/hostPools/pool1/registrationInfo/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/hostPools/pool1/registrationInfo/default",
			Expected: &HostPoolRegistrationInfoId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resGroup1",
				HostPoolName:         "pool1",
				RegistrationInfoName: "default",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/hostpools/pool1/registrationinfo/default",
			Expected: &HostPoolRegistrationInfoId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resGroup1",
				HostPoolName:         "pool1",
				RegistrationInfoName: "default",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/HOSTPOOLS/pool1/REGISTRATIONINFO/default",
			Expected: &HostPoolRegistrationInfoId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resGroup1",
				HostPoolName:         "pool1",
				RegistrationInfoName: "default",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/HoStPoOlS/pool1/ReGiStRaTiOnInFo/default",
			Expected: &HostPoolRegistrationInfoId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resGroup1",
				HostPoolName:         "pool1",
				RegistrationInfoName: "default",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := HostPoolRegistrationInfoIDInsensitively(v.Input)
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
		if actual.HostPoolName != v.Expected.HostPoolName {
			t.Fatalf("Expected %q but got %q for HostPoolName", v.Expected.HostPoolName, actual.HostPoolName)
		}
		if actual.RegistrationInfoName != v.Expected.RegistrationInfoName {
			t.Fatalf("Expected %q but got %q for RegistrationInfoName", v.Expected.RegistrationInfoName, actual.RegistrationInfoName)
		}
	}
}
