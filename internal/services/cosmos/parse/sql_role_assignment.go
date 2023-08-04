// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SqlRoleAssignmentId struct {
	SubscriptionId      string
	ResourceGroup       string
	DatabaseAccountName string
	Name                string
}

func NewSqlRoleAssignmentID(subscriptionId, resourceGroup, databaseAccountName, name string) SqlRoleAssignmentId {
	return SqlRoleAssignmentId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		DatabaseAccountName: databaseAccountName,
		Name:                name,
	}
}

func (id SqlRoleAssignmentId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Database Account Name %q", id.DatabaseAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Sql Role Assignment", segmentsStr)
}

func (id SqlRoleAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/sqlRoleAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DatabaseAccountName, id.Name)
}

// SqlRoleAssignmentID parses a SqlRoleAssignment ID into an SqlRoleAssignmentId struct
func SqlRoleAssignmentID(input string) (*SqlRoleAssignmentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SqlRoleAssignment ID: %+v", input, err)
	}

	resourceId := SqlRoleAssignmentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DatabaseAccountName, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("sqlRoleAssignments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
