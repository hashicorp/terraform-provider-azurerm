package integrationruntimeenableinteractivequery

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&IntegrationRuntimeId{})
}

var _ resourceids.ResourceId = &IntegrationRuntimeId{}

// IntegrationRuntimeId is a struct representing the Resource ID for a Integration Runtime
type IntegrationRuntimeId struct {
	SubscriptionId         string
	ResourceGroupName      string
	FactoryName            string
	IntegrationRuntimeName string
}

// NewIntegrationRuntimeID returns a new IntegrationRuntimeId struct
func NewIntegrationRuntimeID(subscriptionId string, resourceGroupName string, factoryName string, integrationRuntimeName string) IntegrationRuntimeId {
	return IntegrationRuntimeId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		FactoryName:            factoryName,
		IntegrationRuntimeName: integrationRuntimeName,
	}
}

// ParseIntegrationRuntimeID parses 'input' into a IntegrationRuntimeId
func ParseIntegrationRuntimeID(input string) (*IntegrationRuntimeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IntegrationRuntimeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IntegrationRuntimeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseIntegrationRuntimeIDInsensitively parses 'input' case-insensitively into a IntegrationRuntimeId
// note: this method should only be used for API response data and not user input
func ParseIntegrationRuntimeIDInsensitively(input string) (*IntegrationRuntimeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IntegrationRuntimeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IntegrationRuntimeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *IntegrationRuntimeId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FactoryName, ok = input.Parsed["factoryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "factoryName", input)
	}

	if id.IntegrationRuntimeName, ok = input.Parsed["integrationRuntimeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "integrationRuntimeName", input)
	}

	return nil
}

// ValidateIntegrationRuntimeID checks that 'input' can be parsed as a Integration Runtime ID
func ValidateIntegrationRuntimeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIntegrationRuntimeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Integration Runtime ID
func (id IntegrationRuntimeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataFactory/factories/%s/integrationRuntimes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FactoryName, id.IntegrationRuntimeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Integration Runtime ID
func (id IntegrationRuntimeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataFactory", "Microsoft.DataFactory", "Microsoft.DataFactory"),
		resourceids.StaticSegment("staticFactories", "factories", "factories"),
		resourceids.UserSpecifiedSegment("factoryName", "factoryName"),
		resourceids.StaticSegment("staticIntegrationRuntimes", "integrationRuntimes", "integrationRuntimes"),
		resourceids.UserSpecifiedSegment("integrationRuntimeName", "integrationRuntimeName"),
	}
}

// String returns a human-readable description of this Integration Runtime ID
func (id IntegrationRuntimeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Factory Name: %q", id.FactoryName),
		fmt.Sprintf("Integration Runtime Name: %q", id.IntegrationRuntimeName),
	}
	return fmt.Sprintf("Integration Runtime (%s)", strings.Join(components, "\n"))
}
