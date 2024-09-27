package watchlistitems

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WatchlistItemId{})
}

var _ resourceids.ResourceId = &WatchlistItemId{}

// WatchlistItemId is a struct representing the Resource ID for a Watchlist Item
type WatchlistItemId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	WatchlistAlias    string
	WatchlistItemId   string
}

// NewWatchlistItemID returns a new WatchlistItemId struct
func NewWatchlistItemID(subscriptionId string, resourceGroupName string, workspaceName string, watchlistAlias string, watchlistItemId string) WatchlistItemId {
	return WatchlistItemId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		WatchlistAlias:    watchlistAlias,
		WatchlistItemId:   watchlistItemId,
	}
}

// ParseWatchlistItemID parses 'input' into a WatchlistItemId
func ParseWatchlistItemID(input string) (*WatchlistItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WatchlistItemId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WatchlistItemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWatchlistItemIDInsensitively parses 'input' case-insensitively into a WatchlistItemId
// note: this method should only be used for API response data and not user input
func ParseWatchlistItemIDInsensitively(input string) (*WatchlistItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WatchlistItemId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WatchlistItemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WatchlistItemId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.WorkspaceName, ok = input.Parsed["workspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", input)
	}

	if id.WatchlistAlias, ok = input.Parsed["watchlistAlias"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "watchlistAlias", input)
	}

	if id.WatchlistItemId, ok = input.Parsed["watchlistItemId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "watchlistItemId", input)
	}

	return nil
}

// ValidateWatchlistItemID checks that 'input' can be parsed as a Watchlist Item ID
func ValidateWatchlistItemID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWatchlistItemID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Watchlist Item ID
func (id WatchlistItemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/watchlists/%s/watchlistItems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.WatchlistAlias, id.WatchlistItemId)
}

// Segments returns a slice of Resource ID Segments which comprise this Watchlist Item ID
func (id WatchlistItemId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceName"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSecurityInsights", "Microsoft.SecurityInsights", "Microsoft.SecurityInsights"),
		resourceids.StaticSegment("staticWatchlists", "watchlists", "watchlists"),
		resourceids.UserSpecifiedSegment("watchlistAlias", "watchlistAlias"),
		resourceids.StaticSegment("staticWatchlistItems", "watchlistItems", "watchlistItems"),
		resourceids.UserSpecifiedSegment("watchlistItemId", "watchlistItemId"),
	}
}

// String returns a human-readable description of this Watchlist Item ID
func (id WatchlistItemId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Watchlist Alias: %q", id.WatchlistAlias),
		fmt.Sprintf("Watchlist Item: %q", id.WatchlistItemId),
	}
	return fmt.Sprintf("Watchlist Item (%s)", strings.Join(components, "\n"))
}
