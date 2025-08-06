package python2package

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&Python2PackageId{})
}

var _ resourceids.ResourceId = &Python2PackageId{}

// Python2PackageId is a struct representing the Resource ID for a Python 2 Package
type Python2PackageId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	Python2PackageName    string
}

// NewPython2PackageID returns a new Python2PackageId struct
func NewPython2PackageID(subscriptionId string, resourceGroupName string, automationAccountName string, python2PackageName string) Python2PackageId {
	return Python2PackageId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		Python2PackageName:    python2PackageName,
	}
}

// ParsePython2PackageID parses 'input' into a Python2PackageId
func ParsePython2PackageID(input string) (*Python2PackageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&Python2PackageId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Python2PackageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePython2PackageIDInsensitively parses 'input' case-insensitively into a Python2PackageId
// note: this method should only be used for API response data and not user input
func ParsePython2PackageIDInsensitively(input string) (*Python2PackageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&Python2PackageId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Python2PackageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *Python2PackageId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.Python2PackageName, ok = input.Parsed["python2PackageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "python2PackageName", input)
	}

	return nil
}

// ValidatePython2PackageID checks that 'input' can be parsed as a Python 2 Package ID
func ValidatePython2PackageID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePython2PackageID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Python 2 Package ID
func (id Python2PackageId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/python2Packages/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.Python2PackageName)
}

// Segments returns a slice of Resource ID Segments which comprise this Python 2 Package ID
func (id Python2PackageId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticPython2Packages", "python2Packages", "python2Packages"),
		resourceids.UserSpecifiedSegment("python2PackageName", "python2PackageName"),
	}
}

// String returns a human-readable description of this Python 2 Package ID
func (id Python2PackageId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Python 2 Package Name: %q", id.Python2PackageName),
	}
	return fmt.Sprintf("Python 2 Package (%s)", strings.Join(components, "\n"))
}
