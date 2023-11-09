package ddoscustompolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DdosCustomPolicyId{}

// DdosCustomPolicyId is a struct representing the Resource ID for a Ddos Custom Policy
type DdosCustomPolicyId struct {
	SubscriptionId       string
	ResourceGroupName    string
	DdosCustomPolicyName string
}

// NewDdosCustomPolicyID returns a new DdosCustomPolicyId struct
func NewDdosCustomPolicyID(subscriptionId string, resourceGroupName string, ddosCustomPolicyName string) DdosCustomPolicyId {
	return DdosCustomPolicyId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		DdosCustomPolicyName: ddosCustomPolicyName,
	}
}

// ParseDdosCustomPolicyID parses 'input' into a DdosCustomPolicyId
func ParseDdosCustomPolicyID(input string) (*DdosCustomPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(DdosCustomPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DdosCustomPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DdosCustomPolicyName, ok = parsed.Parsed["ddosCustomPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "ddosCustomPolicyName", *parsed)
	}

	return &id, nil
}

// ParseDdosCustomPolicyIDInsensitively parses 'input' case-insensitively into a DdosCustomPolicyId
// note: this method should only be used for API response data and not user input
func ParseDdosCustomPolicyIDInsensitively(input string) (*DdosCustomPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(DdosCustomPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DdosCustomPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DdosCustomPolicyName, ok = parsed.Parsed["ddosCustomPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "ddosCustomPolicyName", *parsed)
	}

	return &id, nil
}

// ValidateDdosCustomPolicyID checks that 'input' can be parsed as a Ddos Custom Policy ID
func ValidateDdosCustomPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDdosCustomPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Ddos Custom Policy ID
func (id DdosCustomPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/ddosCustomPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DdosCustomPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Ddos Custom Policy ID
func (id DdosCustomPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticDdosCustomPolicies", "ddosCustomPolicies", "ddosCustomPolicies"),
		resourceids.UserSpecifiedSegment("ddosCustomPolicyName", "ddosCustomPolicyValue"),
	}
}

// String returns a human-readable description of this Ddos Custom Policy ID
func (id DdosCustomPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Ddos Custom Policy Name: %q", id.DdosCustomPolicyName),
	}
	return fmt.Sprintf("Ddos Custom Policy (%s)", strings.Join(components, "\n"))
}
