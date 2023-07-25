// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = TrustedRootCertificateId{}

func TestTrustedRootCertificateIDFormatter(t *testing.T) {
	actual := NewTrustedRootCertificateID("12345678-1234-9876-4563-123456789012", "group1", "applicationGateway1", "rootCert1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/trustedRootCertificates/rootCert1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestTrustedRootCertificateID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *TrustedRootCertificateId
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
			// missing ApplicationGatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for ApplicationGatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/trustedRootCertificates/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/trustedRootCertificates/rootCert1",
			Expected: &TrustedRootCertificateId{
				SubscriptionId:         "12345678-1234-9876-4563-123456789012",
				ResourceGroup:          "group1",
				ApplicationGatewayName: "applicationGateway1",
				Name:                   "rootCert1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/MICROSOFT.NETWORK/APPLICATIONGATEWAYS/APPLICATIONGATEWAY1/TRUSTEDROOTCERTIFICATES/ROOTCERT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := TrustedRootCertificateID(v.Input)
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
		if actual.ApplicationGatewayName != v.Expected.ApplicationGatewayName {
			t.Fatalf("Expected %q but got %q for ApplicationGatewayName", v.Expected.ApplicationGatewayName, actual.ApplicationGatewayName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}

func TestTrustedRootCertificateIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *TrustedRootCertificateId
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
			// missing ApplicationGatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for ApplicationGatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/trustedRootCertificates/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/trustedRootCertificates/rootCert1",
			Expected: &TrustedRootCertificateId{
				SubscriptionId:         "12345678-1234-9876-4563-123456789012",
				ResourceGroup:          "group1",
				ApplicationGatewayName: "applicationGateway1",
				Name:                   "rootCert1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationgateways/applicationGateway1/trustedrootcertificates/rootCert1",
			Expected: &TrustedRootCertificateId{
				SubscriptionId:         "12345678-1234-9876-4563-123456789012",
				ResourceGroup:          "group1",
				ApplicationGatewayName: "applicationGateway1",
				Name:                   "rootCert1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/APPLICATIONGATEWAYS/applicationGateway1/TRUSTEDROOTCERTIFICATES/rootCert1",
			Expected: &TrustedRootCertificateId{
				SubscriptionId:         "12345678-1234-9876-4563-123456789012",
				ResourceGroup:          "group1",
				ApplicationGatewayName: "applicationGateway1",
				Name:                   "rootCert1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/ApPlIcAtIoNgAtEwAyS/applicationGateway1/TrUsTeDrOoTcErTiFiCaTeS/rootCert1",
			Expected: &TrustedRootCertificateId{
				SubscriptionId:         "12345678-1234-9876-4563-123456789012",
				ResourceGroup:          "group1",
				ApplicationGatewayName: "applicationGateway1",
				Name:                   "rootCert1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := TrustedRootCertificateIDInsensitively(v.Input)
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
		if actual.ApplicationGatewayName != v.Expected.ApplicationGatewayName {
			t.Fatalf("Expected %q but got %q for ApplicationGatewayName", v.Expected.ApplicationGatewayName, actual.ApplicationGatewayName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
