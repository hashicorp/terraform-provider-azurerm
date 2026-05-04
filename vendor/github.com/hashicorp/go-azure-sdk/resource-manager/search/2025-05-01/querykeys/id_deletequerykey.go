package querykeys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DeleteQueryKeyId{})
}

var _ resourceids.ResourceId = &DeleteQueryKeyId{}

// DeleteQueryKeyId is a struct representing the Resource ID for a Delete Query Key
type DeleteQueryKeyId struct {
	SubscriptionId     string
	ResourceGroupName  string
	SearchServiceName  string
	DeleteQueryKeyName string
}

// NewDeleteQueryKeyID returns a new DeleteQueryKeyId struct
func NewDeleteQueryKeyID(subscriptionId string, resourceGroupName string, searchServiceName string, deleteQueryKeyName string) DeleteQueryKeyId {
	return DeleteQueryKeyId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		SearchServiceName:  searchServiceName,
		DeleteQueryKeyName: deleteQueryKeyName,
	}
}

// ParseDeleteQueryKeyID parses 'input' into a DeleteQueryKeyId
func ParseDeleteQueryKeyID(input string) (*DeleteQueryKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeleteQueryKeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeleteQueryKeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDeleteQueryKeyIDInsensitively parses 'input' case-insensitively into a DeleteQueryKeyId
// note: this method should only be used for API response data and not user input
func ParseDeleteQueryKeyIDInsensitively(input string) (*DeleteQueryKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeleteQueryKeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeleteQueryKeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DeleteQueryKeyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SearchServiceName, ok = input.Parsed["searchServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "searchServiceName", input)
	}

	if id.DeleteQueryKeyName, ok = input.Parsed["deleteQueryKeyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deleteQueryKeyName", input)
	}

	return nil
}

// ValidateDeleteQueryKeyID checks that 'input' can be parsed as a Delete Query Key ID
func ValidateDeleteQueryKeyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeleteQueryKeyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Delete Query Key ID
func (id DeleteQueryKeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Search/searchServices/%s/deleteQueryKey/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SearchServiceName, id.DeleteQueryKeyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Delete Query Key ID
func (id DeleteQueryKeyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSearch", "Microsoft.Search", "Microsoft.Search"),
		resourceids.StaticSegment("staticSearchServices", "searchServices", "searchServices"),
		resourceids.UserSpecifiedSegment("searchServiceName", "searchServiceName"),
		resourceids.StaticSegment("staticDeleteQueryKey", "deleteQueryKey", "deleteQueryKey"),
		resourceids.UserSpecifiedSegment("deleteQueryKeyName", "deleteQueryKeyName"),
	}
}

// String returns a human-readable description of this Delete Query Key ID
func (id DeleteQueryKeyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Search Service Name: %q", id.SearchServiceName),
		fmt.Sprintf("Delete Query Key Name: %q", id.DeleteQueryKeyName),
	}
	return fmt.Sprintf("Delete Query Key (%s)", strings.Join(components, "\n"))
}
