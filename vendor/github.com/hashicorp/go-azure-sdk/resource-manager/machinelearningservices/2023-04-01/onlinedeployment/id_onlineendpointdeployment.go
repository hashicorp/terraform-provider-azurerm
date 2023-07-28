package onlinedeployment

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = OnlineEndpointDeploymentId{}

// OnlineEndpointDeploymentId is a struct representing the Resource ID for a Online Endpoint Deployment
type OnlineEndpointDeploymentId struct {
	SubscriptionId     string
	ResourceGroupName  string
	WorkspaceName      string
	OnlineEndpointName string
	DeploymentName     string
}

// NewOnlineEndpointDeploymentID returns a new OnlineEndpointDeploymentId struct
func NewOnlineEndpointDeploymentID(subscriptionId string, resourceGroupName string, workspaceName string, onlineEndpointName string, deploymentName string) OnlineEndpointDeploymentId {
	return OnlineEndpointDeploymentId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		WorkspaceName:      workspaceName,
		OnlineEndpointName: onlineEndpointName,
		DeploymentName:     deploymentName,
	}
}

// ParseOnlineEndpointDeploymentID parses 'input' into a OnlineEndpointDeploymentId
func ParseOnlineEndpointDeploymentID(input string) (*OnlineEndpointDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(OnlineEndpointDeploymentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OnlineEndpointDeploymentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.OnlineEndpointName, ok = parsed.Parsed["onlineEndpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "onlineEndpointName", *parsed)
	}

	if id.DeploymentName, ok = parsed.Parsed["deploymentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "deploymentName", *parsed)
	}

	return &id, nil
}

// ParseOnlineEndpointDeploymentIDInsensitively parses 'input' case-insensitively into a OnlineEndpointDeploymentId
// note: this method should only be used for API response data and not user input
func ParseOnlineEndpointDeploymentIDInsensitively(input string) (*OnlineEndpointDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(OnlineEndpointDeploymentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OnlineEndpointDeploymentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.OnlineEndpointName, ok = parsed.Parsed["onlineEndpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "onlineEndpointName", *parsed)
	}

	if id.DeploymentName, ok = parsed.Parsed["deploymentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "deploymentName", *parsed)
	}

	return &id, nil
}

// ValidateOnlineEndpointDeploymentID checks that 'input' can be parsed as a Online Endpoint Deployment ID
func ValidateOnlineEndpointDeploymentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOnlineEndpointDeploymentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Online Endpoint Deployment ID
func (id OnlineEndpointDeploymentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MachineLearningServices/workspaces/%s/onlineEndpoints/%s/deployments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.OnlineEndpointName, id.DeploymentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Online Endpoint Deployment ID
func (id OnlineEndpointDeploymentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMachineLearningServices", "Microsoft.MachineLearningServices", "Microsoft.MachineLearningServices"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticOnlineEndpoints", "onlineEndpoints", "onlineEndpoints"),
		resourceids.UserSpecifiedSegment("onlineEndpointName", "onlineEndpointValue"),
		resourceids.StaticSegment("staticDeployments", "deployments", "deployments"),
		resourceids.UserSpecifiedSegment("deploymentName", "deploymentValue"),
	}
}

// String returns a human-readable description of this Online Endpoint Deployment ID
func (id OnlineEndpointDeploymentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Online Endpoint Name: %q", id.OnlineEndpointName),
		fmt.Sprintf("Deployment Name: %q", id.DeploymentName),
	}
	return fmt.Sprintf("Online Endpoint Deployment (%s)", strings.Join(components, "\n"))
}
