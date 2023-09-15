package softwareupdateconfigurationrun

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SoftwareUpdateConfigurationRunId{}

// SoftwareUpdateConfigurationRunId is a struct representing the Resource ID for a Software Update Configuration Run
type SoftwareUpdateConfigurationRunId struct {
	SubscriptionId                   string
	ResourceGroupName                string
	AutomationAccountName            string
	SoftwareUpdateConfigurationRunId string
}

// NewSoftwareUpdateConfigurationRunID returns a new SoftwareUpdateConfigurationRunId struct
func NewSoftwareUpdateConfigurationRunID(subscriptionId string, resourceGroupName string, automationAccountName string, softwareUpdateConfigurationRunId string) SoftwareUpdateConfigurationRunId {
	return SoftwareUpdateConfigurationRunId{
		SubscriptionId:                   subscriptionId,
		ResourceGroupName:                resourceGroupName,
		AutomationAccountName:            automationAccountName,
		SoftwareUpdateConfigurationRunId: softwareUpdateConfigurationRunId,
	}
}

// ParseSoftwareUpdateConfigurationRunID parses 'input' into a SoftwareUpdateConfigurationRunId
func ParseSoftwareUpdateConfigurationRunID(input string) (*SoftwareUpdateConfigurationRunId, error) {
	parser := resourceids.NewParserFromResourceIdType(SoftwareUpdateConfigurationRunId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SoftwareUpdateConfigurationRunId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.SoftwareUpdateConfigurationRunId, ok = parsed.Parsed["softwareUpdateConfigurationRunId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "softwareUpdateConfigurationRunId", *parsed)
	}

	return &id, nil
}

// ParseSoftwareUpdateConfigurationRunIDInsensitively parses 'input' case-insensitively into a SoftwareUpdateConfigurationRunId
// note: this method should only be used for API response data and not user input
func ParseSoftwareUpdateConfigurationRunIDInsensitively(input string) (*SoftwareUpdateConfigurationRunId, error) {
	parser := resourceids.NewParserFromResourceIdType(SoftwareUpdateConfigurationRunId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SoftwareUpdateConfigurationRunId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.SoftwareUpdateConfigurationRunId, ok = parsed.Parsed["softwareUpdateConfigurationRunId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "softwareUpdateConfigurationRunId", *parsed)
	}

	return &id, nil
}

// ValidateSoftwareUpdateConfigurationRunID checks that 'input' can be parsed as a Software Update Configuration Run ID
func ValidateSoftwareUpdateConfigurationRunID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSoftwareUpdateConfigurationRunID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Software Update Configuration Run ID
func (id SoftwareUpdateConfigurationRunId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/softwareUpdateConfigurationRuns/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.SoftwareUpdateConfigurationRunId)
}

// Segments returns a slice of Resource ID Segments which comprise this Software Update Configuration Run ID
func (id SoftwareUpdateConfigurationRunId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountValue"),
		resourceids.StaticSegment("staticSoftwareUpdateConfigurationRuns", "softwareUpdateConfigurationRuns", "softwareUpdateConfigurationRuns"),
		resourceids.UserSpecifiedSegment("softwareUpdateConfigurationRunId", "softwareUpdateConfigurationRunIdValue"),
	}
}

// String returns a human-readable description of this Software Update Configuration Run ID
func (id SoftwareUpdateConfigurationRunId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Software Update Configuration Run: %q", id.SoftwareUpdateConfigurationRunId),
	}
	return fmt.Sprintf("Software Update Configuration Run (%s)", strings.Join(components, "\n"))
}
