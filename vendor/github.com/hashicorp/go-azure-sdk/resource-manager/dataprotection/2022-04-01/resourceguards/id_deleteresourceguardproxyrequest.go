package resourceguards

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DeleteResourceGuardProxyRequestId{}

// DeleteResourceGuardProxyRequestId is a struct representing the Resource ID for a Delete Resource Guard Proxy Request
type DeleteResourceGuardProxyRequestId struct {
	SubscriptionId                      string
	ResourceGroupName                   string
	ResourceGuardName                   string
	DeleteResourceGuardProxyRequestName string
}

// NewDeleteResourceGuardProxyRequestID returns a new DeleteResourceGuardProxyRequestId struct
func NewDeleteResourceGuardProxyRequestID(subscriptionId string, resourceGroupName string, resourceGuardName string, deleteResourceGuardProxyRequestName string) DeleteResourceGuardProxyRequestId {
	return DeleteResourceGuardProxyRequestId{
		SubscriptionId:                      subscriptionId,
		ResourceGroupName:                   resourceGroupName,
		ResourceGuardName:                   resourceGuardName,
		DeleteResourceGuardProxyRequestName: deleteResourceGuardProxyRequestName,
	}
}

// ParseDeleteResourceGuardProxyRequestID parses 'input' into a DeleteResourceGuardProxyRequestId
func ParseDeleteResourceGuardProxyRequestID(input string) (*DeleteResourceGuardProxyRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(DeleteResourceGuardProxyRequestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DeleteResourceGuardProxyRequestId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ResourceGuardName, ok = parsed.Parsed["resourceGuardName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGuardName", *parsed)
	}

	if id.DeleteResourceGuardProxyRequestName, ok = parsed.Parsed["deleteResourceGuardProxyRequestName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "deleteResourceGuardProxyRequestName", *parsed)
	}

	return &id, nil
}

// ParseDeleteResourceGuardProxyRequestIDInsensitively parses 'input' case-insensitively into a DeleteResourceGuardProxyRequestId
// note: this method should only be used for API response data and not user input
func ParseDeleteResourceGuardProxyRequestIDInsensitively(input string) (*DeleteResourceGuardProxyRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(DeleteResourceGuardProxyRequestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DeleteResourceGuardProxyRequestId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ResourceGuardName, ok = parsed.Parsed["resourceGuardName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGuardName", *parsed)
	}

	if id.DeleteResourceGuardProxyRequestName, ok = parsed.Parsed["deleteResourceGuardProxyRequestName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "deleteResourceGuardProxyRequestName", *parsed)
	}

	return &id, nil
}

// ValidateDeleteResourceGuardProxyRequestID checks that 'input' can be parsed as a Delete Resource Guard Proxy Request ID
func ValidateDeleteResourceGuardProxyRequestID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeleteResourceGuardProxyRequestID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Delete Resource Guard Proxy Request ID
func (id DeleteResourceGuardProxyRequestId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/resourceGuards/%s/deleteResourceGuardProxyRequests/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceGuardName, id.DeleteResourceGuardProxyRequestName)
}

// Segments returns a slice of Resource ID Segments which comprise this Delete Resource Guard Proxy Request ID
func (id DeleteResourceGuardProxyRequestId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticResourceGuards", "resourceGuards", "resourceGuards"),
		resourceids.UserSpecifiedSegment("resourceGuardName", "resourceGuardValue"),
		resourceids.StaticSegment("staticDeleteResourceGuardProxyRequests", "deleteResourceGuardProxyRequests", "deleteResourceGuardProxyRequests"),
		resourceids.UserSpecifiedSegment("deleteResourceGuardProxyRequestName", "deleteResourceGuardProxyRequestValue"),
	}
}

// String returns a human-readable description of this Delete Resource Guard Proxy Request ID
func (id DeleteResourceGuardProxyRequestId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Guard Name: %q", id.ResourceGuardName),
		fmt.Sprintf("Delete Resource Guard Proxy Request Name: %q", id.DeleteResourceGuardProxyRequestName),
	}
	return fmt.Sprintf("Delete Resource Guard Proxy Request (%s)", strings.Join(components, "\n"))
}
