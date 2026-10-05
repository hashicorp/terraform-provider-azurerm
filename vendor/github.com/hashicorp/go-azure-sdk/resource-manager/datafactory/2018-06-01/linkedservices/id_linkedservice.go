package linkedservices

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LinkedServiceId{})
}

var _ resourceids.ResourceId = &LinkedServiceId{}

// LinkedServiceId is a struct representing the Resource ID for a Linked Service
type LinkedServiceId struct {
	SubscriptionId    string
	ResourceGroupName string
	FactoryName       string
	LinkedServiceName string
}

// NewLinkedServiceID returns a new LinkedServiceId struct
func NewLinkedServiceID(subscriptionId string, resourceGroupName string, factoryName string, linkedServiceName string) LinkedServiceId {
	return LinkedServiceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		FactoryName:       factoryName,
		LinkedServiceName: linkedServiceName,
	}
}

// ParseLinkedServiceID parses 'input' into a LinkedServiceId
func ParseLinkedServiceID(input string) (*LinkedServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LinkedServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LinkedServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLinkedServiceIDInsensitively parses 'input' case-insensitively into a LinkedServiceId
// note: this method should only be used for API response data and not user input
func ParseLinkedServiceIDInsensitively(input string) (*LinkedServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LinkedServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LinkedServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LinkedServiceId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.LinkedServiceName, ok = input.Parsed["linkedServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "linkedServiceName", input)
	}

	return nil
}

// ValidateLinkedServiceID checks that 'input' can be parsed as a Linked Service ID
func ValidateLinkedServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLinkedServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Linked Service ID
func (id LinkedServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataFactory/factories/%s/linkedServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FactoryName, id.LinkedServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Linked Service ID
func (id LinkedServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataFactory", "Microsoft.DataFactory", "Microsoft.DataFactory"),
		resourceids.StaticSegment("staticFactories", "factories", "factories"),
		resourceids.UserSpecifiedSegment("factoryName", "factoryName"),
		resourceids.StaticSegment("staticLinkedServices", "linkedServices", "linkedServices"),
		resourceids.UserSpecifiedSegment("linkedServiceName", "linkedServiceName"),
	}
}

// String returns a human-readable description of this Linked Service ID
func (id LinkedServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Factory Name: %q", id.FactoryName),
		fmt.Sprintf("Linked Service Name: %q", id.LinkedServiceName),
	}
	return fmt.Sprintf("Linked Service (%s)", strings.Join(components, "\n"))
}
