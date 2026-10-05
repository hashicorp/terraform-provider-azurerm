package logicalnetworks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LogicalNetworkId{})
}

var _ resourceids.ResourceId = &LogicalNetworkId{}

// LogicalNetworkId is a struct representing the Resource ID for a Logical Network
type LogicalNetworkId struct {
	SubscriptionId     string
	ResourceGroupName  string
	LogicalNetworkName string
}

// NewLogicalNetworkID returns a new LogicalNetworkId struct
func NewLogicalNetworkID(subscriptionId string, resourceGroupName string, logicalNetworkName string) LogicalNetworkId {
	return LogicalNetworkId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		LogicalNetworkName: logicalNetworkName,
	}
}

// ParseLogicalNetworkID parses 'input' into a LogicalNetworkId
func ParseLogicalNetworkID(input string) (*LogicalNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LogicalNetworkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LogicalNetworkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLogicalNetworkIDInsensitively parses 'input' case-insensitively into a LogicalNetworkId
// note: this method should only be used for API response data and not user input
func ParseLogicalNetworkIDInsensitively(input string) (*LogicalNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LogicalNetworkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LogicalNetworkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LogicalNetworkId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.LogicalNetworkName, ok = input.Parsed["logicalNetworkName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "logicalNetworkName", input)
	}

	return nil
}

// ValidateLogicalNetworkID checks that 'input' can be parsed as a Logical Network ID
func ValidateLogicalNetworkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLogicalNetworkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Logical Network ID
func (id LogicalNetworkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHCI/logicalNetworks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LogicalNetworkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Logical Network ID
func (id LogicalNetworkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAzureStackHCI", "Microsoft.AzureStackHCI", "Microsoft.AzureStackHCI"),
		resourceids.StaticSegment("staticLogicalNetworks", "logicalNetworks", "logicalNetworks"),
		resourceids.UserSpecifiedSegment("logicalNetworkName", "logicalNetworkName"),
	}
}

// String returns a human-readable description of this Logical Network ID
func (id LogicalNetworkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Logical Network Name: %q", id.LogicalNetworkName),
	}
	return fmt.Sprintf("Logical Network (%s)", strings.Join(components, "\n"))
}
