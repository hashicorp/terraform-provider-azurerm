// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudDeploymentId struct {
	SubscriptionId string
	ResourceGroup  string
	SpringName     string
	AppName        string
	DeploymentName string
}

func NewSpringCloudDeploymentID(subscriptionId, resourceGroup, springName, appName, deploymentName string) SpringCloudDeploymentId {
	return SpringCloudDeploymentId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SpringName:     springName,
		AppName:        appName,
		DeploymentName: deploymentName,
	}
}

func (id SpringCloudDeploymentId) String() string {
	segments := []string{
		fmt.Sprintf("Deployment Name %q", id.DeploymentName),
		fmt.Sprintf("App Name %q", id.AppName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Deployment", segmentsStr)
}

func (id SpringCloudDeploymentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/apps/%s/deployments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.AppName, id.DeploymentName)
}

// SpringCloudDeploymentID parses a SpringCloudDeployment ID into an SpringCloudDeploymentId struct
func SpringCloudDeploymentID(input string) (*SpringCloudDeploymentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudDeployment ID: %+v", input, err)
	}

	resourceId := SpringCloudDeploymentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SpringName, err = id.PopSegment("spring"); err != nil {
		return nil, err
	}
	if resourceId.AppName, err = id.PopSegment("apps"); err != nil {
		return nil, err
	}
	if resourceId.DeploymentName, err = id.PopSegment("deployments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SpringCloudDeploymentIDInsensitively parses an SpringCloudDeployment ID into an SpringCloudDeploymentId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudDeploymentID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudDeploymentIDInsensitively(input string) (*SpringCloudDeploymentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudDeploymentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'spring' segment
	springKey := "spring"
	for key := range id.Path {
		if strings.EqualFold(key, springKey) {
			springKey = key
			break
		}
	}
	if resourceId.SpringName, err = id.PopSegment(springKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'apps' segment
	appsKey := "apps"
	for key := range id.Path {
		if strings.EqualFold(key, appsKey) {
			appsKey = key
			break
		}
	}
	if resourceId.AppName, err = id.PopSegment(appsKey); err != nil {
		return nil, err
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
