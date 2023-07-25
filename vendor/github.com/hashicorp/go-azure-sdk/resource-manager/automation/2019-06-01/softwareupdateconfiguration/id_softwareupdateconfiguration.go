package softwareupdateconfiguration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SoftwareUpdateConfigurationId{}

// SoftwareUpdateConfigurationId is a struct representing the Resource ID for a Software Update Configuration
type SoftwareUpdateConfigurationId struct {
	SubscriptionId                  string
	ResourceGroupName               string
	AutomationAccountName           string
	SoftwareUpdateConfigurationName string
}

// NewSoftwareUpdateConfigurationID returns a new SoftwareUpdateConfigurationId struct
func NewSoftwareUpdateConfigurationID(subscriptionId string, resourceGroupName string, automationAccountName string, softwareUpdateConfigurationName string) SoftwareUpdateConfigurationId {
	return SoftwareUpdateConfigurationId{
		SubscriptionId:                  subscriptionId,
		ResourceGroupName:               resourceGroupName,
		AutomationAccountName:           automationAccountName,
		SoftwareUpdateConfigurationName: softwareUpdateConfigurationName,
	}
}

// ParseSoftwareUpdateConfigurationID parses 'input' into a SoftwareUpdateConfigurationId
func ParseSoftwareUpdateConfigurationID(input string) (*SoftwareUpdateConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(SoftwareUpdateConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SoftwareUpdateConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.SoftwareUpdateConfigurationName, ok = parsed.Parsed["softwareUpdateConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "softwareUpdateConfigurationName", *parsed)
	}

	return &id, nil
}

// ParseSoftwareUpdateConfigurationIDInsensitively parses 'input' case-insensitively into a SoftwareUpdateConfigurationId
// note: this method should only be used for API response data and not user input
func ParseSoftwareUpdateConfigurationIDInsensitively(input string) (*SoftwareUpdateConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(SoftwareUpdateConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SoftwareUpdateConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.SoftwareUpdateConfigurationName, ok = parsed.Parsed["softwareUpdateConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "softwareUpdateConfigurationName", *parsed)
	}

	return &id, nil
}

// ValidateSoftwareUpdateConfigurationID checks that 'input' can be parsed as a Software Update Configuration ID
func ValidateSoftwareUpdateConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSoftwareUpdateConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Software Update Configuration ID
func (id SoftwareUpdateConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/softwareUpdateConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.SoftwareUpdateConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Software Update Configuration ID
func (id SoftwareUpdateConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountValue"),
		resourceids.StaticSegment("staticSoftwareUpdateConfigurations", "softwareUpdateConfigurations", "softwareUpdateConfigurations"),
		resourceids.UserSpecifiedSegment("softwareUpdateConfigurationName", "softwareUpdateConfigurationValue"),
	}
}

// String returns a human-readable description of this Software Update Configuration ID
func (id SoftwareUpdateConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Software Update Configuration Name: %q", id.SoftwareUpdateConfigurationName),
	}
	return fmt.Sprintf("Software Update Configuration (%s)", strings.Join(components, "\n"))
}
