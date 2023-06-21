// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ClusterPrincipalAssignmentId struct {
	SubscriptionId          string
	ResourceGroup           string
	ClusterName             string
	PrincipalAssignmentName string
}

func NewClusterPrincipalAssignmentID(subscriptionId, resourceGroup, clusterName, principalAssignmentName string) ClusterPrincipalAssignmentId {
	return ClusterPrincipalAssignmentId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ClusterName:             clusterName,
		PrincipalAssignmentName: principalAssignmentName,
	}
}

func (id ClusterPrincipalAssignmentId) String() string {
	segments := []string{
		fmt.Sprintf("Principal Assignment Name %q", id.PrincipalAssignmentName),
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Cluster Principal Assignment", segmentsStr)
}

func (id ClusterPrincipalAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/clusters/%s/principalAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.PrincipalAssignmentName)
}

// ClusterPrincipalAssignmentID parses a ClusterPrincipalAssignment ID into an ClusterPrincipalAssignmentId struct
func ClusterPrincipalAssignmentID(input string) (*ClusterPrincipalAssignmentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ClusterPrincipalAssignment ID: %+v", input, err)
	}

	resourceId := ClusterPrincipalAssignmentId{
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
	if resourceId.PrincipalAssignmentName, err = id.PopSegment("principalAssignments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ClusterPrincipalAssignmentIDInsensitively parses an ClusterPrincipalAssignment ID into an ClusterPrincipalAssignmentId struct, insensitively
// This should only be used to parse an ID for rewriting, the ClusterPrincipalAssignmentID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func ClusterPrincipalAssignmentIDInsensitively(input string) (*ClusterPrincipalAssignmentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ClusterPrincipalAssignmentId{
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
