package configurationassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConfigurationAssignmentId{})
}

var _ resourceids.ResourceId = &ConfigurationAssignmentId{}

// ConfigurationAssignmentId is a struct representing the Resource ID for a Configuration Assignment
type ConfigurationAssignmentId struct {
	SubscriptionId              string
	ConfigurationAssignmentName string
}

// NewConfigurationAssignmentID returns a new ConfigurationAssignmentId struct
func NewConfigurationAssignmentID(subscriptionId string, configurationAssignmentName string) ConfigurationAssignmentId {
	return ConfigurationAssignmentId{
		SubscriptionId:              subscriptionId,
		ConfigurationAssignmentName: configurationAssignmentName,
	}
}

// ParseConfigurationAssignmentID parses 'input' into a ConfigurationAssignmentId
func ParseConfigurationAssignmentID(input string) (*ConfigurationAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConfigurationAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigurationAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConfigurationAssignmentIDInsensitively parses 'input' case-insensitively into a ConfigurationAssignmentId
// note: this method should only be used for API response data and not user input
func ParseConfigurationAssignmentIDInsensitively(input string) (*ConfigurationAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConfigurationAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigurationAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConfigurationAssignmentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ConfigurationAssignmentName, ok = input.Parsed["configurationAssignmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "configurationAssignmentName", input)
	}

	return nil
}

// ValidateConfigurationAssignmentID checks that 'input' can be parsed as a Configuration Assignment ID
func ValidateConfigurationAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConfigurationAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Configuration Assignment ID
func (id ConfigurationAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Maintenance/configurationAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ConfigurationAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Configuration Assignment ID
func (id ConfigurationAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMaintenance", "Microsoft.Maintenance", "Microsoft.Maintenance"),
		resourceids.StaticSegment("staticConfigurationAssignments", "configurationAssignments", "configurationAssignments"),
		resourceids.UserSpecifiedSegment("configurationAssignmentName", "configurationAssignmentName"),
	}
}

// String returns a human-readable description of this Configuration Assignment ID
func (id ConfigurationAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Configuration Assignment Name: %q", id.ConfigurationAssignmentName),
	}
	return fmt.Sprintf("Configuration Assignment (%s)", strings.Join(components, "\n"))
}
