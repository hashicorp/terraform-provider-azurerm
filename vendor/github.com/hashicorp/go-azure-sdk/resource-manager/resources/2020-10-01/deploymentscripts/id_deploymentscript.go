package deploymentscripts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DeploymentScriptId{}

// DeploymentScriptId is a struct representing the Resource ID for a Deployment Script
type DeploymentScriptId struct {
	SubscriptionId       string
	ResourceGroupName    string
	DeploymentScriptName string
}

// NewDeploymentScriptID returns a new DeploymentScriptId struct
func NewDeploymentScriptID(subscriptionId string, resourceGroupName string, deploymentScriptName string) DeploymentScriptId {
	return DeploymentScriptId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		DeploymentScriptName: deploymentScriptName,
	}
}

// ParseDeploymentScriptID parses 'input' into a DeploymentScriptId
func ParseDeploymentScriptID(input string) (*DeploymentScriptId, error) {
	parser := resourceids.NewParserFromResourceIdType(DeploymentScriptId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DeploymentScriptId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DeploymentScriptName, ok = parsed.Parsed["deploymentScriptName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "deploymentScriptName", *parsed)
	}

	return &id, nil
}

// ParseDeploymentScriptIDInsensitively parses 'input' case-insensitively into a DeploymentScriptId
// note: this method should only be used for API response data and not user input
func ParseDeploymentScriptIDInsensitively(input string) (*DeploymentScriptId, error) {
	parser := resourceids.NewParserFromResourceIdType(DeploymentScriptId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DeploymentScriptId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DeploymentScriptName, ok = parsed.Parsed["deploymentScriptName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "deploymentScriptName", *parsed)
	}

	return &id, nil
}

// ValidateDeploymentScriptID checks that 'input' can be parsed as a Deployment Script ID
func ValidateDeploymentScriptID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeploymentScriptID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deployment Script ID
func (id DeploymentScriptId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Resources/deploymentScripts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DeploymentScriptName)
}

// Segments returns a slice of Resource ID Segments which comprise this Deployment Script ID
func (id DeploymentScriptId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftResources", "Microsoft.Resources", "Microsoft.Resources"),
		resourceids.StaticSegment("staticDeploymentScripts", "deploymentScripts", "deploymentScripts"),
		resourceids.UserSpecifiedSegment("deploymentScriptName", "deploymentScriptValue"),
	}
}

// String returns a human-readable description of this Deployment Script ID
func (id DeploymentScriptId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Deployment Script Name: %q", id.DeploymentScriptName),
	}
	return fmt.Sprintf("Deployment Script (%s)", strings.Join(components, "\n"))
}
