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
	recaser.RegisterResourceId(&TriggerId{})
}

var _ resourceids.ResourceId = &TriggerId{}

// TriggerId is a struct representing the Resource ID for a Trigger
type TriggerId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	SqlDatabaseName     string
	ContainerName       string
	TriggerName         string
}

// NewTriggerID returns a new TriggerId struct
func NewTriggerID(subscriptionId string, resourceGroupName string, databaseAccountName string, sqlDatabaseName string, containerName string, triggerName string) TriggerId {
	return TriggerId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		SqlDatabaseName:     sqlDatabaseName,
		ContainerName:       containerName,
		TriggerName:         triggerName,
	}
}

// ParseTriggerID parses 'input' into a TriggerId
func ParseTriggerID(input string) (*TriggerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TriggerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TriggerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTriggerIDInsensitively parses 'input' case-insensitively into a TriggerId
// note: this method should only be used for API response data and not user input
func ParseTriggerIDInsensitively(input string) (*TriggerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TriggerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TriggerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TriggerId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.SqlDatabaseName, ok = input.Parsed["sqlDatabaseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sqlDatabaseName", input)
	}

	if id.ContainerName, ok = input.Parsed["containerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "containerName", input)
	}

	if id.TriggerName, ok = input.Parsed["triggerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "triggerName", input)
	}

	return nil
}

// ValidateTriggerID checks that 'input' can be parsed as a Trigger ID
func ValidateTriggerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTriggerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Trigger ID
func (id TriggerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/sqlDatabases/%s/containers/%s/triggers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.TriggerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Trigger ID
func (id TriggerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountName"),
		resourceids.StaticSegment("staticSqlDatabases", "sqlDatabases", "sqlDatabases"),
		resourceids.UserSpecifiedSegment("sqlDatabaseName", "sqlDatabaseName"),
		resourceids.StaticSegment("staticContainers", "containers", "containers"),
		resourceids.UserSpecifiedSegment("containerName", "containerName"),
		resourceids.StaticSegment("staticTriggers", "triggers", "triggers"),
		resourceids.UserSpecifiedSegment("triggerName", "triggerName"),
	}
}

// String returns a human-readable description of this Trigger ID
func (id TriggerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Sql Database Name: %q", id.SqlDatabaseName),
		fmt.Sprintf("Container Name: %q", id.ContainerName),
		fmt.Sprintf("Trigger Name: %q", id.TriggerName),
	}
	return fmt.Sprintf("Trigger (%s)", strings.Join(components, "\n"))
}
