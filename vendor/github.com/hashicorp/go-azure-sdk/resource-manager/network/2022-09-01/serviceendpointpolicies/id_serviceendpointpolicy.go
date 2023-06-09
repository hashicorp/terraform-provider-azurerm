package serviceendpointpolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ServiceEndpointPolicyId{}

// ServiceEndpointPolicyId is a struct representing the Resource ID for a Service Endpoint Policy
type ServiceEndpointPolicyId struct {
	SubscriptionId            string
	ResourceGroupName         string
	ServiceEndpointPolicyName string
}

// NewServiceEndpointPolicyID returns a new ServiceEndpointPolicyId struct
func NewServiceEndpointPolicyID(subscriptionId string, resourceGroupName string, serviceEndpointPolicyName string) ServiceEndpointPolicyId {
	return ServiceEndpointPolicyId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		ServiceEndpointPolicyName: serviceEndpointPolicyName,
	}
}

// ParseServiceEndpointPolicyID parses 'input' into a ServiceEndpointPolicyId
func ParseServiceEndpointPolicyID(input string) (*ServiceEndpointPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(ServiceEndpointPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ServiceEndpointPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceEndpointPolicyName, ok = parsed.Parsed["serviceEndpointPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceEndpointPolicyName", *parsed)
	}

	return &id, nil
}

// ParseServiceEndpointPolicyIDInsensitively parses 'input' case-insensitively into a ServiceEndpointPolicyId
// note: this method should only be used for API response data and not user input
func ParseServiceEndpointPolicyIDInsensitively(input string) (*ServiceEndpointPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(ServiceEndpointPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ServiceEndpointPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceEndpointPolicyName, ok = parsed.Parsed["serviceEndpointPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceEndpointPolicyName", *parsed)
	}

	return &id, nil
}

// ValidateServiceEndpointPolicyID checks that 'input' can be parsed as a Service Endpoint Policy ID
func ValidateServiceEndpointPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseServiceEndpointPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Service Endpoint Policy ID
func (id ServiceEndpointPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/serviceEndpointPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceEndpointPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Service Endpoint Policy ID
func (id ServiceEndpointPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticServiceEndpointPolicies", "serviceEndpointPolicies", "serviceEndpointPolicies"),
		resourceids.UserSpecifiedSegment("serviceEndpointPolicyName", "serviceEndpointPolicyValue"),
	}
}

// String returns a human-readable description of this Service Endpoint Policy ID
func (id ServiceEndpointPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Endpoint Policy Name: %q", id.ServiceEndpointPolicyName),
	}
	return fmt.Sprintf("Service Endpoint Policy (%s)", strings.Join(components, "\n"))
}
