// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ManagedHSMNestedKeyWithVersionId struct {
	SubscriptionId string
	ResourceGroup  string
	ManagedHSMName string
	KeyName        string
	VersionName    string
}

func NewManagedHSMNestedKeyWithVersionID(subscriptionId, resourceGroup, managedHSMName, keyName, versionName string) ManagedHSMNestedKeyWithVersionId {
	return ManagedHSMNestedKeyWithVersionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ManagedHSMName: managedHSMName,
		KeyName:        keyName,
		VersionName:    versionName,
	}
}

func (id ManagedHSMNestedKeyWithVersionId) String() string {
	segments := []string{
		fmt.Sprintf("Version Name %q", id.VersionName),
		fmt.Sprintf("Key Name %q", id.KeyName),
		fmt.Sprintf("Managed H S M Name %q", id.ManagedHSMName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Managed H S M Nested Key With Version", segmentsStr)
}

func (id ManagedHSMNestedKeyWithVersionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/managedHSMs/%s/keys/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedHSMName, id.KeyName, id.VersionName)
}

// ManagedHSMNestedKeyWithVersionID parses a ManagedHSMNestedKeyWithVersion ID into an ManagedHSMNestedKeyWithVersionId struct
func ManagedHSMNestedKeyWithVersionID(input string) (*ManagedHSMNestedKeyWithVersionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ManagedHSMNestedKeyWithVersion ID: %+v", input, err)
	}

	resourceId := ManagedHSMNestedKeyWithVersionId{
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
	if resourceId.KeyName, err = id.PopSegment("keys"); err != nil {
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
