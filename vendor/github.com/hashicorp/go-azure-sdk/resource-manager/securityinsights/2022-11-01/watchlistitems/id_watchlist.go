package watchlistitems

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = WatchlistId{}

// WatchlistId is a struct representing the Resource ID for a Watchlist
type WatchlistId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	WatchlistAlias    string
}

// NewWatchlistID returns a new WatchlistId struct
func NewWatchlistID(subscriptionId string, resourceGroupName string, workspaceName string, watchlistAlias string) WatchlistId {
	return WatchlistId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		WatchlistAlias:    watchlistAlias,
	}
}

// ParseWatchlistID parses 'input' into a WatchlistId
func ParseWatchlistID(input string) (*WatchlistId, error) {
	parser := resourceids.NewParserFromResourceIdType(WatchlistId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WatchlistId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.WatchlistAlias, ok = parsed.Parsed["watchlistAlias"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "watchlistAlias", *parsed)
	}

	return &id, nil
}

// ParseWatchlistIDInsensitively parses 'input' case-insensitively into a WatchlistId
// note: this method should only be used for API response data and not user input
func ParseWatchlistIDInsensitively(input string) (*WatchlistId, error) {
	parser := resourceids.NewParserFromResourceIdType(WatchlistId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WatchlistId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.WatchlistAlias, ok = parsed.Parsed["watchlistAlias"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "watchlistAlias", *parsed)
	}

	return &id, nil
}

// ValidateWatchlistID checks that 'input' can be parsed as a Watchlist ID
func ValidateWatchlistID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWatchlistID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Watchlist ID
func (id WatchlistId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/watchlists/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.WatchlistAlias)
}

// Segments returns a slice of Resource ID Segments which comprise this Watchlist ID
func (id WatchlistId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSecurityInsights", "Microsoft.SecurityInsights", "Microsoft.SecurityInsights"),
		resourceids.StaticSegment("staticWatchlists", "watchlists", "watchlists"),
		resourceids.UserSpecifiedSegment("watchlistAlias", "watchlistAliasValue"),
	}
}

// String returns a human-readable description of this Watchlist ID
func (id WatchlistId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Watchlist Alias: %q", id.WatchlistAlias),
	}
	return fmt.Sprintf("Watchlist (%s)", strings.Join(components, "\n"))
}
