package tables

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = TableId{}

// TableId is a struct representing the Resource ID for a Table
type TableId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	TableName         string
}

// NewTableID returns a new TableId struct
func NewTableID(subscriptionId string, resourceGroupName string, workspaceName string, tableName string) TableId {
	return TableId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		TableName:         tableName,
	}
}

// ParseTableID parses 'input' into a TableId
func ParseTableID(input string) (*TableId, error) {
	parser := resourceids.NewParserFromResourceIdType(TableId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TableId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'workspaceName' was not found in the resource id %q", input)
	}

	if id.TableName, ok = parsed.Parsed["tableName"]; !ok {
		return nil, fmt.Errorf("the segment 'tableName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseTableIDInsensitively parses 'input' case-insensitively into a TableId
// note: this method should only be used for API response data and not user input
func ParseTableIDInsensitively(input string) (*TableId, error) {
	parser := resourceids.NewParserFromResourceIdType(TableId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TableId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'workspaceName' was not found in the resource id %q", input)
	}

	if id.TableName, ok = parsed.Parsed["tableName"]; !ok {
		return nil, fmt.Errorf("the segment 'tableName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateTableID checks that 'input' can be parsed as a Table ID
func ValidateTableID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTableID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Table ID
func (id TableId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/tables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.TableName)
}

// Segments returns a slice of Resource ID Segments which comprise this Table ID
func (id TableId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticTables", "tables", "tables"),
		resourceids.UserSpecifiedSegment("tableName", "tableValue"),
	}
}

// String returns a human-readable description of this Table ID
func (id TableId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Table Name: %q", id.TableName),
	}
	return fmt.Sprintf("Table (%s)", strings.Join(components, "\n"))
}
