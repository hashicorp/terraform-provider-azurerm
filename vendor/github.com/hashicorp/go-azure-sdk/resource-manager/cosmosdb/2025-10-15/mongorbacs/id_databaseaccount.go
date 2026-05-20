package mongorbacs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DatabaseAccountId{})
}

var _ resourceids.ResourceId = &DatabaseAccountId{}

// DatabaseAccountId is a struct representing the Resource ID for a Database Account
type DatabaseAccountId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
}

// NewDatabaseAccountID returns a new DatabaseAccountId struct
func NewDatabaseAccountID(subscriptionId string, resourceGroupName string, databaseAccountName string) DatabaseAccountId {
	return DatabaseAccountId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
	}
}

// ParseDatabaseAccountID parses 'input' into a DatabaseAccountId
func ParseDatabaseAccountID(input string) (*DatabaseAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabaseAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabaseAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDatabaseAccountIDInsensitively parses 'input' case-insensitively into a DatabaseAccountId
// note: this method should only be used for API response data and not user input
func ParseDatabaseAccountIDInsensitively(input string) (*DatabaseAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabaseAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabaseAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DatabaseAccountId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DatabaseAccountName, ok = input.Parsed["databaseAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", input)
	}

	return nil
}

// ValidateDatabaseAccountID checks that 'input' can be parsed as a Database Account ID
func ValidateDatabaseAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDatabaseAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Database Account ID
func (id DatabaseAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Database Account ID
func (id DatabaseAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountName"),
	}
}

// String returns a human-readable description of this Database Account ID
func (id DatabaseAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
	}
	return fmt.Sprintf("Database Account (%s)", strings.Join(components, "\n"))
}
