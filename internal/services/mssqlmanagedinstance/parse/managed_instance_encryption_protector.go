// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ManagedInstanceEncryptionProtectorId struct {
	SubscriptionId          string
	ResourceGroup           string
	ManagedInstanceName     string
	EncryptionProtectorName string
}

func NewManagedInstanceEncryptionProtectorID(subscriptionId, resourceGroup, managedInstanceName, encryptionProtectorName string) ManagedInstanceEncryptionProtectorId {
	return ManagedInstanceEncryptionProtectorId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ManagedInstanceName:     managedInstanceName,
		EncryptionProtectorName: encryptionProtectorName,
	}
}

func (id ManagedInstanceEncryptionProtectorId) String() string {
	segments := []string{
		fmt.Sprintf("Encryption Protector Name %q", id.EncryptionProtectorName),
		fmt.Sprintf("Managed Instance Name %q", id.ManagedInstanceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Managed Instance Encryption Protector", segmentsStr)
}

func (id ManagedInstanceEncryptionProtectorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/managedInstances/%s/encryptionProtector/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName, id.EncryptionProtectorName)
}

// ManagedInstanceEncryptionProtectorID parses a ManagedInstanceEncryptionProtector ID into an ManagedInstanceEncryptionProtectorId struct
func ManagedInstanceEncryptionProtectorID(input string) (*ManagedInstanceEncryptionProtectorId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ManagedInstanceEncryptionProtector ID: %+v", input, err)
	}

	resourceId := ManagedInstanceEncryptionProtectorId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ManagedInstanceName, err = id.PopSegment("managedInstances"); err != nil {
		return nil, err
	}
	if resourceId.EncryptionProtectorName, err = id.PopSegment("encryptionProtector"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
