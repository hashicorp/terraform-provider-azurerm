package module

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PowerShell72ModuleId{})
}

var _ resourceids.ResourceId = &PowerShell72ModuleId{}

// PowerShell72ModuleId is a struct representing the Resource ID for a Power Shell 7 2 Module
type PowerShell72ModuleId struct {
	SubscriptionId         string
	ResourceGroupName      string
	AutomationAccountName  string
	PowerShell72ModuleName string
}

// NewPowerShell72ModuleID returns a new PowerShell72ModuleId struct
func NewPowerShell72ModuleID(subscriptionId string, resourceGroupName string, automationAccountName string, powerShell72ModuleName string) PowerShell72ModuleId {
	return PowerShell72ModuleId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		AutomationAccountName:  automationAccountName,
		PowerShell72ModuleName: powerShell72ModuleName,
	}
}

// ParsePowerShell72ModuleID parses 'input' into a PowerShell72ModuleId
func ParsePowerShell72ModuleID(input string) (*PowerShell72ModuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PowerShell72ModuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PowerShell72ModuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePowerShell72ModuleIDInsensitively parses 'input' case-insensitively into a PowerShell72ModuleId
// note: this method should only be used for API response data and not user input
func ParsePowerShell72ModuleIDInsensitively(input string) (*PowerShell72ModuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PowerShell72ModuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PowerShell72ModuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PowerShell72ModuleId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.PowerShell72ModuleName, ok = input.Parsed["powerShell72ModuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "powerShell72ModuleName", input)
	}

	return nil
}

// ValidatePowerShell72ModuleID checks that 'input' can be parsed as a Power Shell 7 2 Module ID
func ValidatePowerShell72ModuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePowerShell72ModuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Power Shell 7 2 Module ID
func (id PowerShell72ModuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/powerShell72Modules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.PowerShell72ModuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Power Shell 7 2 Module ID
func (id PowerShell72ModuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticPowerShell72Modules", "powerShell72Modules", "powerShell72Modules"),
		resourceids.UserSpecifiedSegment("powerShell72ModuleName", "powerShell72ModuleName"),
	}
}

// String returns a human-readable description of this Power Shell 7 2 Module ID
func (id PowerShell72ModuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Power Shell 7 2 Module Name: %q", id.PowerShell72ModuleName),
	}
	return fmt.Sprintf("Power Shell 7 2 Module (%s)", strings.Join(components, "\n"))
}
