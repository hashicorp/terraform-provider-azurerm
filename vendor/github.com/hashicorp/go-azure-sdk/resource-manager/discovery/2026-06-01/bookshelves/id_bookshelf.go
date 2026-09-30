package bookshelves

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BookshelfId{})
}

var _ resourceids.ResourceId = &BookshelfId{}

// BookshelfId is a struct representing the Resource ID for a Bookshelf
type BookshelfId struct {
	SubscriptionId    string
	ResourceGroupName string
	BookshelfName     string
}

// NewBookshelfID returns a new BookshelfId struct
func NewBookshelfID(subscriptionId string, resourceGroupName string, bookshelfName string) BookshelfId {
	return BookshelfId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		BookshelfName:     bookshelfName,
	}
}

// ParseBookshelfID parses 'input' into a BookshelfId
func ParseBookshelfID(input string) (*BookshelfId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BookshelfId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BookshelfId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBookshelfIDInsensitively parses 'input' case-insensitively into a BookshelfId
// note: this method should only be used for API response data and not user input
func ParseBookshelfIDInsensitively(input string) (*BookshelfId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BookshelfId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BookshelfId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BookshelfId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateBookshelfID checks that 'input' can be parsed as a Bookshelf ID
func ValidateBookshelfID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBookshelfID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Bookshelf ID
func (id BookshelfId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Discovery/bookshelves/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.BookshelfName)
}

// Segments returns a slice of Resource ID Segments which comprise this Bookshelf ID
func (id BookshelfId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDiscovery", "Microsoft.Discovery", "Microsoft.Discovery"),
		resourceids.StaticSegment("staticBookshelves", "bookshelves", "bookshelves"),
		resourceids.UserSpecifiedSegment("bookshelfName", "bookshelfName"),
	}
}

// String returns a human-readable description of this Bookshelf ID
func (id BookshelfId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Bookshelf Name: %q", id.BookshelfName),
	}
	return fmt.Sprintf("Bookshelf (%s)", strings.Join(components, "\n"))
}
