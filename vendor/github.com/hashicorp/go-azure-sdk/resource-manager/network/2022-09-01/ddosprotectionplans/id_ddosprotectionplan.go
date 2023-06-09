package ddosprotectionplans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DdosProtectionPlanId{}

// DdosProtectionPlanId is a struct representing the Resource ID for a Ddos Protection Plan
type DdosProtectionPlanId struct {
	SubscriptionId         string
	ResourceGroupName      string
	DdosProtectionPlanName string
}

// NewDdosProtectionPlanID returns a new DdosProtectionPlanId struct
func NewDdosProtectionPlanID(subscriptionId string, resourceGroupName string, ddosProtectionPlanName string) DdosProtectionPlanId {
	return DdosProtectionPlanId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		DdosProtectionPlanName: ddosProtectionPlanName,
	}
}

// ParseDdosProtectionPlanID parses 'input' into a DdosProtectionPlanId
func ParseDdosProtectionPlanID(input string) (*DdosProtectionPlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(DdosProtectionPlanId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DdosProtectionPlanId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DdosProtectionPlanName, ok = parsed.Parsed["ddosProtectionPlanName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "ddosProtectionPlanName", *parsed)
	}

	return &id, nil
}

// ParseDdosProtectionPlanIDInsensitively parses 'input' case-insensitively into a DdosProtectionPlanId
// note: this method should only be used for API response data and not user input
func ParseDdosProtectionPlanIDInsensitively(input string) (*DdosProtectionPlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(DdosProtectionPlanId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DdosProtectionPlanId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DdosProtectionPlanName, ok = parsed.Parsed["ddosProtectionPlanName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "ddosProtectionPlanName", *parsed)
	}

	return &id, nil
}

// ValidateDdosProtectionPlanID checks that 'input' can be parsed as a Ddos Protection Plan ID
func ValidateDdosProtectionPlanID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDdosProtectionPlanID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Ddos Protection Plan ID
func (id DdosProtectionPlanId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/ddosProtectionPlans/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DdosProtectionPlanName)
}

// Segments returns a slice of Resource ID Segments which comprise this Ddos Protection Plan ID
func (id DdosProtectionPlanId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticDdosProtectionPlans", "ddosProtectionPlans", "ddosProtectionPlans"),
		resourceids.UserSpecifiedSegment("ddosProtectionPlanName", "ddosProtectionPlanValue"),
	}
}

// String returns a human-readable description of this Ddos Protection Plan ID
func (id DdosProtectionPlanId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Ddos Protection Plan Name: %q", id.DdosProtectionPlanName),
	}
	return fmt.Sprintf("Ddos Protection Plan (%s)", strings.Join(components, "\n"))
}
