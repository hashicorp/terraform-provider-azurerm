// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ManagedInstanceAzureActiveDirectoryAdministratorId struct {
	SubscriptionId      string
	ResourceGroup       string
	ManagedInstanceName string
	AdministratorName   string
}

func NewManagedInstanceAzureActiveDirectoryAdministratorID(subscriptionId, resourceGroup, managedInstanceName, administratorName string) ManagedInstanceAzureActiveDirectoryAdministratorId {
	return ManagedInstanceAzureActiveDirectoryAdministratorId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		ManagedInstanceName: managedInstanceName,
		AdministratorName:   administratorName,
	}
}

func (id ManagedInstanceAzureActiveDirectoryAdministratorId) String() string {
	segments := []string{
		fmt.Sprintf("Administrator Name %q", id.AdministratorName),
		fmt.Sprintf("Managed Instance Name %q", id.ManagedInstanceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Managed Instance Azure Active Directory Administrator", segmentsStr)
}

func (id ManagedInstanceAzureActiveDirectoryAdministratorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/managedInstances/%s/administrators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName, id.AdministratorName)
}

// ManagedInstanceAzureActiveDirectoryAdministratorID parses a ManagedInstanceAzureActiveDirectoryAdministrator ID into an ManagedInstanceAzureActiveDirectoryAdministratorId struct
func ManagedInstanceAzureActiveDirectoryAdministratorID(input string) (*ManagedInstanceAzureActiveDirectoryAdministratorId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ManagedInstanceAzureActiveDirectoryAdministrator ID: %+v", input, err)
	}

	resourceId := ManagedInstanceAzureActiveDirectoryAdministratorId{
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
	if resourceId.AdministratorName, err = id.PopSegment("administrators"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
