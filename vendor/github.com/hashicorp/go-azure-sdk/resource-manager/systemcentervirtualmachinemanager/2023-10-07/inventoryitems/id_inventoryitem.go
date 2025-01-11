package inventoryitems

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&InventoryItemId{})
}

var _ resourceids.ResourceId = &InventoryItemId{}

// InventoryItemId is a struct representing the Resource ID for a Inventory Item
type InventoryItemId struct {
	SubscriptionId    string
	ResourceGroupName string
	VmmServerName     string
	InventoryItemName string
}

// NewInventoryItemID returns a new InventoryItemId struct
func NewInventoryItemID(subscriptionId string, resourceGroupName string, vmmServerName string, inventoryItemName string) InventoryItemId {
	return InventoryItemId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VmmServerName:     vmmServerName,
		InventoryItemName: inventoryItemName,
	}
}

// ParseInventoryItemID parses 'input' into a InventoryItemId
func ParseInventoryItemID(input string) (*InventoryItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InventoryItemId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InventoryItemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseInventoryItemIDInsensitively parses 'input' case-insensitively into a InventoryItemId
// note: this method should only be used for API response data and not user input
func ParseInventoryItemIDInsensitively(input string) (*InventoryItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InventoryItemId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InventoryItemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *InventoryItemId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VmmServerName, ok = input.Parsed["vmmServerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vmmServerName", input)
	}

	if id.InventoryItemName, ok = input.Parsed["inventoryItemName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "inventoryItemName", input)
	}

	return nil
}

// ValidateInventoryItemID checks that 'input' can be parsed as a Inventory Item ID
func ValidateInventoryItemID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseInventoryItemID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Inventory Item ID
func (id InventoryItemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ScVmm/vmmServers/%s/inventoryItems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VmmServerName, id.InventoryItemName)
}

// Segments returns a slice of Resource ID Segments which comprise this Inventory Item ID
func (id InventoryItemId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftScVmm", "Microsoft.ScVmm", "Microsoft.ScVmm"),
		resourceids.StaticSegment("staticVmmServers", "vmmServers", "vmmServers"),
		resourceids.UserSpecifiedSegment("vmmServerName", "vmmServerName"),
		resourceids.StaticSegment("staticInventoryItems", "inventoryItems", "inventoryItems"),
		resourceids.UserSpecifiedSegment("inventoryItemName", "inventoryItemName"),
	}
}

// String returns a human-readable description of this Inventory Item ID
func (id InventoryItemId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vmm Server Name: %q", id.VmmServerName),
		fmt.Sprintf("Inventory Item Name: %q", id.InventoryItemName),
	}
	return fmt.Sprintf("Inventory Item (%s)", strings.Join(components, "\n"))
}
