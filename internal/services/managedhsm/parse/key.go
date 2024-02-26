// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type KeyId struct {
	SubscriptionId string
	ResourceGroup  string
	ManagedHSMName string
	Name           string
	VersionName    string
}

func NewKeyID(subscriptionId, resourceGroup, managedHSMName, name, versionName string) KeyId {
	return KeyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ManagedHSMName: managedHSMName,
		Name:           name,
		VersionName:    versionName,
	}
}

func (id KeyId) String() string {
	segments := []string{
		fmt.Sprintf("Version Name %q", id.VersionName),
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Managed H S M Name %q", id.ManagedHSMName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Key", segmentsStr)
}

func (id KeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/managedHSMs/%s/keys/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedHSMName, id.Name, id.VersionName)
}

// KeyID parses a Key ID into an KeyId struct
func KeyID(input string) (*KeyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an Key ID: %+v", input, err)
	}

	resourceId := KeyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ManagedHSMName, err = id.PopSegment("managedHSMs"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("keys"); err != nil {
		return nil, err
	}
	if resourceId.VersionName, err = id.PopSegment("versions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
