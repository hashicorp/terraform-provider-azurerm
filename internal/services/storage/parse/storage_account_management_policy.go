// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type StorageAccountManagementPolicyId struct {
	SubscriptionId       string
	ResourceGroup        string
	StorageAccountName   string
	ManagementPolicyName string
}

func NewStorageAccountManagementPolicyID(subscriptionId, resourceGroup, storageAccountName, managementPolicyName string) StorageAccountManagementPolicyId {
	return StorageAccountManagementPolicyId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		StorageAccountName:   storageAccountName,
		ManagementPolicyName: managementPolicyName,
	}
}

func (id StorageAccountManagementPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Management Policy Name %q", id.ManagementPolicyName),
		fmt.Sprintf("Storage Account Name %q", id.StorageAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Account Management Policy", segmentsStr)
}

func (id StorageAccountManagementPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/managementPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.ManagementPolicyName)
}

// StorageAccountManagementPolicyID parses a StorageAccountManagementPolicy ID into an StorageAccountManagementPolicyId struct
func StorageAccountManagementPolicyID(input string) (*StorageAccountManagementPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an StorageAccountManagementPolicy ID: %+v", input, err)
	}

	resourceId := StorageAccountManagementPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.StorageAccountName, err = id.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}
	if resourceId.ManagementPolicyName, err = id.PopSegment("managementPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
