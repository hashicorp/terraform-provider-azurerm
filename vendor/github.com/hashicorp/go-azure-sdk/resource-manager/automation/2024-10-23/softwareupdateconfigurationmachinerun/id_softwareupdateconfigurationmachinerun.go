package softwareupdateconfigurationmachinerun

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SoftwareUpdateConfigurationMachineRunId{})
}

var _ resourceids.ResourceId = &SoftwareUpdateConfigurationMachineRunId{}

// SoftwareUpdateConfigurationMachineRunId is a struct representing the Resource ID for a Software Update Configuration Machine Run
type SoftwareUpdateConfigurationMachineRunId struct {
	SubscriptionId                          string
	ResourceGroupName                       string
	AutomationAccountName                   string
	SoftwareUpdateConfigurationMachineRunId string
}

// NewSoftwareUpdateConfigurationMachineRunID returns a new SoftwareUpdateConfigurationMachineRunId struct
func NewSoftwareUpdateConfigurationMachineRunID(subscriptionId string, resourceGroupName string, automationAccountName string, softwareUpdateConfigurationMachineRunId string) SoftwareUpdateConfigurationMachineRunId {
	return SoftwareUpdateConfigurationMachineRunId{
		SubscriptionId:                          subscriptionId,
		ResourceGroupName:                       resourceGroupName,
		AutomationAccountName:                   automationAccountName,
		SoftwareUpdateConfigurationMachineRunId: softwareUpdateConfigurationMachineRunId,
	}
}

// ParseSoftwareUpdateConfigurationMachineRunID parses 'input' into a SoftwareUpdateConfigurationMachineRunId
func ParseSoftwareUpdateConfigurationMachineRunID(input string) (*SoftwareUpdateConfigurationMachineRunId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SoftwareUpdateConfigurationMachineRunId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SoftwareUpdateConfigurationMachineRunId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSoftwareUpdateConfigurationMachineRunIDInsensitively parses 'input' case-insensitively into a SoftwareUpdateConfigurationMachineRunId
// note: this method should only be used for API response data and not user input
func ParseSoftwareUpdateConfigurationMachineRunIDInsensitively(input string) (*SoftwareUpdateConfigurationMachineRunId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SoftwareUpdateConfigurationMachineRunId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SoftwareUpdateConfigurationMachineRunId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SoftwareUpdateConfigurationMachineRunId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AutomationAccountName, ok = input.Parsed["automationAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", input)
	}

	if id.SoftwareUpdateConfigurationMachineRunId, ok = input.Parsed["softwareUpdateConfigurationMachineRunId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "softwareUpdateConfigurationMachineRunId", input)
	}

	return nil
}

// ValidateSoftwareUpdateConfigurationMachineRunID checks that 'input' can be parsed as a Software Update Configuration Machine Run ID
func ValidateSoftwareUpdateConfigurationMachineRunID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSoftwareUpdateConfigurationMachineRunID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Software Update Configuration Machine Run ID
func (id SoftwareUpdateConfigurationMachineRunId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/softwareUpdateConfigurationMachineRuns/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.SoftwareUpdateConfigurationMachineRunId)
}

// Segments returns a slice of Resource ID Segments which comprise this Software Update Configuration Machine Run ID
func (id SoftwareUpdateConfigurationMachineRunId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticSoftwareUpdateConfigurationMachineRuns", "softwareUpdateConfigurationMachineRuns", "softwareUpdateConfigurationMachineRuns"),
		resourceids.UserSpecifiedSegment("softwareUpdateConfigurationMachineRunId", "softwareUpdateConfigurationMachineRunId"),
	}
}

// String returns a human-readable description of this Software Update Configuration Machine Run ID
func (id SoftwareUpdateConfigurationMachineRunId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Software Update Configuration Machine Run: %q", id.SoftwareUpdateConfigurationMachineRunId),
	}
	return fmt.Sprintf("Software Update Configuration Machine Run (%s)", strings.Join(components, "\n"))
}
