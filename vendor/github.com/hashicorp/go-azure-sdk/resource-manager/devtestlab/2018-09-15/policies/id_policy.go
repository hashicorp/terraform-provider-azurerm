package policies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PolicyId{}

// PolicyId is a struct representing the Resource ID for a Policy
type PolicyId struct {
	SubscriptionId    string
	ResourceGroupName string
	LabName           string
	PolicySetName     string
	PolicyName        string
}

// NewPolicyID returns a new PolicyId struct
func NewPolicyID(subscriptionId string, resourceGroupName string, labName string, policySetName string, policyName string) PolicyId {
	return PolicyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LabName:           labName,
		PolicySetName:     policySetName,
		PolicyName:        policyName,
	}
}

// ParsePolicyID parses 'input' into a PolicyId
func ParsePolicyID(input string) (*PolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(PolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PolicyId{}

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

	if id.PolicyName, ok = parsed.Parsed["policyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "policyName", *parsed)
	}

	return &id, nil
}

// ParsePolicyIDInsensitively parses 'input' case-insensitively into a PolicyId
// note: this method should only be used for API response data and not user input
func ParsePolicyIDInsensitively(input string) (*PolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(PolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PolicyId{}

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

	if id.PolicyName, ok = parsed.Parsed["policyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "policyName", *parsed)
	}

	return &id, nil
}

// ValidatePolicyID checks that 'input' can be parsed as a Policy ID
func ValidatePolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Policy ID
func (id PolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevTestLab/labs/%s/policySets/%s/policies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LabName, id.PolicySetName, id.PolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Policy ID
func (id PolicyId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticPolicies", "policies", "policies"),
		resourceids.UserSpecifiedSegment("policyName", "policyValue"),
	}
}

// String returns a human-readable description of this Policy ID
func (id PolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Lab Name: %q", id.LabName),
		fmt.Sprintf("Policy Set Name: %q", id.PolicySetName),
		fmt.Sprintf("Policy Name: %q", id.PolicyName),
	}
	return fmt.Sprintf("Policy (%s)", strings.Join(components, "\n"))
}
