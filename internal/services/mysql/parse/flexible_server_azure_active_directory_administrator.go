// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FlexibleServerAzureActiveDirectoryAdministratorId struct {
	SubscriptionId     string
	ResourceGroup      string
	FlexibleServerName string
	AdministratorName  string
}

func NewFlexibleServerAzureActiveDirectoryAdministratorID(subscriptionId, resourceGroup, flexibleServerName, administratorName string) FlexibleServerAzureActiveDirectoryAdministratorId {
	return FlexibleServerAzureActiveDirectoryAdministratorId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		FlexibleServerName: flexibleServerName,
		AdministratorName:  administratorName,
	}
}

func (id FlexibleServerAzureActiveDirectoryAdministratorId) String() string {
	segments := []string{
		fmt.Sprintf("Administrator Name %q", id.AdministratorName),
		fmt.Sprintf("Flexible Server Name %q", id.FlexibleServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Flexible Server Azure Active Directory Administrator", segmentsStr)
}

func (id FlexibleServerAzureActiveDirectoryAdministratorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforMySQL/flexibleServers/%s/administrators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FlexibleServerName, id.AdministratorName)
}

// FlexibleServerAzureActiveDirectoryAdministratorID parses a FlexibleServerAzureActiveDirectoryAdministrator ID into an FlexibleServerAzureActiveDirectoryAdministratorId struct
func FlexibleServerAzureActiveDirectoryAdministratorID(input string) (*FlexibleServerAzureActiveDirectoryAdministratorId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an FlexibleServerAzureActiveDirectoryAdministrator ID: %+v", input, err)
	}

	resourceId := FlexibleServerAzureActiveDirectoryAdministratorId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FlexibleServerName, err = id.PopSegment("flexibleServers"); err != nil {
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
