package sqlvirtualmachines

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SqlVirtualMachineId{})
}

var _ resourceids.ResourceId = &SqlVirtualMachineId{}

// SqlVirtualMachineId is a struct representing the Resource ID for a Sql Virtual Machine
type SqlVirtualMachineId struct {
	SubscriptionId        string
	ResourceGroupName     string
	SqlVirtualMachineName string
}

// NewSqlVirtualMachineID returns a new SqlVirtualMachineId struct
func NewSqlVirtualMachineID(subscriptionId string, resourceGroupName string, sqlVirtualMachineName string) SqlVirtualMachineId {
	return SqlVirtualMachineId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		SqlVirtualMachineName: sqlVirtualMachineName,
	}
}

// ParseSqlVirtualMachineID parses 'input' into a SqlVirtualMachineId
func ParseSqlVirtualMachineID(input string) (*SqlVirtualMachineId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SqlVirtualMachineId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SqlVirtualMachineId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSqlVirtualMachineIDInsensitively parses 'input' case-insensitively into a SqlVirtualMachineId
// note: this method should only be used for API response data and not user input
func ParseSqlVirtualMachineIDInsensitively(input string) (*SqlVirtualMachineId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SqlVirtualMachineId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SqlVirtualMachineId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SqlVirtualMachineId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SqlVirtualMachineName, ok = input.Parsed["sqlVirtualMachineName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sqlVirtualMachineName", input)
	}

	return nil
}

// ValidateSqlVirtualMachineID checks that 'input' can be parsed as a Sql Virtual Machine ID
func ValidateSqlVirtualMachineID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSqlVirtualMachineID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Sql Virtual Machine ID
func (id SqlVirtualMachineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SqlVirtualMachineName)
}

// Segments returns a slice of Resource ID Segments which comprise this Sql Virtual Machine ID
func (id SqlVirtualMachineId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSqlVirtualMachine", "Microsoft.SqlVirtualMachine", "Microsoft.SqlVirtualMachine"),
		resourceids.StaticSegment("staticSqlVirtualMachines", "sqlVirtualMachines", "sqlVirtualMachines"),
		resourceids.UserSpecifiedSegment("sqlVirtualMachineName", "sqlVirtualMachineName"),
	}
}

// String returns a human-readable description of this Sql Virtual Machine ID
func (id SqlVirtualMachineId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Sql Virtual Machine Name: %q", id.SqlVirtualMachineName),
	}
	return fmt.Sprintf("Sql Virtual Machine (%s)", strings.Join(components, "\n"))
}
