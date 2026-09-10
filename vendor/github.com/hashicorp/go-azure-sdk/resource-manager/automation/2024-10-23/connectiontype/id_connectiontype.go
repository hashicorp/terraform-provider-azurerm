package connectiontype

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConnectionTypeId{})
}

var _ resourceids.ResourceId = &ConnectionTypeId{}

// ConnectionTypeId is a struct representing the Resource ID for a Connection Type
type ConnectionTypeId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	ConnectionTypeName    string
}

// NewConnectionTypeID returns a new ConnectionTypeId struct
func NewConnectionTypeID(subscriptionId string, resourceGroupName string, automationAccountName string, connectionTypeName string) ConnectionTypeId {
	return ConnectionTypeId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		ConnectionTypeName:    connectionTypeName,
	}
}

// ParseConnectionTypeID parses 'input' into a ConnectionTypeId
func ParseConnectionTypeID(input string) (*ConnectionTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectionTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectionTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConnectionTypeIDInsensitively parses 'input' case-insensitively into a ConnectionTypeId
// note: this method should only be used for API response data and not user input
func ParseConnectionTypeIDInsensitively(input string) (*ConnectionTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectionTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectionTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConnectionTypeId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ConnectionTypeName, ok = input.Parsed["connectionTypeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "connectionTypeName", input)
	}

	return nil
}

// ValidateConnectionTypeID checks that 'input' can be parsed as a Connection Type ID
func ValidateConnectionTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConnectionTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Connection Type ID
func (id ConnectionTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/connectionTypes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.ConnectionTypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Connection Type ID
func (id ConnectionTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticConnectionTypes", "connectionTypes", "connectionTypes"),
		resourceids.UserSpecifiedSegment("connectionTypeName", "connectionTypeName"),
	}
}

// String returns a human-readable description of this Connection Type ID
func (id ConnectionTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Connection Type Name: %q", id.ConnectionTypeName),
	}
	return fmt.Sprintf("Connection Type (%s)", strings.Join(components, "\n"))
}
