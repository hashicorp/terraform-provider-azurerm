package virtualmachineruncommands

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RunCommandId{})
}

var _ resourceids.ResourceId = &RunCommandId{}

// RunCommandId is a struct representing the Resource ID for a Run Command
type RunCommandId struct {
	SubscriptionId string
	LocationName   string
	CommandId      string
}

// NewRunCommandID returns a new RunCommandId struct
func NewRunCommandID(subscriptionId string, locationName string, commandId string) RunCommandId {
	return RunCommandId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		CommandId:      commandId,
	}
}

// ParseRunCommandID parses 'input' into a RunCommandId
func ParseRunCommandID(input string) (*RunCommandId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RunCommandId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RunCommandId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRunCommandIDInsensitively parses 'input' case-insensitively into a RunCommandId
// note: this method should only be used for API response data and not user input
func ParseRunCommandIDInsensitively(input string) (*RunCommandId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RunCommandId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RunCommandId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RunCommandId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.CommandId, ok = input.Parsed["commandId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "commandId", input)
	}

	return nil
}

// ValidateRunCommandID checks that 'input' can be parsed as a Run Command ID
func ValidateRunCommandID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRunCommandID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Run Command ID
func (id RunCommandId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Compute/locations/%s/runCommands/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.CommandId)
}

// Segments returns a slice of Resource ID Segments which comprise this Run Command ID
func (id RunCommandId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticRunCommands", "runCommands", "runCommands"),
		resourceids.UserSpecifiedSegment("commandId", "commandId"),
	}
}

// String returns a human-readable description of this Run Command ID
func (id RunCommandId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Command: %q", id.CommandId),
	}
	return fmt.Sprintf("Run Command (%s)", strings.Join(components, "\n"))
}
