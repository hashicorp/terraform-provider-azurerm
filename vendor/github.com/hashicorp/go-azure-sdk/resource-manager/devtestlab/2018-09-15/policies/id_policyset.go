package policies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PolicySetId{}

// PolicySetId is a struct representing the Resource ID for a Policy Set
type PolicySetId struct {
	SubscriptionId    string
	ResourceGroupName string
	LabName           string
	PolicySetName     string
}

// NewPolicySetID returns a new PolicySetId struct
func NewPolicySetID(subscriptionId string, resourceGroupName string, labName string, policySetName string) PolicySetId {
	return PolicySetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LabName:           labName,
		PolicySetName:     policySetName,
	}
}

// ParsePolicySetID parses 'input' into a PolicySetId
func ParsePolicySetID(input string) (*PolicySetId, error) {
	parser := resourceids.NewParserFromResourceIdType(PolicySetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PolicySetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LabName, ok = parsed.Parsed["labName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "labName", *parsed)
	}

	if id.PolicySetName, ok = parsed.Parsed["policySetName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "policySetName", *parsed)
	}

	return &id, nil
}

// ParsePolicySetIDInsensitively parses 'input' case-insensitively into a PolicySetId
// note: this method should only be used for API response data and not user input
func ParsePolicySetIDInsensitively(input string) (*PolicySetId, error) {
	parser := resourceids.NewParserFromResourceIdType(PolicySetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PolicySetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LabName, ok = parsed.Parsed["labName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "labName", *parsed)
	}

	if id.PolicySetName, ok = parsed.Parsed["policySetName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "policySetName", *parsed)
	}

	return &id, nil
}

// ValidatePolicySetID checks that 'input' can be parsed as a Policy Set ID
func ValidatePolicySetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePolicySetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Policy Set ID
func (id PolicySetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevTestLab/labs/%s/policySets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LabName, id.PolicySetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Policy Set ID
func (id PolicySetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevTestLab", "Microsoft.DevTestLab", "Microsoft.DevTestLab"),
		resourceids.StaticSegment("staticLabs", "labs", "labs"),
		resourceids.UserSpecifiedSegment("labName", "labValue"),
		resourceids.StaticSegment("staticPolicySets", "policySets", "policySets"),
		resourceids.UserSpecifiedSegment("policySetName", "policySetValue"),
	}
}

// String returns a human-readable description of this Policy Set ID
func (id PolicySetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Lab Name: %q", id.LabName),
		fmt.Sprintf("Policy Set Name: %q", id.PolicySetName),
	}
	return fmt.Sprintf("Policy Set (%s)", strings.Join(components, "\n"))
}
