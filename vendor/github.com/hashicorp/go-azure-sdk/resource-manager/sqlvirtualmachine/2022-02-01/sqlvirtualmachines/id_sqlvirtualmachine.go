package sqlvirtualmachines

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SqlVirtualMachineId{}

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
	parser := resourceids.NewParserFromResourceIdType(SqlVirtualMachineId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SqlVirtualMachineId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.SqlVirtualMachineName, ok = parsed.Parsed["sqlVirtualMachineName"]; !ok {
		return nil, fmt.Errorf("the segment 'sqlVirtualMachineName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseSqlVirtualMachineIDInsensitively parses 'input' case-insensitively into a SqlVirtualMachineId
// note: this method should only be used for API response data and not user input
func ParseSqlVirtualMachineIDInsensitively(input string) (*SqlVirtualMachineId, error) {
	parser := resourceids.NewParserFromResourceIdType(SqlVirtualMachineId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SqlVirtualMachineId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.SqlVirtualMachineName, ok = parsed.Parsed["sqlVirtualMachineName"]; !ok {
		return nil, fmt.Errorf("the segment 'sqlVirtualMachineName' was not found in the resource id %q", input)
	}

	return &id, nil
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
		resourceids.UserSpecifiedSegment("sqlVirtualMachineName", "sqlVirtualMachineValue"),
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
