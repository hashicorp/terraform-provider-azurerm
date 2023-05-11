package communicationservices

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CommunicationServiceId{}

// CommunicationServiceId is a struct representing the Resource ID for a Communication Service
type CommunicationServiceId struct {
	SubscriptionId           string
	ResourceGroupName        string
	CommunicationServiceName string
}

// NewCommunicationServiceID returns a new CommunicationServiceId struct
func NewCommunicationServiceID(subscriptionId string, resourceGroupName string, communicationServiceName string) CommunicationServiceId {
	return CommunicationServiceId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		CommunicationServiceName: communicationServiceName,
	}
}

// ParseCommunicationServiceID parses 'input' into a CommunicationServiceId
func ParseCommunicationServiceID(input string) (*CommunicationServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(CommunicationServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CommunicationServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CommunicationServiceName, ok = parsed.Parsed["communicationServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "communicationServiceName", *parsed)
	}

	return &id, nil
}

// ParseCommunicationServiceIDInsensitively parses 'input' case-insensitively into a CommunicationServiceId
// note: this method should only be used for API response data and not user input
func ParseCommunicationServiceIDInsensitively(input string) (*CommunicationServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(CommunicationServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CommunicationServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CommunicationServiceName, ok = parsed.Parsed["communicationServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "communicationServiceName", *parsed)
	}

	return &id, nil
}

// ValidateCommunicationServiceID checks that 'input' can be parsed as a Communication Service ID
func ValidateCommunicationServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCommunicationServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Communication Service ID
func (id CommunicationServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Communication/communicationServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CommunicationServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Communication Service ID
func (id CommunicationServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCommunication", "Microsoft.Communication", "Microsoft.Communication"),
		resourceids.StaticSegment("staticCommunicationServices", "communicationServices", "communicationServices"),
		resourceids.UserSpecifiedSegment("communicationServiceName", "communicationServiceValue"),
	}
}

// String returns a human-readable description of this Communication Service ID
func (id CommunicationServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Communication Service Name: %q", id.CommunicationServiceName),
	}
	return fmt.Sprintf("Communication Service (%s)", strings.Join(components, "\n"))
}
