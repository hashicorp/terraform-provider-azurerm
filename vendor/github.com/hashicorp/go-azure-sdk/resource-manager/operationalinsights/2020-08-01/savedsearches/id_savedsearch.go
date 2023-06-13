package savedsearches

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SavedSearchId{}

// SavedSearchId is a struct representing the Resource ID for a Saved Search
type SavedSearchId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	SavedSearchId     string
}

// NewSavedSearchID returns a new SavedSearchId struct
func NewSavedSearchID(subscriptionId string, resourceGroupName string, workspaceName string, savedSearchId string) SavedSearchId {
	return SavedSearchId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		SavedSearchId:     savedSearchId,
	}
}

// ParseSavedSearchID parses 'input' into a SavedSearchId
func ParseSavedSearchID(input string) (*SavedSearchId, error) {
	parser := resourceids.NewParserFromResourceIdType(SavedSearchId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SavedSearchId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.SavedSearchId, ok = parsed.Parsed["savedSearchId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "savedSearchId", *parsed)
	}

	return &id, nil
}

// ParseSavedSearchIDInsensitively parses 'input' case-insensitively into a SavedSearchId
// note: this method should only be used for API response data and not user input
func ParseSavedSearchIDInsensitively(input string) (*SavedSearchId, error) {
	parser := resourceids.NewParserFromResourceIdType(SavedSearchId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SavedSearchId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.SavedSearchId, ok = parsed.Parsed["savedSearchId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "savedSearchId", *parsed)
	}

	return &id, nil
}

// ValidateSavedSearchID checks that 'input' can be parsed as a Saved Search ID
func ValidateSavedSearchID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSavedSearchID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Saved Search ID
func (id SavedSearchId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/savedSearches/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.SavedSearchId)
}

// Segments returns a slice of Resource ID Segments which comprise this Saved Search ID
func (id SavedSearchId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticSavedSearches", "savedSearches", "savedSearches"),
		resourceids.UserSpecifiedSegment("savedSearchId", "savedSearchIdValue"),
	}
}

// String returns a human-readable description of this Saved Search ID
func (id SavedSearchId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Saved Search: %q", id.SavedSearchId),
	}
	return fmt.Sprintf("Saved Search (%s)", strings.Join(components, "\n"))
}
