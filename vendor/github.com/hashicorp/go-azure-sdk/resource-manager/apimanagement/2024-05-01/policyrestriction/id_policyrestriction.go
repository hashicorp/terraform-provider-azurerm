package policyrestriction

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PolicyRestrictionId{})
}

var _ resourceids.ResourceId = &PolicyRestrictionId{}

// PolicyRestrictionId is a struct representing the Resource ID for a Policy Restriction
type PolicyRestrictionId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ServiceName         string
	PolicyRestrictionId string
}

// NewPolicyRestrictionID returns a new PolicyRestrictionId struct
func NewPolicyRestrictionID(subscriptionId string, resourceGroupName string, serviceName string, policyRestrictionId string) PolicyRestrictionId {
	return PolicyRestrictionId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ServiceName:         serviceName,
		PolicyRestrictionId: policyRestrictionId,
	}
}

// ParsePolicyRestrictionID parses 'input' into a PolicyRestrictionId
func ParsePolicyRestrictionID(input string) (*PolicyRestrictionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PolicyRestrictionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PolicyRestrictionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePolicyRestrictionIDInsensitively parses 'input' case-insensitively into a PolicyRestrictionId
// note: this method should only be used for API response data and not user input
func ParsePolicyRestrictionIDInsensitively(input string) (*PolicyRestrictionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PolicyRestrictionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PolicyRestrictionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PolicyRestrictionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.PolicyRestrictionId, ok = input.Parsed["policyRestrictionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "policyRestrictionId", input)
	}

	return nil
}

// ValidatePolicyRestrictionID checks that 'input' can be parsed as a Policy Restriction ID
func ValidatePolicyRestrictionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePolicyRestrictionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Policy Restriction ID
func (id PolicyRestrictionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/policyRestrictions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.PolicyRestrictionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Policy Restriction ID
func (id PolicyRestrictionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticPolicyRestrictions", "policyRestrictions", "policyRestrictions"),
		resourceids.UserSpecifiedSegment("policyRestrictionId", "policyRestrictionId"),
	}
}

// String returns a human-readable description of this Policy Restriction ID
func (id PolicyRestrictionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Policy Restriction: %q", id.PolicyRestrictionId),
	}
	return fmt.Sprintf("Policy Restriction (%s)", strings.Join(components, "\n"))
}
