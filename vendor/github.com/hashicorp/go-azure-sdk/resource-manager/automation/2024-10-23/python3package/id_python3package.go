package python3package

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&Python3PackageId{})
}

var _ resourceids.ResourceId = &Python3PackageId{}

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
	parser := resourceids.NewParserFromResourceIdType(&Python3PackageId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Python3PackageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePython3PackageIDInsensitively parses 'input' case-insensitively into a Python3PackageId
// note: this method should only be used for API response data and not user input
func ParsePython3PackageIDInsensitively(input string) (*Python3PackageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&Python3PackageId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Python3PackageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *Python3PackageId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.Python3PackageName, ok = input.Parsed["python3PackageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "python3PackageName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticPython3Packages", "python3Packages", "python3Packages"),
		resourceids.UserSpecifiedSegment("python3PackageName", "python3PackageName"),
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
