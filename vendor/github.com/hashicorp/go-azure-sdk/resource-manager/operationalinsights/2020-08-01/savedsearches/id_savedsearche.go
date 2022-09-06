package savedsearches

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SavedSearcheId{}

// SavedSearcheId is a struct representing the Resource ID for a Saved Searche
type SavedSearcheId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	SavedSearchId     string
}

// NewSavedSearcheID returns a new SavedSearcheId struct
func NewSavedSearcheID(subscriptionId string, resourceGroupName string, workspaceName string, savedSearchId string) SavedSearcheId {
	return SavedSearcheId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		SavedSearchId:     savedSearchId,
	}
}

// ParseSavedSearcheID parses 'input' into a SavedSearcheId
func ParseSavedSearcheID(input string) (*SavedSearcheId, error) {
	parser := resourceids.NewParserFromResourceIdType(SavedSearcheId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SavedSearcheId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'workspaceName' was not found in the resource id %q", input)
	}

	if id.SavedSearchId, ok = parsed.Parsed["savedSearchId"]; !ok {
		return nil, fmt.Errorf("the segment 'savedSearchId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseSavedSearcheIDInsensitively parses 'input' case-insensitively into a SavedSearcheId
// note: this method should only be used for API response data and not user input
func ParseSavedSearcheIDInsensitively(input string) (*SavedSearcheId, error) {
	parser := resourceids.NewParserFromResourceIdType(SavedSearcheId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SavedSearcheId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'workspaceName' was not found in the resource id %q", input)
	}

	if id.SavedSearchId, ok = parsed.Parsed["savedSearchId"]; !ok {
		return nil, fmt.Errorf("the segment 'savedSearchId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateSavedSearcheID checks that 'input' can be parsed as a Saved Searche ID
func ValidateSavedSearcheID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSavedSearcheID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Saved Searche ID
func (id SavedSearcheId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/savedSearches/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.SavedSearchId)
}

// Segments returns a slice of Resource ID Segments which comprise this Saved Searche ID
func (id SavedSearcheId) Segments() []resourceids.Segment {
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

// String returns a human-readable description of this Saved Searche ID
func (id SavedSearcheId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Saved Search: %q", id.SavedSearchId),
	}
	return fmt.Sprintf("Saved Searche (%s)", strings.Join(components, "\n"))
}
