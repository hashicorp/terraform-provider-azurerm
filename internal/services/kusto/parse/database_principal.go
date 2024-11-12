// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DatabasePrincipalId struct {
	SubscriptionId string
	ResourceGroup  string
	ClusterName    string
	DatabaseName   string
	RoleName       string
	FQNName        string
}

func NewDatabasePrincipalID(subscriptionId, resourceGroup, clusterName, databaseName, roleName, fQNName string) DatabasePrincipalId {
	return DatabasePrincipalId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ClusterName:    clusterName,
		DatabaseName:   databaseName,
		RoleName:       roleName,
		FQNName:        fQNName,
	}
}

func (id DatabasePrincipalId) String() string {
	segments := []string{
		fmt.Sprintf("FQN Name %q", id.FQNName),
		fmt.Sprintf("Role Name %q", id.RoleName),
		fmt.Sprintf("Database Name %q", id.DatabaseName),
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Database Principal", segmentsStr)
}

func (id DatabasePrincipalId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/Clusters/%s/Databases/%s/Role/%s/FQN/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.RoleName, id.FQNName)
}

// DatabasePrincipalID parses a DatabasePrincipal ID into an DatabasePrincipalId struct
func DatabasePrincipalID(input string) (*DatabasePrincipalId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an DatabasePrincipal ID: %+v", input, err)
	}

	resourceId := DatabasePrincipalId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ClusterName, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}
	if resourceId.DatabaseName, err = id.PopSegment("Databases"); err != nil {
		return nil, err
	}
	if resourceId.RoleName, err = id.PopSegment("Role"); err != nil {
		return nil, err
	}
	if resourceId.FQNName, err = id.PopSegment("FQN"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
