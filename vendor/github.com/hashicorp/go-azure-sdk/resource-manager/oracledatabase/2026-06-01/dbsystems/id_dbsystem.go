package dbsystems

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DbSystemId{})
}

var _ resourceids.ResourceId = &DbSystemId{}

// DbSystemId is a struct representing the Resource ID for a Db System
type DbSystemId struct {
	SubscriptionId    string
	ResourceGroupName string
	DbSystemName      string
}

// NewDbSystemID returns a new DbSystemId struct
func NewDbSystemID(subscriptionId string, resourceGroupName string, dbSystemName string) DbSystemId {
	return DbSystemId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DbSystemName:      dbSystemName,
	}
}

// ParseDbSystemID parses 'input' into a DbSystemId
func ParseDbSystemID(input string) (*DbSystemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DbSystemId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DbSystemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDbSystemIDInsensitively parses 'input' case-insensitively into a DbSystemId
// note: this method should only be used for API response data and not user input
func ParseDbSystemIDInsensitively(input string) (*DbSystemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DbSystemId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DbSystemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DbSystemId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DbSystemName, ok = input.Parsed["dbSystemName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dbSystemName", input)
	}

	return nil
}

// ValidateDbSystemID checks that 'input' can be parsed as a Db System ID
func ValidateDbSystemID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDbSystemID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Db System ID
func (id DbSystemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/dbSystems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DbSystemName)
}

// Segments returns a slice of Resource ID Segments which comprise this Db System ID
func (id DbSystemId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticDbSystems", "dbSystems", "dbSystems"),
		resourceids.UserSpecifiedSegment("dbSystemName", "dbSystemName"),
	}
}

// String returns a human-readable description of this Db System ID
func (id DbSystemId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Db System Name: %q", id.DbSystemName),
	}
	return fmt.Sprintf("Db System (%s)", strings.Join(components, "\n"))
}
