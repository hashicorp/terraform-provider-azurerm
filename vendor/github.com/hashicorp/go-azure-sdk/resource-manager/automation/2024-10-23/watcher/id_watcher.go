package watcher

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WatcherId{})
}

var _ resourceids.ResourceId = &WatcherId{}

// WatcherId is a struct representing the Resource ID for a Watcher
type WatcherId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	WatcherName           string
}

// NewWatcherID returns a new WatcherId struct
func NewWatcherID(subscriptionId string, resourceGroupName string, automationAccountName string, watcherName string) WatcherId {
	return WatcherId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		WatcherName:           watcherName,
	}
}

// ParseWatcherID parses 'input' into a WatcherId
func ParseWatcherID(input string) (*WatcherId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WatcherId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WatcherId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWatcherIDInsensitively parses 'input' case-insensitively into a WatcherId
// note: this method should only be used for API response data and not user input
func ParseWatcherIDInsensitively(input string) (*WatcherId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WatcherId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WatcherId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WatcherId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.WatcherName, ok = input.Parsed["watcherName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "watcherName", input)
	}

	return nil
}

// ValidateWatcherID checks that 'input' can be parsed as a Watcher ID
func ValidateWatcherID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWatcherID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Watcher ID
func (id WatcherId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/watchers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.WatcherName)
}

// Segments returns a slice of Resource ID Segments which comprise this Watcher ID
func (id WatcherId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticWatchers", "watchers", "watchers"),
		resourceids.UserSpecifiedSegment("watcherName", "watcherName"),
	}
}

// String returns a human-readable description of this Watcher ID
func (id WatcherId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Watcher Name: %q", id.WatcherName),
	}
	return fmt.Sprintf("Watcher (%s)", strings.Join(components, "\n"))
}
