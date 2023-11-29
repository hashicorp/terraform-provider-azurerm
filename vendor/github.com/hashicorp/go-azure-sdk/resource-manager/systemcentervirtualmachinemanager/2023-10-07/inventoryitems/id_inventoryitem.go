package inventoryitems

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = InventoryItemId{}

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
	parser := resourceids.NewParserFromResourceIdType(InventoryItemId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := InventoryItemId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VmmServerName, ok = parsed.Parsed["vmmServerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vmmServerName", *parsed)
	}

	if id.InventoryItemName, ok = parsed.Parsed["inventoryItemName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "inventoryItemName", *parsed)
	}

	return &id, nil
}

// ParseInventoryItemIDInsensitively parses 'input' case-insensitively into a InventoryItemId
// note: this method should only be used for API response data and not user input
func ParseInventoryItemIDInsensitively(input string) (*InventoryItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(InventoryItemId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := InventoryItemId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VmmServerName, ok = parsed.Parsed["vmmServerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vmmServerName", *parsed)
	}

	if id.InventoryItemName, ok = parsed.Parsed["inventoryItemName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "inventoryItemName", *parsed)
	}

	return &id, nil
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
		resourceids.UserSpecifiedSegment("vmmServerName", "vmmServerValue"),
		resourceids.StaticSegment("staticInventoryItems", "inventoryItems", "inventoryItems"),
		resourceids.UserSpecifiedSegment("inventoryItemName", "inventoryItemValue"),
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
