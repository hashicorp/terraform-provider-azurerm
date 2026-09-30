package bookshelfprivateendpointconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BookshelfPrivateEndpointConnectionId{})
}

var _ resourceids.ResourceId = &BookshelfPrivateEndpointConnectionId{}

// BookshelfPrivateEndpointConnectionId is a struct representing the Resource ID for a Bookshelf Private Endpoint Connection
type BookshelfPrivateEndpointConnectionId struct {
	SubscriptionId                string
	ResourceGroupName             string
	BookshelfName                 string
	PrivateEndpointConnectionName string
}

// NewBookshelfPrivateEndpointConnectionID returns a new BookshelfPrivateEndpointConnectionId struct
func NewBookshelfPrivateEndpointConnectionID(subscriptionId string, resourceGroupName string, bookshelfName string, privateEndpointConnectionName string) BookshelfPrivateEndpointConnectionId {
	return BookshelfPrivateEndpointConnectionId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		BookshelfName:                 bookshelfName,
		PrivateEndpointConnectionName: privateEndpointConnectionName,
	}
}

// ParseBookshelfPrivateEndpointConnectionID parses 'input' into a BookshelfPrivateEndpointConnectionId
func ParseBookshelfPrivateEndpointConnectionID(input string) (*BookshelfPrivateEndpointConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BookshelfPrivateEndpointConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BookshelfPrivateEndpointConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBookshelfPrivateEndpointConnectionIDInsensitively parses 'input' case-insensitively into a BookshelfPrivateEndpointConnectionId
// note: this method should only be used for API response data and not user input
func ParseBookshelfPrivateEndpointConnectionIDInsensitively(input string) (*BookshelfPrivateEndpointConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BookshelfPrivateEndpointConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BookshelfPrivateEndpointConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BookshelfPrivateEndpointConnectionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.PrivateEndpointConnectionName, ok = input.Parsed["privateEndpointConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointConnectionName", input)
	}

	return nil
}

// ValidateBookshelfPrivateEndpointConnectionID checks that 'input' can be parsed as a Bookshelf Private Endpoint Connection ID
func ValidateBookshelfPrivateEndpointConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBookshelfPrivateEndpointConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Bookshelf Private Endpoint Connection ID
func (id BookshelfPrivateEndpointConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Discovery/bookshelves/%s/privateEndpointConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.BookshelfName, id.PrivateEndpointConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Bookshelf Private Endpoint Connection ID
func (id BookshelfPrivateEndpointConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDiscovery", "Microsoft.Discovery", "Microsoft.Discovery"),
		resourceids.StaticSegment("staticBookshelves", "bookshelves", "bookshelves"),
		resourceids.UserSpecifiedSegment("bookshelfName", "bookshelfName"),
		resourceids.StaticSegment("staticPrivateEndpointConnections", "privateEndpointConnections", "privateEndpointConnections"),
		resourceids.UserSpecifiedSegment("privateEndpointConnectionName", "privateEndpointConnectionName"),
	}
}

// String returns a human-readable description of this Bookshelf Private Endpoint Connection ID
func (id BookshelfPrivateEndpointConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Bookshelf Name: %q", id.BookshelfName),
		fmt.Sprintf("Private Endpoint Connection Name: %q", id.PrivateEndpointConnectionName),
	}
	return fmt.Sprintf("Bookshelf Private Endpoint Connection (%s)", strings.Join(components, "\n"))
}
