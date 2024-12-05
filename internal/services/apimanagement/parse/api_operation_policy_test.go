// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ApiOperationPolicyId{}

func TestApiOperationPolicyIDFormatter(t *testing.T) {
	actual := NewApiOperationPolicyID("12345678-1234-9876-4563-123456789012", "resGroup1", "service1", "api1", "operation1", "policy1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/apis/api1/operations/operation1/policies/policy1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestApiOperationPolicyID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ApiOperationPolicyId
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
			// missing ServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/",
			Error: true,
		},

		{
			// missing value for ServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/",
			Error: true,
		},

		{
			// missing ApiName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/",
			Error: true,
		},

		{
			// missing value for ApiName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/apis/",
			Error: true,
		},

		{
			// missing OperationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/apis/api1/",
			Error: true,
		},

		{
			// missing value for OperationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/apis/api1/operations/",
			Error: true,
		},

		{
			// missing PolicyName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/apis/api1/operations/operation1/",
			Error: true,
		},

		{
			// missing value for PolicyName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/apis/api1/operations/operation1/policies/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/apis/api1/operations/operation1/policies/policy1",
			Expected: &ApiOperationPolicyId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resGroup1",
				ServiceName:    "service1",
				ApiName:        "api1",
				OperationName:  "operation1",
				PolicyName:     "policy1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.APIMANAGEMENT/SERVICE/SERVICE1/APIS/API1/OPERATIONS/OPERATION1/POLICIES/POLICY1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ApiOperationPolicyID(v.Input)
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
		if actual.ServiceName != v.Expected.ServiceName {
			t.Fatalf("Expected %q but got %q for ServiceName", v.Expected.ServiceName, actual.ServiceName)
		}
		if actual.ApiName != v.Expected.ApiName {
			t.Fatalf("Expected %q but got %q for ApiName", v.Expected.ApiName, actual.ApiName)
		}
		if actual.OperationName != v.Expected.OperationName {
			t.Fatalf("Expected %q but got %q for OperationName", v.Expected.OperationName, actual.OperationName)
		}
		if actual.PolicyName != v.Expected.PolicyName {
			t.Fatalf("Expected %q but got %q for PolicyName", v.Expected.PolicyName, actual.PolicyName)
		}
	}
}
