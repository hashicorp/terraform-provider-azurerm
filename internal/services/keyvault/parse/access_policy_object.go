// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AccessPolicyObjectId struct {
	SubscriptionId string
	ResourceGroup  string
	VaultName      string
	ObjectIdName   string
}

func NewAccessPolicyObjectID(subscriptionId, resourceGroup, vaultName, objectIdName string) AccessPolicyObjectId {
	return AccessPolicyObjectId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		VaultName:      vaultName,
		ObjectIdName:   objectIdName,
	}
}

func (id AccessPolicyObjectId) String() string {
	segments := []string{
		fmt.Sprintf("Object Id Name %q", id.ObjectIdName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Access Policy Object", segmentsStr)
}

func (id AccessPolicyObjectId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/objectId/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.ObjectIdName)
}

// AccessPolicyObjectID parses a AccessPolicyObject ID into an AccessPolicyObjectId struct
func AccessPolicyObjectID(input string) (*AccessPolicyObjectId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AccessPolicyObject ID: %+v", input, err)
	}

	resourceId := AccessPolicyObjectId{
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

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
