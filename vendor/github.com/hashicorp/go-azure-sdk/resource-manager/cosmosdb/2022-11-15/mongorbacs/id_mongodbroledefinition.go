package mongorbacs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = MongodbRoleDefinitionId{}

// MongodbRoleDefinitionId is a struct representing the Resource ID for a Mongodb Role Definition
type MongodbRoleDefinitionId struct {
	SubscriptionId        string
	ResourceGroupName     string
	DatabaseAccountName   string
	MongoRoleDefinitionId string
}

// NewMongodbRoleDefinitionID returns a new MongodbRoleDefinitionId struct
func NewMongodbRoleDefinitionID(subscriptionId string, resourceGroupName string, databaseAccountName string, mongoRoleDefinitionId string) MongodbRoleDefinitionId {
	return MongodbRoleDefinitionId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		DatabaseAccountName:   databaseAccountName,
		MongoRoleDefinitionId: mongoRoleDefinitionId,
	}
}

// ParseMongodbRoleDefinitionID parses 'input' into a MongodbRoleDefinitionId
func ParseMongodbRoleDefinitionID(input string) (*MongodbRoleDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(MongodbRoleDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MongodbRoleDefinitionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.MongoRoleDefinitionId, ok = parsed.Parsed["mongoRoleDefinitionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mongoRoleDefinitionId", *parsed)
	}

	return &id, nil
}

// ParseMongodbRoleDefinitionIDInsensitively parses 'input' case-insensitively into a MongodbRoleDefinitionId
// note: this method should only be used for API response data and not user input
func ParseMongodbRoleDefinitionIDInsensitively(input string) (*MongodbRoleDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(MongodbRoleDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MongodbRoleDefinitionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.MongoRoleDefinitionId, ok = parsed.Parsed["mongoRoleDefinitionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mongoRoleDefinitionId", *parsed)
	}

	return &id, nil
}

// ValidateMongodbRoleDefinitionID checks that 'input' can be parsed as a Mongodb Role Definition ID
func ValidateMongodbRoleDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMongodbRoleDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Mongodb Role Definition ID
func (id MongodbRoleDefinitionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/mongodbRoleDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.MongoRoleDefinitionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Mongodb Role Definition ID
func (id MongodbRoleDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountValue"),
		resourceids.StaticSegment("staticMongodbRoleDefinitions", "mongodbRoleDefinitions", "mongodbRoleDefinitions"),
		resourceids.UserSpecifiedSegment("mongoRoleDefinitionId", "mongoRoleDefinitionIdValue"),
	}
}

// String returns a human-readable description of this Mongodb Role Definition ID
func (id MongodbRoleDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Mongo Role Definition: %q", id.MongoRoleDefinitionId),
	}
	return fmt.Sprintf("Mongodb Role Definition (%s)", strings.Join(components, "\n"))
}
