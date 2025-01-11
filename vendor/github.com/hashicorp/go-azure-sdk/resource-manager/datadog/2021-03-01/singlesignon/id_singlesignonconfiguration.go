package singlesignon

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SingleSignOnConfigurationId{})
}

var _ resourceids.ResourceId = &SingleSignOnConfigurationId{}

// SingleSignOnConfigurationId is a struct representing the Resource ID for a Single Sign On Configuration
type SingleSignOnConfigurationId struct {
	SubscriptionId                string
	ResourceGroupName             string
	MonitorName                   string
	SingleSignOnConfigurationName string
}

// NewSingleSignOnConfigurationID returns a new SingleSignOnConfigurationId struct
func NewSingleSignOnConfigurationID(subscriptionId string, resourceGroupName string, monitorName string, singleSignOnConfigurationName string) SingleSignOnConfigurationId {
	return SingleSignOnConfigurationId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		MonitorName:                   monitorName,
		SingleSignOnConfigurationName: singleSignOnConfigurationName,
	}
}

// ParseSingleSignOnConfigurationID parses 'input' into a SingleSignOnConfigurationId
func ParseSingleSignOnConfigurationID(input string) (*SingleSignOnConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SingleSignOnConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SingleSignOnConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSingleSignOnConfigurationIDInsensitively parses 'input' case-insensitively into a SingleSignOnConfigurationId
// note: this method should only be used for API response data and not user input
func ParseSingleSignOnConfigurationIDInsensitively(input string) (*SingleSignOnConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SingleSignOnConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SingleSignOnConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SingleSignOnConfigurationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MonitorName, ok = input.Parsed["monitorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "monitorName", input)
	}

	if id.SingleSignOnConfigurationName, ok = input.Parsed["singleSignOnConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "singleSignOnConfigurationName", input)
	}

	return nil
}

// ValidateSingleSignOnConfigurationID checks that 'input' can be parsed as a Single Sign On Configuration ID
func ValidateSingleSignOnConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSingleSignOnConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Single Sign On Configuration ID
func (id SingleSignOnConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Datadog/monitors/%s/singleSignOnConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MonitorName, id.SingleSignOnConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Single Sign On Configuration ID
func (id SingleSignOnConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDatadog", "Microsoft.Datadog", "Microsoft.Datadog"),
		resourceids.StaticSegment("staticMonitors", "monitors", "monitors"),
		resourceids.UserSpecifiedSegment("monitorName", "monitorName"),
		resourceids.StaticSegment("staticSingleSignOnConfigurations", "singleSignOnConfigurations", "singleSignOnConfigurations"),
		resourceids.UserSpecifiedSegment("singleSignOnConfigurationName", "singleSignOnConfigurationName"),
	}
}

// String returns a human-readable description of this Single Sign On Configuration ID
func (id SingleSignOnConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Monitor Name: %q", id.MonitorName),
		fmt.Sprintf("Single Sign On Configuration Name: %q", id.SingleSignOnConfigurationName),
	}
	return fmt.Sprintf("Single Sign On Configuration (%s)", strings.Join(components, "\n"))
}
