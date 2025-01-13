package sim

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SimId{})
}

var _ resourceids.ResourceId = &SimId{}

// SimId is a struct representing the Resource ID for a Sim
type SimId struct {
	SubscriptionId    string
	ResourceGroupName string
	SimGroupName      string
	SimName           string
}

// NewSimID returns a new SimId struct
func NewSimID(subscriptionId string, resourceGroupName string, simGroupName string, simName string) SimId {
	return SimId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SimGroupName:      simGroupName,
		SimName:           simName,
	}
}

// ParseSimID parses 'input' into a SimId
func ParseSimID(input string) (*SimId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SimId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SimId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSimIDInsensitively parses 'input' case-insensitively into a SimId
// note: this method should only be used for API response data and not user input
func ParseSimIDInsensitively(input string) (*SimId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SimId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SimId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SimId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SimGroupName, ok = input.Parsed["simGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "simGroupName", input)
	}

	if id.SimName, ok = input.Parsed["simName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "simName", input)
	}

	return nil
}

// ValidateSimID checks that 'input' can be parsed as a Sim ID
func ValidateSimID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSimID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Sim ID
func (id SimId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MobileNetwork/simGroups/%s/sims/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SimGroupName, id.SimName)
}

// Segments returns a slice of Resource ID Segments which comprise this Sim ID
func (id SimId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMobileNetwork", "Microsoft.MobileNetwork", "Microsoft.MobileNetwork"),
		resourceids.StaticSegment("staticSimGroups", "simGroups", "simGroups"),
		resourceids.UserSpecifiedSegment("simGroupName", "simGroupName"),
		resourceids.StaticSegment("staticSims", "sims", "sims"),
		resourceids.UserSpecifiedSegment("simName", "simName"),
	}
}

// String returns a human-readable description of this Sim ID
func (id SimId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Sim Group Name: %q", id.SimGroupName),
		fmt.Sprintf("Sim Name: %q", id.SimName),
	}
	return fmt.Sprintf("Sim (%s)", strings.Join(components, "\n"))
}
