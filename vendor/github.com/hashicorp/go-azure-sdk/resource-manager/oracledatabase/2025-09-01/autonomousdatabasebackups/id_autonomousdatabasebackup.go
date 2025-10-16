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
	recaser.RegisterResourceId(&AutonomousDatabaseBackupId{})
}

var _ resourceids.ResourceId = &AutonomousDatabaseBackupId{}

// AutonomousDatabaseBackupId is a struct representing the Resource ID for a Autonomous Database Backup
type AutonomousDatabaseBackupId struct {
	SubscriptionId               string
	ResourceGroupName            string
	AutonomousDatabaseName       string
	AutonomousDatabaseBackupName string
}

// NewAutonomousDatabaseBackupID returns a new AutonomousDatabaseBackupId struct
func NewAutonomousDatabaseBackupID(subscriptionId string, resourceGroupName string, autonomousDatabaseName string, autonomousDatabaseBackupName string) AutonomousDatabaseBackupId {
	return AutonomousDatabaseBackupId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		AutonomousDatabaseName:       autonomousDatabaseName,
		AutonomousDatabaseBackupName: autonomousDatabaseBackupName,
	}
}

// ParseAutonomousDatabaseBackupID parses 'input' into a AutonomousDatabaseBackupId
func ParseAutonomousDatabaseBackupID(input string) (*AutonomousDatabaseBackupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutonomousDatabaseBackupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutonomousDatabaseBackupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAutonomousDatabaseBackupIDInsensitively parses 'input' case-insensitively into a AutonomousDatabaseBackupId
// note: this method should only be used for API response data and not user input
func ParseAutonomousDatabaseBackupIDInsensitively(input string) (*AutonomousDatabaseBackupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutonomousDatabaseBackupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutonomousDatabaseBackupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AutonomousDatabaseBackupId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AutonomousDatabaseBackupName, ok = input.Parsed["autonomousDatabaseBackupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "autonomousDatabaseBackupName", input)
	}

	return nil
}

// ValidateAutonomousDatabaseBackupID checks that 'input' can be parsed as a Autonomous Database Backup ID
func ValidateAutonomousDatabaseBackupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAutonomousDatabaseBackupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Autonomous Database Backup ID
func (id AutonomousDatabaseBackupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/autonomousDatabases/%s/autonomousDatabaseBackups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutonomousDatabaseName, id.AutonomousDatabaseBackupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Autonomous Database Backup ID
func (id AutonomousDatabaseBackupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticAutonomousDatabases", "autonomousDatabases", "autonomousDatabases"),
		resourceids.UserSpecifiedSegment("autonomousDatabaseName", "autonomousDatabaseName"),
		resourceids.StaticSegment("staticAutonomousDatabaseBackups", "autonomousDatabaseBackups", "autonomousDatabaseBackups"),
		resourceids.UserSpecifiedSegment("autonomousDatabaseBackupName", "autonomousDatabaseBackupName"),
	}
}

// String returns a human-readable description of this Autonomous Database Backup ID
func (id AutonomousDatabaseBackupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Autonomous Database Name: %q", id.AutonomousDatabaseName),
		fmt.Sprintf("Autonomous Database Backup Name: %q", id.AutonomousDatabaseBackupName),
	}
	return fmt.Sprintf("Autonomous Database Backup (%s)", strings.Join(components, "\n"))
}
