package machineruncommands

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
	SubscriptionId    string
	ResourceGroupName string
	MachineName       string
	RunCommandName    string
}

// NewRunCommandID returns a new RunCommandId struct
func NewRunCommandID(subscriptionId string, resourceGroupName string, machineName string, runCommandName string) RunCommandId {
	return RunCommandId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MachineName:       machineName,
		RunCommandName:    runCommandName,
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
	if err := id.FromParseResult(*parsed); err != nil {
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
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RunCommandId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MachineName, ok = input.Parsed["machineName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "machineName", input)
	}

	if id.RunCommandName, ok = input.Parsed["runCommandName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "runCommandName", input)
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HybridCompute/machines/%s/runCommands/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MachineName, id.RunCommandName)
}

// Segments returns a slice of Resource ID Segments which comprise this Run Command ID
func (id RunCommandId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHybridCompute", "Microsoft.HybridCompute", "Microsoft.HybridCompute"),
		resourceids.StaticSegment("staticMachines", "machines", "machines"),
		resourceids.UserSpecifiedSegment("machineName", "machineValue"),
		resourceids.StaticSegment("staticRunCommands", "runCommands", "runCommands"),
		resourceids.UserSpecifiedSegment("runCommandName", "runCommandValue"),
	}
}

// String returns a human-readable description of this Run Command ID
func (id RunCommandId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Machine Name: %q", id.MachineName),
		fmt.Sprintf("Run Command Name: %q", id.RunCommandName),
	}
	return fmt.Sprintf("Run Command (%s)", strings.Join(components, "\n"))
}
