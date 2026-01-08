package policyfragment

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PolicyFragmentId{})
}

var _ resourceids.ResourceId = &PolicyFragmentId{}

// PolicyFragmentId is a struct representing the Resource ID for a Policy Fragment
type PolicyFragmentId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ServiceName        string
	PolicyFragmentName string
}

// NewPolicyFragmentID returns a new PolicyFragmentId struct
func NewPolicyFragmentID(subscriptionId string, resourceGroupName string, serviceName string, policyFragmentName string) PolicyFragmentId {
	return PolicyFragmentId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ServiceName:        serviceName,
		PolicyFragmentName: policyFragmentName,
	}
}

// ParsePolicyFragmentID parses 'input' into a PolicyFragmentId
func ParsePolicyFragmentID(input string) (*PolicyFragmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PolicyFragmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PolicyFragmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePolicyFragmentIDInsensitively parses 'input' case-insensitively into a PolicyFragmentId
// note: this method should only be used for API response data and not user input
func ParsePolicyFragmentIDInsensitively(input string) (*PolicyFragmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PolicyFragmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PolicyFragmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PolicyFragmentId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.PolicyFragmentName, ok = input.Parsed["policyFragmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "policyFragmentName", input)
	}

	return nil
}

// ValidatePolicyFragmentID checks that 'input' can be parsed as a Policy Fragment ID
func ValidatePolicyFragmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePolicyFragmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Policy Fragment ID
func (id PolicyFragmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/policyFragments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.PolicyFragmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Policy Fragment ID
func (id PolicyFragmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticPolicyFragments", "policyFragments", "policyFragments"),
		resourceids.UserSpecifiedSegment("policyFragmentName", "policyFragmentName"),
	}
}

// String returns a human-readable description of this Policy Fragment ID
func (id PolicyFragmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Policy Fragment Name: %q", id.PolicyFragmentName),
	}
	return fmt.Sprintf("Policy Fragment (%s)", strings.Join(components, "\n"))
}
