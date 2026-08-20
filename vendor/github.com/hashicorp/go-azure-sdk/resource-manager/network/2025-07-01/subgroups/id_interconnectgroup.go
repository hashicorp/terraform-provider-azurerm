package subgroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&InterconnectGroupId{})
}

var _ resourceids.ResourceId = &InterconnectGroupId{}

// InterconnectGroupId is a struct representing the Resource ID for a Interconnect Group
type InterconnectGroupId struct {
	SubscriptionId        string
	ResourceGroupName     string
	InterconnectGroupName string
}

// NewInterconnectGroupID returns a new InterconnectGroupId struct
func NewInterconnectGroupID(subscriptionId string, resourceGroupName string, interconnectGroupName string) InterconnectGroupId {
	return InterconnectGroupId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		InterconnectGroupName: interconnectGroupName,
	}
}

// ParseInterconnectGroupID parses 'input' into a InterconnectGroupId
func ParseInterconnectGroupID(input string) (*InterconnectGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InterconnectGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InterconnectGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseInterconnectGroupIDInsensitively parses 'input' case-insensitively into a InterconnectGroupId
// note: this method should only be used for API response data and not user input
func ParseInterconnectGroupIDInsensitively(input string) (*InterconnectGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InterconnectGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InterconnectGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *InterconnectGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.InterconnectGroupName, ok = input.Parsed["interconnectGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "interconnectGroupName", input)
	}

	return nil
}

// ValidateInterconnectGroupID checks that 'input' can be parsed as a Interconnect Group ID
func ValidateInterconnectGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseInterconnectGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Interconnect Group ID
func (id InterconnectGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/interconnectGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.InterconnectGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Interconnect Group ID
func (id InterconnectGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticInterconnectGroups", "interconnectGroups", "interconnectGroups"),
		resourceids.UserSpecifiedSegment("interconnectGroupName", "interconnectGroupName"),
	}
}

// String returns a human-readable description of this Interconnect Group ID
func (id InterconnectGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Interconnect Group Name: %q", id.InterconnectGroupName),
	}
	return fmt.Sprintf("Interconnect Group (%s)", strings.Join(components, "\n"))
}
