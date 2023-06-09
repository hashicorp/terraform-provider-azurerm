package serviceendpointpolicydefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ServiceEndpointPolicyDefinitionId{}

// ServiceEndpointPolicyDefinitionId is a struct representing the Resource ID for a Service Endpoint Policy Definition
type ServiceEndpointPolicyDefinitionId struct {
	SubscriptionId                      string
	ResourceGroupName                   string
	ServiceEndpointPolicyName           string
	ServiceEndpointPolicyDefinitionName string
}

// NewServiceEndpointPolicyDefinitionID returns a new ServiceEndpointPolicyDefinitionId struct
func NewServiceEndpointPolicyDefinitionID(subscriptionId string, resourceGroupName string, serviceEndpointPolicyName string, serviceEndpointPolicyDefinitionName string) ServiceEndpointPolicyDefinitionId {
	return ServiceEndpointPolicyDefinitionId{
		SubscriptionId:                      subscriptionId,
		ResourceGroupName:                   resourceGroupName,
		ServiceEndpointPolicyName:           serviceEndpointPolicyName,
		ServiceEndpointPolicyDefinitionName: serviceEndpointPolicyDefinitionName,
	}
}

// ParseServiceEndpointPolicyDefinitionID parses 'input' into a ServiceEndpointPolicyDefinitionId
func ParseServiceEndpointPolicyDefinitionID(input string) (*ServiceEndpointPolicyDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ServiceEndpointPolicyDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ServiceEndpointPolicyDefinitionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceEndpointPolicyName, ok = parsed.Parsed["serviceEndpointPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceEndpointPolicyName", *parsed)
	}

	if id.ServiceEndpointPolicyDefinitionName, ok = parsed.Parsed["serviceEndpointPolicyDefinitionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceEndpointPolicyDefinitionName", *parsed)
	}

	return &id, nil
}

// ParseServiceEndpointPolicyDefinitionIDInsensitively parses 'input' case-insensitively into a ServiceEndpointPolicyDefinitionId
// note: this method should only be used for API response data and not user input
func ParseServiceEndpointPolicyDefinitionIDInsensitively(input string) (*ServiceEndpointPolicyDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ServiceEndpointPolicyDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ServiceEndpointPolicyDefinitionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceEndpointPolicyName, ok = parsed.Parsed["serviceEndpointPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceEndpointPolicyName", *parsed)
	}

	if id.ServiceEndpointPolicyDefinitionName, ok = parsed.Parsed["serviceEndpointPolicyDefinitionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceEndpointPolicyDefinitionName", *parsed)
	}

	return &id, nil
}

// ValidateServiceEndpointPolicyDefinitionID checks that 'input' can be parsed as a Service Endpoint Policy Definition ID
func ValidateServiceEndpointPolicyDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseServiceEndpointPolicyDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Service Endpoint Policy Definition ID
func (id ServiceEndpointPolicyDefinitionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/serviceEndpointPolicies/%s/serviceEndpointPolicyDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceEndpointPolicyName, id.ServiceEndpointPolicyDefinitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Service Endpoint Policy Definition ID
func (id ServiceEndpointPolicyDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticServiceEndpointPolicies", "serviceEndpointPolicies", "serviceEndpointPolicies"),
		resourceids.UserSpecifiedSegment("serviceEndpointPolicyName", "serviceEndpointPolicyValue"),
		resourceids.StaticSegment("staticServiceEndpointPolicyDefinitions", "serviceEndpointPolicyDefinitions", "serviceEndpointPolicyDefinitions"),
		resourceids.UserSpecifiedSegment("serviceEndpointPolicyDefinitionName", "serviceEndpointPolicyDefinitionValue"),
	}
}

// String returns a human-readable description of this Service Endpoint Policy Definition ID
func (id ServiceEndpointPolicyDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Endpoint Policy Name: %q", id.ServiceEndpointPolicyName),
		fmt.Sprintf("Service Endpoint Policy Definition Name: %q", id.ServiceEndpointPolicyDefinitionName),
	}
	return fmt.Sprintf("Service Endpoint Policy Definition (%s)", strings.Join(components, "\n"))
}
