package deploymentstacksatsubscription

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DeploymentStackId{})
}

var _ resourceids.ResourceId = &DeploymentStackId{}

// DeploymentStackId is a struct representing the Resource ID for a Deployment Stack
type DeploymentStackId struct {
	SubscriptionId      string
	DeploymentStackName string
}

// NewDeploymentStackID returns a new DeploymentStackId struct
func NewDeploymentStackID(subscriptionId string, deploymentStackName string) DeploymentStackId {
	return DeploymentStackId{
		SubscriptionId:      subscriptionId,
		DeploymentStackName: deploymentStackName,
	}
}

// ParseDeploymentStackID parses 'input' into a DeploymentStackId
func ParseDeploymentStackID(input string) (*DeploymentStackId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeploymentStackId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeploymentStackId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDeploymentStackIDInsensitively parses 'input' case-insensitively into a DeploymentStackId
// note: this method should only be used for API response data and not user input
func ParseDeploymentStackIDInsensitively(input string) (*DeploymentStackId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeploymentStackId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeploymentStackId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DeploymentStackId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.DeploymentStackName, ok = input.Parsed["deploymentStackName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deploymentStackName", input)
	}

	return nil
}

// ValidateDeploymentStackID checks that 'input' can be parsed as a Deployment Stack ID
func ValidateDeploymentStackID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeploymentStackID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deployment Stack ID
func (id DeploymentStackId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Resources/deploymentStacks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.DeploymentStackName)
}

// Segments returns a slice of Resource ID Segments which comprise this Deployment Stack ID
func (id DeploymentStackId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftResources", "Microsoft.Resources", "Microsoft.Resources"),
		resourceids.StaticSegment("staticDeploymentStacks", "deploymentStacks", "deploymentStacks"),
		resourceids.UserSpecifiedSegment("deploymentStackName", "deploymentStackName"),
	}
}

// String returns a human-readable description of this Deployment Stack ID
func (id DeploymentStackId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Deployment Stack Name: %q", id.DeploymentStackName),
	}
	return fmt.Sprintf("Deployment Stack (%s)", strings.Join(components, "\n"))
}
