// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SqlManagedInstanceId{}

// SqlManagedInstanceId is a struct representing the Resource ID for a Sql Managed Instance
type SqlManagedInstanceId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ManagedInstanceName string
}

// NewSqlManagedInstanceID returns a new SqlManagedInstanceId struct
func NewSqlManagedInstanceID(subscriptionId string, resourceGroupName string, managedInstanceName string) SqlManagedInstanceId {
	return SqlManagedInstanceId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ManagedInstanceName: managedInstanceName,
	}
}

// ParseSqlManagedInstanceID parses 'input' into a SqlManagedInstanceId
func ParseSqlManagedInstanceID(input string) (*SqlManagedInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(SqlManagedInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SqlManagedInstanceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedInstanceName, ok = parsed.Parsed["managedInstanceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedInstanceName", *parsed)
	}

	return &id, nil
}

// ParseSqlManagedInstanceIDInsensitively parses 'input' case-insensitively into a SqlManagedInstanceId
// note: this method should only be used for API response data and not user input
func ParseSqlManagedInstanceIDInsensitively(input string) (*SqlManagedInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(SqlManagedInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SqlManagedInstanceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedInstanceName, ok = parsed.Parsed["managedInstanceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedInstanceName", *parsed)
	}

	return &id, nil
}

// ValidateSqlManagedInstanceID checks that 'input' can be parsed as a Sql Managed Instance ID
func ValidateSqlManagedInstanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSqlManagedInstanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Sql Managed Instance ID
func (id SqlManagedInstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/managedInstances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedInstanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Sql Managed Instance ID
func (id SqlManagedInstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSql", "Microsoft.Sql", "Microsoft.Sql"),
		resourceids.StaticSegment("staticManagedInstances", "managedInstances", "managedInstances"),
		resourceids.UserSpecifiedSegment("managedInstanceName", "managedInstanceValue"),
	}
}

// String returns a human-readable description of this Sql Managed Instance ID
func (id SqlManagedInstanceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Instance Name: %q", id.ManagedInstanceName),
	}
	return fmt.Sprintf("Managed Instance (%s)", strings.Join(components, "\n"))
}
