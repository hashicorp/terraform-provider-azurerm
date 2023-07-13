// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ResourceGroupTemplateDeploymentId struct {
	SubscriptionId string
	ResourceGroup  string
	DeploymentName string
}

func NewResourceGroupTemplateDeploymentID(subscriptionId, resourceGroup, deploymentName string) ResourceGroupTemplateDeploymentId {
	return ResourceGroupTemplateDeploymentId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		DeploymentName: deploymentName,
	}
}

func (id ResourceGroupTemplateDeploymentId) String() string {
	segments := []string{
		fmt.Sprintf("Deployment Name %q", id.DeploymentName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Resource Group Template Deployment", segmentsStr)
}

func (id ResourceGroupTemplateDeploymentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Resources/deployments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DeploymentName)
}

// ResourceGroupTemplateDeploymentID parses a ResourceGroupTemplateDeployment ID into an ResourceGroupTemplateDeploymentId struct
func ResourceGroupTemplateDeploymentID(input string) (*ResourceGroupTemplateDeploymentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ResourceGroupTemplateDeployment ID: %+v", input, err)
	}

	resourceId := ResourceGroupTemplateDeploymentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DeploymentName, err = id.PopSegment("deployments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ResourceGroupTemplateDeploymentIDInsensitively parses an ResourceGroupTemplateDeployment ID into an ResourceGroupTemplateDeploymentId struct, insensitively
// This should only be used to parse an ID for rewriting, the ResourceGroupTemplateDeploymentID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func ResourceGroupTemplateDeploymentIDInsensitively(input string) (*ResourceGroupTemplateDeploymentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ResourceGroupTemplateDeploymentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'deployments' segment
	deploymentsKey := "deployments"
	for key := range id.Path {
		if strings.EqualFold(key, deploymentsKey) {
			deploymentsKey = key
			break
		}
	}
	if resourceId.DeploymentName, err = id.PopSegment(deploymentsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
