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
	recaser.RegisterResourceId(&SubgroupId{})
}

var _ resourceids.ResourceId = &SubgroupId{}

// SubgroupId is a struct representing the Resource ID for a Subgroup
type SubgroupId struct {
	SubscriptionId        string
	ResourceGroupName     string
	InterconnectGroupName string
	SubgroupName          string
}

// NewSubgroupID returns a new SubgroupId struct
func NewSubgroupID(subscriptionId string, resourceGroupName string, interconnectGroupName string, subgroupName string) SubgroupId {
	return SubgroupId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		InterconnectGroupName: interconnectGroupName,
		SubgroupName:          subgroupName,
	}
}

// ParseSubgroupID parses 'input' into a SubgroupId
func ParseSubgroupID(input string) (*SubgroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SubgroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SubgroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSubgroupIDInsensitively parses 'input' case-insensitively into a SubgroupId
// note: this method should only be used for API response data and not user input
func ParseSubgroupIDInsensitively(input string) (*SubgroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SubgroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SubgroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SubgroupId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.SubgroupName, ok = input.Parsed["subgroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subgroupName", input)
	}

	return nil
}

// ValidateSubgroupID checks that 'input' can be parsed as a Subgroup ID
func ValidateSubgroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSubgroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Subgroup ID
func (id SubgroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/interconnectGroups/%s/subgroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.InterconnectGroupName, id.SubgroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Subgroup ID
func (id SubgroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticInterconnectGroups", "interconnectGroups", "interconnectGroups"),
		resourceids.UserSpecifiedSegment("interconnectGroupName", "interconnectGroupName"),
		resourceids.StaticSegment("staticSubgroups", "subgroups", "subgroups"),
		resourceids.UserSpecifiedSegment("subgroupName", "subgroupName"),
	}
}

// String returns a human-readable description of this Subgroup ID
func (id SubgroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Interconnect Group Name: %q", id.InterconnectGroupName),
		fmt.Sprintf("Subgroup Name: %q", id.SubgroupName),
	}
	return fmt.Sprintf("Subgroup (%s)", strings.Join(components, "\n"))
}
