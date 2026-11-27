package nodepools

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SupercomputerId{})
}

var _ resourceids.ResourceId = &SupercomputerId{}

// SupercomputerId is a struct representing the Resource ID for a Supercomputer
type SupercomputerId struct {
	SubscriptionId    string
	ResourceGroupName string
	SupercomputerName string
}

// NewSupercomputerID returns a new SupercomputerId struct
func NewSupercomputerID(subscriptionId string, resourceGroupName string, supercomputerName string) SupercomputerId {
	return SupercomputerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SupercomputerName: supercomputerName,
	}
}

// ParseSupercomputerID parses 'input' into a SupercomputerId
func ParseSupercomputerID(input string) (*SupercomputerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SupercomputerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SupercomputerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSupercomputerIDInsensitively parses 'input' case-insensitively into a SupercomputerId
// note: this method should only be used for API response data and not user input
func ParseSupercomputerIDInsensitively(input string) (*SupercomputerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SupercomputerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SupercomputerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SupercomputerId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SupercomputerName, ok = input.Parsed["supercomputerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "supercomputerName", input)
	}

	return nil
}

// ValidateSupercomputerID checks that 'input' can be parsed as a Supercomputer ID
func ValidateSupercomputerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSupercomputerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Supercomputer ID
func (id SupercomputerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Discovery/supercomputers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SupercomputerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Supercomputer ID
func (id SupercomputerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDiscovery", "Microsoft.Discovery", "Microsoft.Discovery"),
		resourceids.StaticSegment("staticSupercomputers", "supercomputers", "supercomputers"),
		resourceids.UserSpecifiedSegment("supercomputerName", "supercomputerName"),
	}
}

// String returns a human-readable description of this Supercomputer ID
func (id SupercomputerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Supercomputer Name: %q", id.SupercomputerName),
	}
	return fmt.Sprintf("Supercomputer (%s)", strings.Join(components, "\n"))
}
