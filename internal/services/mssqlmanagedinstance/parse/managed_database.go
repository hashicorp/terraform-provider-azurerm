// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ManagedDatabaseId struct {
	SubscriptionId      string
	ResourceGroup       string
	ManagedInstanceName string
	DatabaseName        string
}

func NewManagedDatabaseID(subscriptionId, resourceGroup, managedInstanceName, databaseName string) ManagedDatabaseId {
	return ManagedDatabaseId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		ManagedInstanceName: managedInstanceName,
		DatabaseName:        databaseName,
	}
}

func (id ManagedDatabaseId) String() string {
	segments := []string{
		fmt.Sprintf("Database Name %q", id.DatabaseName),
		fmt.Sprintf("Managed Instance Name %q", id.ManagedInstanceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Managed Database", segmentsStr)
}

func (id ManagedDatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/managedInstances/%s/databases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
}

// ManagedDatabaseID parses a ManagedDatabase ID into an ManagedDatabaseId struct
func ManagedDatabaseID(input string) (*ManagedDatabaseId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ManagedDatabase ID: %+v", input, err)
	}

	resourceId := ManagedDatabaseId{
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
	if resourceId.DatabaseName, err = id.PopSegment("databases"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
