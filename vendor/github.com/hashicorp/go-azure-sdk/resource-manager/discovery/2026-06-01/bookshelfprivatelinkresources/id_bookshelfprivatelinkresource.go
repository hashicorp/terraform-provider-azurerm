package bookshelfprivatelinkresources

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BookshelfPrivateLinkResourceId{})
}

var _ resourceids.ResourceId = &BookshelfPrivateLinkResourceId{}

// BookshelfPrivateLinkResourceId is a struct representing the Resource ID for a Bookshelf Private Link Resource
type BookshelfPrivateLinkResourceId struct {
	SubscriptionId          string
	ResourceGroupName       string
	BookshelfName           string
	PrivateLinkResourceName string
}

// NewBookshelfPrivateLinkResourceID returns a new BookshelfPrivateLinkResourceId struct
func NewBookshelfPrivateLinkResourceID(subscriptionId string, resourceGroupName string, bookshelfName string, privateLinkResourceName string) BookshelfPrivateLinkResourceId {
	return BookshelfPrivateLinkResourceId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		BookshelfName:           bookshelfName,
		PrivateLinkResourceName: privateLinkResourceName,
	}
}

// ParseBookshelfPrivateLinkResourceID parses 'input' into a BookshelfPrivateLinkResourceId
func ParseBookshelfPrivateLinkResourceID(input string) (*BookshelfPrivateLinkResourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BookshelfPrivateLinkResourceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BookshelfPrivateLinkResourceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBookshelfPrivateLinkResourceIDInsensitively parses 'input' case-insensitively into a BookshelfPrivateLinkResourceId
// note: this method should only be used for API response data and not user input
func ParseBookshelfPrivateLinkResourceIDInsensitively(input string) (*BookshelfPrivateLinkResourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BookshelfPrivateLinkResourceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BookshelfPrivateLinkResourceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BookshelfPrivateLinkResourceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.BookshelfName, ok = input.Parsed["bookshelfName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "bookshelfName", input)
	}

	if id.PrivateLinkResourceName, ok = input.Parsed["privateLinkResourceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateLinkResourceName", input)
	}

	return nil
}

// ValidateBookshelfPrivateLinkResourceID checks that 'input' can be parsed as a Bookshelf Private Link Resource ID
func ValidateBookshelfPrivateLinkResourceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBookshelfPrivateLinkResourceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Bookshelf Private Link Resource ID
func (id BookshelfPrivateLinkResourceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Discovery/bookshelves/%s/privateLinkResources/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.BookshelfName, id.PrivateLinkResourceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Bookshelf Private Link Resource ID
func (id BookshelfPrivateLinkResourceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDiscovery", "Microsoft.Discovery", "Microsoft.Discovery"),
		resourceids.StaticSegment("staticBookshelves", "bookshelves", "bookshelves"),
		resourceids.UserSpecifiedSegment("bookshelfName", "bookshelfName"),
		resourceids.StaticSegment("staticPrivateLinkResources", "privateLinkResources", "privateLinkResources"),
		resourceids.UserSpecifiedSegment("privateLinkResourceName", "privateLinkResourceName"),
	}
}

// String returns a human-readable description of this Bookshelf Private Link Resource ID
func (id BookshelfPrivateLinkResourceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Bookshelf Name: %q", id.BookshelfName),
		fmt.Sprintf("Private Link Resource Name: %q", id.PrivateLinkResourceName),
	}
	return fmt.Sprintf("Bookshelf Private Link Resource (%s)", strings.Join(components, "\n"))
}
