// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SqlAzureActiveDirectoryAdministratorId struct {
	SubscriptionId    string
	ResourceGroup     string
	ServerName        string
	AdministratorName string
}

func NewSqlAzureActiveDirectoryAdministratorID(subscriptionId, resourceGroup, serverName, administratorName string) SqlAzureActiveDirectoryAdministratorId {
	return SqlAzureActiveDirectoryAdministratorId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		ServerName:        serverName,
		AdministratorName: administratorName,
	}
}

func (id SqlAzureActiveDirectoryAdministratorId) String() string {
	segments := []string{
		fmt.Sprintf("Administrator Name %q", id.AdministratorName),
		fmt.Sprintf("Server Name %q", id.ServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Sql Azure Active Directory Administrator", segmentsStr)
}

func (id SqlAzureActiveDirectoryAdministratorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/administrators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerName, id.AdministratorName)
}

// SqlAzureActiveDirectoryAdministratorID parses a SqlAzureActiveDirectoryAdministrator ID into an SqlAzureActiveDirectoryAdministratorId struct
func SqlAzureActiveDirectoryAdministratorID(input string) (*SqlAzureActiveDirectoryAdministratorId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SqlAzureActiveDirectoryAdministrator ID: %+v", input, err)
	}

	resourceId := SqlAzureActiveDirectoryAdministratorId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerName, err = id.PopSegment("servers"); err != nil {
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

// SqlAzureActiveDirectoryAdministratorIDInsensitively parses an SqlAzureActiveDirectoryAdministrator ID into an SqlAzureActiveDirectoryAdministratorId struct, insensitively
// This should only be used to parse an ID for rewriting, the SqlAzureActiveDirectoryAdministratorID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SqlAzureActiveDirectoryAdministratorIDInsensitively(input string) (*SqlAzureActiveDirectoryAdministratorId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SqlAzureActiveDirectoryAdministratorId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'servers' segment
	serversKey := "servers"
	for key := range id.Path {
		if strings.EqualFold(key, serversKey) {
			serversKey = key
			break
		}
	}
	if resourceId.ServerName, err = id.PopSegment(serversKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'administrators' segment
	administratorsKey := "administrators"
	for key := range id.Path {
		if strings.EqualFold(key, administratorsKey) {
			administratorsKey = key
			break
		}
	}
	if resourceId.AdministratorName, err = id.PopSegment(administratorsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
