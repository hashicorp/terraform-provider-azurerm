package namedvalue

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NamedValueId{})
}

var _ resourceids.ResourceId = &NamedValueId{}

// NamedValueId is a struct representing the Resource ID for a Named Value
type NamedValueId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	NamedValueId      string
}

// NewNamedValueID returns a new NamedValueId struct
func NewNamedValueID(subscriptionId string, resourceGroupName string, serviceName string, namedValueId string) NamedValueId {
	return NamedValueId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		NamedValueId:      namedValueId,
	}
}

// ParseNamedValueID parses 'input' into a NamedValueId
func ParseNamedValueID(input string) (*NamedValueId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NamedValueId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NamedValueId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNamedValueIDInsensitively parses 'input' case-insensitively into a NamedValueId
// note: this method should only be used for API response data and not user input
func ParseNamedValueIDInsensitively(input string) (*NamedValueId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NamedValueId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NamedValueId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NamedValueId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.NamedValueId, ok = input.Parsed["namedValueId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "namedValueId", input)
	}

	return nil
}

// ValidateNamedValueID checks that 'input' can be parsed as a Named Value ID
func ValidateNamedValueID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNamedValueID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Named Value ID
func (id NamedValueId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/namedValues/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.NamedValueId)
}

// Segments returns a slice of Resource ID Segments which comprise this Named Value ID
func (id NamedValueId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticNamedValues", "namedValues", "namedValues"),
		resourceids.UserSpecifiedSegment("namedValueId", "namedValueId"),
	}
}

// String returns a human-readable description of this Named Value ID
func (id NamedValueId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Named Value: %q", id.NamedValueId),
	}
	return fmt.Sprintf("Named Value (%s)", strings.Join(components, "\n"))
}
