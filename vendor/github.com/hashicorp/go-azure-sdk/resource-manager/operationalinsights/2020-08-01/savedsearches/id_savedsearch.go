package savedsearches

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SavedSearchId{})
}

var _ resourceids.ResourceId = &SavedSearchId{}

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
	parser := resourceids.NewParserFromResourceIdType(&SavedSearchId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SavedSearchId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSavedSearchIDInsensitively parses 'input' case-insensitively into a SavedSearchId
// note: this method should only be used for API response data and not user input
func ParseSavedSearchIDInsensitively(input string) (*SavedSearchId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SavedSearchId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SavedSearchId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SavedSearchId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.SavedSearchId, ok = input.Parsed["savedSearchId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "savedSearchId", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceName"),
		resourceids.StaticSegment("staticSavedSearches", "savedSearches", "savedSearches"),
		resourceids.UserSpecifiedSegment("savedSearchId", "savedSearchId"),
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
