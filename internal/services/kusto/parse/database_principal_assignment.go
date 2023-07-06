// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DatabasePrincipalAssignmentId struct {
	SubscriptionId          string
	ResourceGroup           string
	ClusterName             string
	DatabaseName            string
	PrincipalAssignmentName string
}

func NewDatabasePrincipalAssignmentID(subscriptionId, resourceGroup, clusterName, databaseName, principalAssignmentName string) DatabasePrincipalAssignmentId {
	return DatabasePrincipalAssignmentId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ClusterName:             clusterName,
		DatabaseName:            databaseName,
		PrincipalAssignmentName: principalAssignmentName,
	}
}

func (id DatabasePrincipalAssignmentId) String() string {
	segments := []string{
		fmt.Sprintf("Principal Assignment Name %q", id.PrincipalAssignmentName),
		fmt.Sprintf("Database Name %q", id.DatabaseName),
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Database Principal Assignment", segmentsStr)
}

func (id DatabasePrincipalAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/clusters/%s/databases/%s/principalAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.PrincipalAssignmentName)
}

// DatabasePrincipalAssignmentID parses a DatabasePrincipalAssignment ID into an DatabasePrincipalAssignmentId struct
func DatabasePrincipalAssignmentID(input string) (*DatabasePrincipalAssignmentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an DatabasePrincipalAssignment ID: %+v", input, err)
	}

	resourceId := DatabasePrincipalAssignmentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ClusterName, err = id.PopSegment("clusters"); err != nil {
		return nil, err
	}
	if resourceId.DatabaseName, err = id.PopSegment("databases"); err != nil {
		return nil, err
	}
	if resourceId.PrincipalAssignmentName, err = id.PopSegment("principalAssignments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// DatabasePrincipalAssignmentIDInsensitively parses an DatabasePrincipalAssignment ID into an DatabasePrincipalAssignmentId struct, insensitively
// This should only be used to parse an ID for rewriting, the DatabasePrincipalAssignmentID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func DatabasePrincipalAssignmentIDInsensitively(input string) (*DatabasePrincipalAssignmentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DatabasePrincipalAssignmentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'clusters' segment
	clustersKey := "clusters"
	for key := range id.Path {
		if strings.EqualFold(key, clustersKey) {
			clustersKey = key
			break
		}
	}
	if resourceId.ClusterName, err = id.PopSegment(clustersKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'databases' segment
	databasesKey := "databases"
	for key := range id.Path {
		if strings.EqualFold(key, databasesKey) {
			databasesKey = key
			break
		}
	}
	if resourceId.DatabaseName, err = id.PopSegment(databasesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'principalAssignments' segment
	principalAssignmentsKey := "principalAssignments"
	for key := range id.Path {
		if strings.EqualFold(key, principalAssignmentsKey) {
			principalAssignmentsKey = key
			break
		}
	}
	if resourceId.PrincipalAssignmentName, err = id.PopSegment(principalAssignmentsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
