package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MongodbDatabaseId{})
}

var _ resourceids.ResourceId = &MongodbDatabaseId{}

// MongodbDatabaseId is a struct representing the Resource ID for a Mongodb Database
type MongodbDatabaseId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	MongodbDatabaseName string
}

// NewMongodbDatabaseID returns a new MongodbDatabaseId struct
func NewMongodbDatabaseID(subscriptionId string, resourceGroupName string, databaseAccountName string, mongodbDatabaseName string) MongodbDatabaseId {
	return MongodbDatabaseId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		MongodbDatabaseName: mongodbDatabaseName,
	}
}

// ParseMongodbDatabaseID parses 'input' into a MongodbDatabaseId
func ParseMongodbDatabaseID(input string) (*MongodbDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MongodbDatabaseId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MongodbDatabaseId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMongodbDatabaseIDInsensitively parses 'input' case-insensitively into a MongodbDatabaseId
// note: this method should only be used for API response data and not user input
func ParseMongodbDatabaseIDInsensitively(input string) (*MongodbDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MongodbDatabaseId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MongodbDatabaseId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MongodbDatabaseId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.MongodbDatabaseName, ok = input.Parsed["mongodbDatabaseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "mongodbDatabaseName", input)
	}

	return nil
}

// ValidateMongodbDatabaseID checks that 'input' can be parsed as a Mongodb Database ID
func ValidateMongodbDatabaseID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMongodbDatabaseID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Mongodb Database ID
func (id MongodbDatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/mongodbDatabases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.MongodbDatabaseName)
}

// Segments returns a slice of Resource ID Segments which comprise this Mongodb Database ID
func (id MongodbDatabaseId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountName"),
		resourceids.StaticSegment("staticMongodbDatabases", "mongodbDatabases", "mongodbDatabases"),
		resourceids.UserSpecifiedSegment("mongodbDatabaseName", "mongodbDatabaseName"),
	}
}

// String returns a human-readable description of this Mongodb Database ID
func (id MongodbDatabaseId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Mongodb Database Name: %q", id.MongodbDatabaseName),
	}
	return fmt.Sprintf("Mongodb Database (%s)", strings.Join(components, "\n"))
}
