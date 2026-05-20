package softwareupdateconfiguration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SoftwareUpdateConfigurationId{})
}

var _ resourceids.ResourceId = &SoftwareUpdateConfigurationId{}

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
	parser := resourceids.NewParserFromResourceIdType(&SoftwareUpdateConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SoftwareUpdateConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSoftwareUpdateConfigurationIDInsensitively parses 'input' case-insensitively into a SoftwareUpdateConfigurationId
// note: this method should only be used for API response data and not user input
func ParseSoftwareUpdateConfigurationIDInsensitively(input string) (*SoftwareUpdateConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SoftwareUpdateConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SoftwareUpdateConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SoftwareUpdateConfigurationId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.SoftwareUpdateConfigurationName, ok = input.Parsed["softwareUpdateConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "softwareUpdateConfigurationName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticSoftwareUpdateConfigurations", "softwareUpdateConfigurations", "softwareUpdateConfigurations"),
		resourceids.UserSpecifiedSegment("softwareUpdateConfigurationName", "softwareUpdateConfigurationName"),
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
