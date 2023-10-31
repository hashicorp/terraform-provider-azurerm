package deletedservice

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DeletedServiceId{}

// DeletedServiceId is a struct representing the Resource ID for a Deleted Service
type DeletedServiceId struct {
	SubscriptionId     string
	LocationName       string
	DeletedServiceName string
}

// NewDeletedServiceID returns a new DeletedServiceId struct
func NewDeletedServiceID(subscriptionId string, locationName string, deletedServiceName string) DeletedServiceId {
	return DeletedServiceId{
		SubscriptionId:     subscriptionId,
		LocationName:       locationName,
		DeletedServiceName: deletedServiceName,
	}
}

// ParseDeletedServiceID parses 'input' into a DeletedServiceId
func ParseDeletedServiceID(input string) (*DeletedServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(DeletedServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DeletedServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.DeletedServiceName, ok = parsed.Parsed["deletedServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "deletedServiceName", *parsed)
	}

	return &id, nil
}

// ParseDeletedServiceIDInsensitively parses 'input' case-insensitively into a DeletedServiceId
// note: this method should only be used for API response data and not user input
func ParseDeletedServiceIDInsensitively(input string) (*DeletedServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(DeletedServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DeletedServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.DeletedServiceName, ok = parsed.Parsed["deletedServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "deletedServiceName", *parsed)
	}

	return &id, nil
}

// ValidateDeletedServiceID checks that 'input' can be parsed as a Deleted Service ID
func ValidateDeletedServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeletedServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deleted Service ID
func (id DeletedServiceId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.ApiManagement/locations/%s/deletedServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.DeletedServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Deleted Service ID
func (id DeletedServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationValue"),
		resourceids.StaticSegment("staticDeletedServices", "deletedServices", "deletedServices"),
		resourceids.UserSpecifiedSegment("deletedServiceName", "deletedServiceValue"),
	}
}

// String returns a human-readable description of this Deleted Service ID
func (id DeletedServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Deleted Service Name: %q", id.DeletedServiceName),
	}
	return fmt.Sprintf("Deleted Service (%s)", strings.Join(components, "\n"))
}
