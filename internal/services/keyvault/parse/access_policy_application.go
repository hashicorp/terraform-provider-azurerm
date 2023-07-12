// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AccessPolicyApplicationId struct {
	SubscriptionId    string
	ResourceGroup     string
	VaultName         string
	ObjectIdName      string
	ApplicationIdName string
}

func NewAccessPolicyApplicationID(subscriptionId, resourceGroup, vaultName, objectIdName, applicationIdName string) AccessPolicyApplicationId {
	return AccessPolicyApplicationId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		VaultName:         vaultName,
		ObjectIdName:      objectIdName,
		ApplicationIdName: applicationIdName,
	}
}

func (id AccessPolicyApplicationId) String() string {
	segments := []string{
		fmt.Sprintf("Application Id Name %q", id.ApplicationIdName),
		fmt.Sprintf("Object Id Name %q", id.ObjectIdName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Access Policy Application", segmentsStr)
}

func (id AccessPolicyApplicationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/objectId/%s/applicationId/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.ObjectIdName, id.ApplicationIdName)
}

// AccessPolicyApplicationID parses a AccessPolicyApplication ID into an AccessPolicyApplicationId struct
func AccessPolicyApplicationID(input string) (*AccessPolicyApplicationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AccessPolicyApplication ID: %+v", input, err)
	}

	resourceId := AccessPolicyApplicationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VaultName, err = id.PopSegment("vaults"); err != nil {
		return nil, err
	}
	if resourceId.ObjectIdName, err = id.PopSegment("objectId"); err != nil {
		return nil, err
	}
	if resourceId.ApplicationIdName, err = id.PopSegment("applicationId"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
