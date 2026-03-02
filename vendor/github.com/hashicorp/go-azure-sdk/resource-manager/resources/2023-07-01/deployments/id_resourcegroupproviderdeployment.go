package deployments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ResourceGroupProviderDeploymentId{})
}

var _ resourceids.ResourceId = &ResourceGroupProviderDeploymentId{}

// ResourceGroupProviderDeploymentId is a struct representing the Resource ID for a Resource Group Provider Deployment
type ResourceGroupProviderDeploymentId struct {
	SubscriptionId    string
	ResourceGroupName string
	DeploymentName    string
}

// NewResourceGroupProviderDeploymentID returns a new ResourceGroupProviderDeploymentId struct
func NewResourceGroupProviderDeploymentID(subscriptionId string, resourceGroupName string, deploymentName string) ResourceGroupProviderDeploymentId {
	return ResourceGroupProviderDeploymentId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DeploymentName:    deploymentName,
	}
}

// ParseResourceGroupProviderDeploymentID parses 'input' into a ResourceGroupProviderDeploymentId
func ParseResourceGroupProviderDeploymentID(input string) (*ResourceGroupProviderDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResourceGroupProviderDeploymentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResourceGroupProviderDeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseResourceGroupProviderDeploymentIDInsensitively parses 'input' case-insensitively into a ResourceGroupProviderDeploymentId
// note: this method should only be used for API response data and not user input
func ParseResourceGroupProviderDeploymentIDInsensitively(input string) (*ResourceGroupProviderDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResourceGroupProviderDeploymentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResourceGroupProviderDeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ResourceGroupProviderDeploymentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DeploymentName, ok = input.Parsed["deploymentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deploymentName", input)
	}

	return nil
}

// ValidateResourceGroupProviderDeploymentID checks that 'input' can be parsed as a Resource Group Provider Deployment ID
func ValidateResourceGroupProviderDeploymentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseResourceGroupProviderDeploymentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Resource Group Provider Deployment ID
func (id ResourceGroupProviderDeploymentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Resources/deployments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DeploymentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Resource Group Provider Deployment ID
func (id ResourceGroupProviderDeploymentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftResources", "Microsoft.Resources", "Microsoft.Resources"),
		resourceids.StaticSegment("staticDeployments", "deployments", "deployments"),
		resourceids.UserSpecifiedSegment("deploymentName", "deploymentName"),
	}
}

// String returns a human-readable description of this Resource Group Provider Deployment ID
func (id ResourceGroupProviderDeploymentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Deployment Name: %q", id.DeploymentName),
	}
	return fmt.Sprintf("Resource Group Provider Deployment (%s)", strings.Join(components, "\n"))
}
