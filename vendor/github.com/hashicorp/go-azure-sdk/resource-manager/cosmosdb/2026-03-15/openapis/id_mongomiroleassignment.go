package openapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MongoMIRoleAssignmentId{})
}

var _ resourceids.ResourceId = &MongoMIRoleAssignmentId{}

// MongoMIRoleAssignmentId is a struct representing the Resource ID for a Mongo MI Role Assignment
type MongoMIRoleAssignmentId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	RoleAssignmentId    string
}

// NewMongoMIRoleAssignmentID returns a new MongoMIRoleAssignmentId struct
func NewMongoMIRoleAssignmentID(subscriptionId string, resourceGroupName string, databaseAccountName string, roleAssignmentId string) MongoMIRoleAssignmentId {
	return MongoMIRoleAssignmentId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		RoleAssignmentId:    roleAssignmentId,
	}
}

// ParseMongoMIRoleAssignmentID parses 'input' into a MongoMIRoleAssignmentId
func ParseMongoMIRoleAssignmentID(input string) (*MongoMIRoleAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MongoMIRoleAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MongoMIRoleAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMongoMIRoleAssignmentIDInsensitively parses 'input' case-insensitively into a MongoMIRoleAssignmentId
// note: this method should only be used for API response data and not user input
func ParseMongoMIRoleAssignmentIDInsensitively(input string) (*MongoMIRoleAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MongoMIRoleAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MongoMIRoleAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MongoMIRoleAssignmentId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.RoleAssignmentId, ok = input.Parsed["roleAssignmentId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleAssignmentId", input)
	}

	return nil
}

// ValidateMongoMIRoleAssignmentID checks that 'input' can be parsed as a Mongo MI Role Assignment ID
func ValidateMongoMIRoleAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMongoMIRoleAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Mongo MI Role Assignment ID
func (id MongoMIRoleAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/mongoMIRoleAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.RoleAssignmentId)
}

// Segments returns a slice of Resource ID Segments which comprise this Mongo MI Role Assignment ID
func (id MongoMIRoleAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountName"),
		resourceids.StaticSegment("staticMongoMIRoleAssignments", "mongoMIRoleAssignments", "mongoMIRoleAssignments"),
		resourceids.UserSpecifiedSegment("roleAssignmentId", "roleAssignmentId"),
	}
}

// String returns a human-readable description of this Mongo MI Role Assignment ID
func (id MongoMIRoleAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Role Assignment: %q", id.RoleAssignmentId),
	}
	return fmt.Sprintf("Mongo MI Role Assignment (%s)", strings.Join(components, "\n"))
}
