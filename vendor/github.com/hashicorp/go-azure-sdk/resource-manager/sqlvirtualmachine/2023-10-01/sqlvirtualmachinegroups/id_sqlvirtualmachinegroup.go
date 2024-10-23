package sqlvirtualmachinegroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SqlVirtualMachineGroupId{})
}

var _ resourceids.ResourceId = &SqlVirtualMachineGroupId{}

// SqlVirtualMachineGroupId is a struct representing the Resource ID for a Sql Virtual Machine Group
type SqlVirtualMachineGroupId struct {
	SubscriptionId             string
	ResourceGroupName          string
	SqlVirtualMachineGroupName string
}

// NewSqlVirtualMachineGroupID returns a new SqlVirtualMachineGroupId struct
func NewSqlVirtualMachineGroupID(subscriptionId string, resourceGroupName string, sqlVirtualMachineGroupName string) SqlVirtualMachineGroupId {
	return SqlVirtualMachineGroupId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		SqlVirtualMachineGroupName: sqlVirtualMachineGroupName,
	}
}

// ParseSqlVirtualMachineGroupID parses 'input' into a SqlVirtualMachineGroupId
func ParseSqlVirtualMachineGroupID(input string) (*SqlVirtualMachineGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SqlVirtualMachineGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SqlVirtualMachineGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSqlVirtualMachineGroupIDInsensitively parses 'input' case-insensitively into a SqlVirtualMachineGroupId
// note: this method should only be used for API response data and not user input
func ParseSqlVirtualMachineGroupIDInsensitively(input string) (*SqlVirtualMachineGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SqlVirtualMachineGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SqlVirtualMachineGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SqlVirtualMachineGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SqlVirtualMachineGroupName, ok = input.Parsed["sqlVirtualMachineGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sqlVirtualMachineGroupName", input)
	}

	return nil
}

// ValidateSqlVirtualMachineGroupID checks that 'input' can be parsed as a Sql Virtual Machine Group ID
func ValidateSqlVirtualMachineGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSqlVirtualMachineGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Sql Virtual Machine Group ID
func (id SqlVirtualMachineGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SqlVirtualMachineGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Sql Virtual Machine Group ID
func (id SqlVirtualMachineGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSqlVirtualMachine", "Microsoft.SqlVirtualMachine", "Microsoft.SqlVirtualMachine"),
		resourceids.StaticSegment("staticSqlVirtualMachineGroups", "sqlVirtualMachineGroups", "sqlVirtualMachineGroups"),
		resourceids.UserSpecifiedSegment("sqlVirtualMachineGroupName", "sqlVirtualMachineGroupName"),
	}
}

// String returns a human-readable description of this Sql Virtual Machine Group ID
func (id SqlVirtualMachineGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Sql Virtual Machine Group Name: %q", id.SqlVirtualMachineGroupName),
	}
	return fmt.Sprintf("Sql Virtual Machine Group (%s)", strings.Join(components, "\n"))
}
