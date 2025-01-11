package integrationserviceenvironments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&IntegrationServiceEnvironmentId{})
}

var _ resourceids.ResourceId = &IntegrationServiceEnvironmentId{}

// IntegrationServiceEnvironmentId is a struct representing the Resource ID for a Integration Service Environment
type IntegrationServiceEnvironmentId struct {
	SubscriptionId                    string
	ResourceGroup                     string
	IntegrationServiceEnvironmentName string
}

// NewIntegrationServiceEnvironmentID returns a new IntegrationServiceEnvironmentId struct
func NewIntegrationServiceEnvironmentID(subscriptionId string, resourceGroup string, integrationServiceEnvironmentName string) IntegrationServiceEnvironmentId {
	return IntegrationServiceEnvironmentId{
		SubscriptionId:                    subscriptionId,
		ResourceGroup:                     resourceGroup,
		IntegrationServiceEnvironmentName: integrationServiceEnvironmentName,
	}
}

// ParseIntegrationServiceEnvironmentID parses 'input' into a IntegrationServiceEnvironmentId
func ParseIntegrationServiceEnvironmentID(input string) (*IntegrationServiceEnvironmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IntegrationServiceEnvironmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IntegrationServiceEnvironmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseIntegrationServiceEnvironmentIDInsensitively parses 'input' case-insensitively into a IntegrationServiceEnvironmentId
// note: this method should only be used for API response data and not user input
func ParseIntegrationServiceEnvironmentIDInsensitively(input string) (*IntegrationServiceEnvironmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IntegrationServiceEnvironmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IntegrationServiceEnvironmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *IntegrationServiceEnvironmentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroup, ok = input.Parsed["resourceGroup"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroup", input)
	}

	if id.IntegrationServiceEnvironmentName, ok = input.Parsed["integrationServiceEnvironmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "integrationServiceEnvironmentName", input)
	}

	return nil
}

// ValidateIntegrationServiceEnvironmentID checks that 'input' can be parsed as a Integration Service Environment ID
func ValidateIntegrationServiceEnvironmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIntegrationServiceEnvironmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Integration Service Environment ID
func (id IntegrationServiceEnvironmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationServiceEnvironments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IntegrationServiceEnvironmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Integration Service Environment ID
func (id IntegrationServiceEnvironmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroup", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLogic", "Microsoft.Logic", "Microsoft.Logic"),
		resourceids.StaticSegment("staticIntegrationServiceEnvironments", "integrationServiceEnvironments", "integrationServiceEnvironments"),
		resourceids.UserSpecifiedSegment("integrationServiceEnvironmentName", "integrationServiceEnvironmentName"),
	}
}

// String returns a human-readable description of this Integration Service Environment ID
func (id IntegrationServiceEnvironmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group: %q", id.ResourceGroup),
		fmt.Sprintf("Integration Service Environment Name: %q", id.IntegrationServiceEnvironmentName),
	}
	return fmt.Sprintf("Integration Service Environment (%s)", strings.Join(components, "\n"))
}
