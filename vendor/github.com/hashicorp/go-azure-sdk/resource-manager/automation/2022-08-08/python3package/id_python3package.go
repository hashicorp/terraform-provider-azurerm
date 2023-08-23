package python3package

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = Python3PackageId{}

// Python3PackageId is a struct representing the Resource ID for a Python 3 Package
type Python3PackageId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	Python3PackageName    string
}

// NewPython3PackageID returns a new Python3PackageId struct
func NewPython3PackageID(subscriptionId string, resourceGroupName string, automationAccountName string, python3PackageName string) Python3PackageId {
	return Python3PackageId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		Python3PackageName:    python3PackageName,
	}
}

// ParsePython3PackageID parses 'input' into a Python3PackageId
func ParsePython3PackageID(input string) (*Python3PackageId, error) {
	parser := resourceids.NewParserFromResourceIdType(Python3PackageId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := Python3PackageId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.Python3PackageName, ok = parsed.Parsed["python3PackageName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "python3PackageName", *parsed)
	}

	return &id, nil
}

// ParsePython3PackageIDInsensitively parses 'input' case-insensitively into a Python3PackageId
// note: this method should only be used for API response data and not user input
func ParsePython3PackageIDInsensitively(input string) (*Python3PackageId, error) {
	parser := resourceids.NewParserFromResourceIdType(Python3PackageId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := Python3PackageId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.Python3PackageName, ok = parsed.Parsed["python3PackageName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "python3PackageName", *parsed)
	}

	return &id, nil
}

// ValidatePython3PackageID checks that 'input' can be parsed as a Python 3 Package ID
func ValidatePython3PackageID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePython3PackageID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Python 3 Package ID
func (id Python3PackageId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/python3Packages/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.Python3PackageName)
}

// Segments returns a slice of Resource ID Segments which comprise this Python 3 Package ID
func (id Python3PackageId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountValue"),
		resourceids.StaticSegment("staticPython3Packages", "python3Packages", "python3Packages"),
		resourceids.UserSpecifiedSegment("python3PackageName", "python3PackageValue"),
	}
}

// String returns a human-readable description of this Python 3 Package ID
func (id Python3PackageId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Python 3 Package Name: %q", id.Python3PackageName),
	}
	return fmt.Sprintf("Python 3 Package (%s)", strings.Join(components, "\n"))
}
