package querykeys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CreateQueryKeyId{}

// CreateQueryKeyId is a struct representing the Resource ID for a Create Query Key
type CreateQueryKeyId struct {
	SubscriptionId     string
	ResourceGroupName  string
	SearchServiceName  string
	CreateQueryKeyName string
}

// NewCreateQueryKeyID returns a new CreateQueryKeyId struct
func NewCreateQueryKeyID(subscriptionId string, resourceGroupName string, searchServiceName string, createQueryKeyName string) CreateQueryKeyId {
	return CreateQueryKeyId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		SearchServiceName:  searchServiceName,
		CreateQueryKeyName: createQueryKeyName,
	}
}

// ParseCreateQueryKeyID parses 'input' into a CreateQueryKeyId
func ParseCreateQueryKeyID(input string) (*CreateQueryKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(CreateQueryKeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CreateQueryKeyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SearchServiceName, ok = parsed.Parsed["searchServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "searchServiceName", *parsed)
	}

	if id.CreateQueryKeyName, ok = parsed.Parsed["createQueryKeyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "createQueryKeyName", *parsed)
	}

	return &id, nil
}

// ParseCreateQueryKeyIDInsensitively parses 'input' case-insensitively into a CreateQueryKeyId
// note: this method should only be used for API response data and not user input
func ParseCreateQueryKeyIDInsensitively(input string) (*CreateQueryKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(CreateQueryKeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CreateQueryKeyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SearchServiceName, ok = parsed.Parsed["searchServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "searchServiceName", *parsed)
	}

	if id.CreateQueryKeyName, ok = parsed.Parsed["createQueryKeyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "createQueryKeyName", *parsed)
	}

	return &id, nil
}

// ValidateCreateQueryKeyID checks that 'input' can be parsed as a Create Query Key ID
func ValidateCreateQueryKeyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCreateQueryKeyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Create Query Key ID
func (id CreateQueryKeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Search/searchServices/%s/createQueryKey/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SearchServiceName, id.CreateQueryKeyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Create Query Key ID
func (id CreateQueryKeyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSearch", "Microsoft.Search", "Microsoft.Search"),
		resourceids.StaticSegment("staticSearchServices", "searchServices", "searchServices"),
		resourceids.UserSpecifiedSegment("searchServiceName", "searchServiceValue"),
		resourceids.StaticSegment("staticCreateQueryKey", "createQueryKey", "createQueryKey"),
		resourceids.UserSpecifiedSegment("createQueryKeyName", "createQueryKeyValue"),
	}
}

// String returns a human-readable description of this Create Query Key ID
func (id CreateQueryKeyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Search Service Name: %q", id.SearchServiceName),
		fmt.Sprintf("Create Query Key Name: %q", id.CreateQueryKeyName),
	}
	return fmt.Sprintf("Create Query Key (%s)", strings.Join(components, "\n"))
}
