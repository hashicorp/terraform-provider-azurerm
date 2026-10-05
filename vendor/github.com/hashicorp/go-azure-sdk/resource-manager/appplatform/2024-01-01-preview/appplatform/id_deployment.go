package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DeploymentId{})
}

var _ resourceids.ResourceId = &DeploymentId{}

// DeploymentId is a struct representing the Resource ID for a Deployment
type DeploymentId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpringName        string
	AppName           string
	DeploymentName    string
}

// NewDeploymentID returns a new DeploymentId struct
func NewDeploymentID(subscriptionId string, resourceGroupName string, springName string, appName string, deploymentName string) DeploymentId {
	return DeploymentId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpringName:        springName,
		AppName:           appName,
		DeploymentName:    deploymentName,
	}
}

// ParseDeploymentID parses 'input' into a DeploymentId
func ParseDeploymentID(input string) (*DeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeploymentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDeploymentIDInsensitively parses 'input' case-insensitively into a DeploymentId
// note: this method should only be used for API response data and not user input
func ParseDeploymentIDInsensitively(input string) (*DeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeploymentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DeploymentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SpringName, ok = input.Parsed["springName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "springName", input)
	}

	if id.AppName, ok = input.Parsed["appName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "appName", input)
	}

	if id.DeploymentName, ok = input.Parsed["deploymentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deploymentName", input)
	}

	return nil
}

// ValidateDeploymentID checks that 'input' can be parsed as a Deployment ID
func ValidateDeploymentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeploymentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deployment ID
func (id DeploymentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/apps/%s/deployments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.AppName, id.DeploymentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Deployment ID
func (id DeploymentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springName"),
		resourceids.StaticSegment("staticApps", "apps", "apps"),
		resourceids.UserSpecifiedSegment("appName", "appName"),
		resourceids.StaticSegment("staticDeployments", "deployments", "deployments"),
		resourceids.UserSpecifiedSegment("deploymentName", "deploymentName"),
	}
}

// String returns a human-readable description of this Deployment ID
func (id DeploymentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("App Name: %q", id.AppName),
		fmt.Sprintf("Deployment Name: %q", id.DeploymentName),
	}
	return fmt.Sprintf("Deployment (%s)", strings.Join(components, "\n"))
}
