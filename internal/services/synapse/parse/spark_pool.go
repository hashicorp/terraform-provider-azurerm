// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SparkPoolId struct {
	SubscriptionId  string
	ResourceGroup   string
	WorkspaceName   string
	BigDataPoolName string
}

func NewSparkPoolID(subscriptionId, resourceGroup, workspaceName, bigDataPoolName string) SparkPoolId {
	return SparkPoolId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		WorkspaceName:   workspaceName,
		BigDataPoolName: bigDataPoolName,
	}
}

func (id SparkPoolId) String() string {
	segments := []string{
		fmt.Sprintf("Big Data Pool Name %q", id.BigDataPoolName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spark Pool", segmentsStr)
}

func (id SparkPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Synapse/workspaces/%s/bigDataPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.BigDataPoolName)
}

// SparkPoolID parses a SparkPool ID into an SparkPoolId struct
func SparkPoolID(input string) (*SparkPoolId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SparkPool ID: %+v", input, err)
	}

	resourceId := SparkPoolId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if resourceId.BigDataPoolName, err = id.PopSegment("bigDataPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SparkPoolIDInsensitively parses an SparkPool ID into an SparkPoolId struct, insensitively
// This should only be used to parse an ID for rewriting, the SparkPoolID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SparkPoolIDInsensitively(input string) (*SparkPoolId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SparkPoolId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'workspaces' segment
	workspacesKey := "workspaces"
	for key := range id.Path {
		if strings.EqualFold(key, workspacesKey) {
			workspacesKey = key
			break
		}
	}
	if resourceId.WorkspaceName, err = id.PopSegment(workspacesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'bigDataPools' segment
	bigDataPoolsKey := "bigDataPools"
	for key := range id.Path {
		if strings.EqualFold(key, bigDataPoolsKey) {
			bigDataPoolsKey = key
			break
		}
	}
	if resourceId.BigDataPoolName, err = id.PopSegment(bigDataPoolsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
