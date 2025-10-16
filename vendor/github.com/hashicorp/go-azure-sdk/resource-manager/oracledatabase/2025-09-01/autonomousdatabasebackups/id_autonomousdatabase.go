package autonomousdatabasebackups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AutonomousDatabaseId{})
}

var _ resourceids.ResourceId = &AutonomousDatabaseId{}

// AutonomousDatabaseId is a struct representing the Resource ID for a Autonomous Database
type AutonomousDatabaseId struct {
	SubscriptionId         string
	ResourceGroupName      string
	AutonomousDatabaseName string
}

// NewAutonomousDatabaseID returns a new AutonomousDatabaseId struct
func NewAutonomousDatabaseID(subscriptionId string, resourceGroupName string, autonomousDatabaseName string) AutonomousDatabaseId {
	return AutonomousDatabaseId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		AutonomousDatabaseName: autonomousDatabaseName,
	}
}

// ParseAutonomousDatabaseID parses 'input' into a AutonomousDatabaseId
func ParseAutonomousDatabaseID(input string) (*AutonomousDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutonomousDatabaseId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutonomousDatabaseId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAutonomousDatabaseIDInsensitively parses 'input' case-insensitively into a AutonomousDatabaseId
// note: this method should only be used for API response data and not user input
func ParseAutonomousDatabaseIDInsensitively(input string) (*AutonomousDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutonomousDatabaseId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutonomousDatabaseId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AutonomousDatabaseId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AutonomousDatabaseName, ok = input.Parsed["autonomousDatabaseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "autonomousDatabaseName", input)
	}

	return nil
}

// ValidateAutonomousDatabaseID checks that 'input' can be parsed as a Autonomous Database ID
func ValidateAutonomousDatabaseID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAutonomousDatabaseID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Autonomous Database ID
func (id AutonomousDatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/autonomousDatabases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutonomousDatabaseName)
}

// Segments returns a slice of Resource ID Segments which comprise this Autonomous Database ID
func (id AutonomousDatabaseId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticAutonomousDatabases", "autonomousDatabases", "autonomousDatabases"),
		resourceids.UserSpecifiedSegment("autonomousDatabaseName", "autonomousDatabaseName"),
	}
}

// String returns a human-readable description of this Autonomous Database ID
func (id AutonomousDatabaseId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Autonomous Database Name: %q", id.AutonomousDatabaseName),
	}
	return fmt.Sprintf("Autonomous Database (%s)", strings.Join(components, "\n"))
}
