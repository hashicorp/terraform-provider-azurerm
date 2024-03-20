// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ManagedHSMNestedKeyVersionlessId struct {
	SubscriptionId string
	ResourceGroup  string
	ManagedHSMName string
	KeyName        string
}

func NewManagedHSMNestedKeyVersionlessID(subscriptionId, resourceGroup, managedHSMName, keyName string) ManagedHSMNestedKeyVersionlessId {
	return ManagedHSMNestedKeyVersionlessId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ManagedHSMName: managedHSMName,
		KeyName:        keyName,
	}
}

func (id ManagedHSMNestedKeyVersionlessId) String() string {
	segments := []string{
		fmt.Sprintf("Key Name %q", id.KeyName),
		fmt.Sprintf("Managed H S M Name %q", id.ManagedHSMName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Managed H S M Nested Key Versionless", segmentsStr)
}

func (id ManagedHSMNestedKeyVersionlessId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/managedHSMs/%s/keys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedHSMName, id.KeyName)
}

// ManagedHSMNestedKeyVersionlessID parses a ManagedHSMNestedKeyVersionless ID into an ManagedHSMNestedKeyVersionlessId struct
func ManagedHSMNestedKeyVersionlessID(input string) (*ManagedHSMNestedKeyVersionlessId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ManagedHSMNestedKeyVersionless ID: %+v", input, err)
	}

	resourceId := ManagedHSMNestedKeyVersionlessId{
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

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
