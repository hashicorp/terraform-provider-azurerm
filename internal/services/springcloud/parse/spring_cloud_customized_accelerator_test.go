// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SpringCloudCustomizedAcceleratorId{}

func TestSpringCloudCustomizedAcceleratorIDFormatter(t *testing.T) {
	actual := NewSpringCloudCustomizedAcceleratorID("12345678-1234-9876-4563-123456789012", "resGroup1", "spring1", "default", "customizedAccelerator1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/applicationAccelerators/default/customizedAccelerators/customizedAccelerator1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSpringCloudCustomizedAcceleratorID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SpringCloudCustomizedAcceleratorId
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
			// missing SpringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/",
			Error: true,
		},

		{
			// missing value for SpringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/",
			Error: true,
		},

		{
			// missing ApplicationAcceleratorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/",
			Error: true,
		},

		{
			// missing value for ApplicationAcceleratorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/applicationAccelerators/",
			Error: true,
		},

		{
			// missing CustomizedAcceleratorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/applicationAccelerators/default/",
			Error: true,
		},

		{
			// missing value for CustomizedAcceleratorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/applicationAccelerators/default/customizedAccelerators/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/applicationAccelerators/default/customizedAccelerators/customizedAccelerator1",
			Expected: &SpringCloudCustomizedAcceleratorId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				SpringName:                 "spring1",
				ApplicationAcceleratorName: "default",
				CustomizedAcceleratorName:  "customizedAccelerator1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.APPPLATFORM/SPRING/SPRING1/APPLICATIONACCELERATORS/DEFAULT/CUSTOMIZEDACCELERATORS/CUSTOMIZEDACCELERATOR1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SpringCloudCustomizedAcceleratorID(v.Input)
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
		if actual.SpringName != v.Expected.SpringName {
			t.Fatalf("Expected %q but got %q for SpringName", v.Expected.SpringName, actual.SpringName)
		}
		if actual.ApplicationAcceleratorName != v.Expected.ApplicationAcceleratorName {
			t.Fatalf("Expected %q but got %q for ApplicationAcceleratorName", v.Expected.ApplicationAcceleratorName, actual.ApplicationAcceleratorName)
		}
		if actual.CustomizedAcceleratorName != v.Expected.CustomizedAcceleratorName {
			t.Fatalf("Expected %q but got %q for CustomizedAcceleratorName", v.Expected.CustomizedAcceleratorName, actual.CustomizedAcceleratorName)
		}
	}
}

func TestSpringCloudCustomizedAcceleratorIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SpringCloudCustomizedAcceleratorId
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
			// missing SpringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/",
			Error: true,
		},

		{
			// missing value for SpringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/",
			Error: true,
		},

		{
			// missing ApplicationAcceleratorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/",
			Error: true,
		},

		{
			// missing value for ApplicationAcceleratorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/applicationAccelerators/",
			Error: true,
		},

		{
			// missing CustomizedAcceleratorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/applicationAccelerators/default/",
			Error: true,
		},

		{
			// missing value for CustomizedAcceleratorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/applicationAccelerators/default/customizedAccelerators/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/applicationAccelerators/default/customizedAccelerators/customizedAccelerator1",
			Expected: &SpringCloudCustomizedAcceleratorId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				SpringName:                 "spring1",
				ApplicationAcceleratorName: "default",
				CustomizedAcceleratorName:  "customizedAccelerator1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/applicationaccelerators/default/customizedaccelerators/customizedAccelerator1",
			Expected: &SpringCloudCustomizedAcceleratorId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				SpringName:                 "spring1",
				ApplicationAcceleratorName: "default",
				CustomizedAcceleratorName:  "customizedAccelerator1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/SPRING/spring1/APPLICATIONACCELERATORS/default/CUSTOMIZEDACCELERATORS/customizedAccelerator1",
			Expected: &SpringCloudCustomizedAcceleratorId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				SpringName:                 "spring1",
				ApplicationAcceleratorName: "default",
				CustomizedAcceleratorName:  "customizedAccelerator1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/SpRiNg/spring1/ApPlIcAtIoNaCcElErAtOrS/default/CuStOmIzEdAcCeLeRaToRs/customizedAccelerator1",
			Expected: &SpringCloudCustomizedAcceleratorId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				SpringName:                 "spring1",
				ApplicationAcceleratorName: "default",
				CustomizedAcceleratorName:  "customizedAccelerator1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SpringCloudCustomizedAcceleratorIDInsensitively(v.Input)
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
		if actual.SpringName != v.Expected.SpringName {
			t.Fatalf("Expected %q but got %q for SpringName", v.Expected.SpringName, actual.SpringName)
		}
		if actual.ApplicationAcceleratorName != v.Expected.ApplicationAcceleratorName {
			t.Fatalf("Expected %q but got %q for ApplicationAcceleratorName", v.Expected.ApplicationAcceleratorName, actual.ApplicationAcceleratorName)
		}
		if actual.CustomizedAcceleratorName != v.Expected.CustomizedAcceleratorName {
			t.Fatalf("Expected %q but got %q for CustomizedAcceleratorName", v.Expected.CustomizedAcceleratorName, actual.CustomizedAcceleratorName)
		}
	}
}
