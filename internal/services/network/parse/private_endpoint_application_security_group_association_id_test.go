// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/privateendpoints"
)

func TestPrivateEndpointApplicationSecurityGroupAssociationID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Expect *PrivateEndpointApplicationSecurityGroupAssociationId
		Error  bool
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "One Segment",
			Input: "hello",
			Error: true,
		},
		{
			Name:  "Two Segments Invalid ID's",
			Input: "hello|world",
			Error: true,
		},
		{
			Name:  "Missing ASG Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/privateEndpoints/endpoints1",
			Error: true,
		},
		{
			Name:  "Private Endpoint Id",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/privateEndpoints/endpoints1",
			Error: true,
		},
		{
			Name:  "Application Security Group ID",
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationSecurityGroups/securityGroup1",
			Error: true,
		},
		{
			Name:  "Nat Gateway / Public IP Association ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/privateEndpoints/endpoints1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/applicationSecurityGroups/securityGroup1",
			Error: false,
			Expect: &PrivateEndpointApplicationSecurityGroupAssociationId{
				ApplicationSecurityGroupId: applicationsecuritygroups.ApplicationSecurityGroupId{
					ResourceGroupName:            "mygroup1",
					SubscriptionId:               "00000000-0000-0000-0000-000000000000",
					ApplicationSecurityGroupName: "securityGroup1",
				},
				PrivateEndpointId: privateendpoints.PrivateEndpointId{
					ResourceGroupName:   "group1",
					SubscriptionId:      "00000000-0000-0000-0000-000000000000",
					PrivateEndpointName: "endpoints1",
				},
			},
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := PrivateEndpointApplicationSecurityGroupAssociationID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.PrivateEndpointId.PrivateEndpointName != v.Expect.PrivateEndpointId.PrivateEndpointName {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.PrivateEndpointId.PrivateEndpointName, actual.PrivateEndpointId.PrivateEndpointName)
		}

		if actual.ApplicationSecurityGroupId.ResourceGroupName != v.Expect.ApplicationSecurityGroupId.ResourceGroupName {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ApplicationSecurityGroupId.ResourceGroupName, actual.ApplicationSecurityGroupId.ResourceGroupName)
		}
	}
}
