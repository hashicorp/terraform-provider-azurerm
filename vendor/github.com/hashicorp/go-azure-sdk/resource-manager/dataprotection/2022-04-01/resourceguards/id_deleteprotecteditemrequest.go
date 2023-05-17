package resourceguards

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DeleteProtectedItemRequestId{}

// DeleteProtectedItemRequestId is a struct representing the Resource ID for a Delete Protected Item Request
type DeleteProtectedItemRequestId struct {
	SubscriptionId                 string
	ResourceGroupName              string
	ResourceGuardName              string
	DeleteProtectedItemRequestName string
}

// NewDeleteProtectedItemRequestID returns a new DeleteProtectedItemRequestId struct
func NewDeleteProtectedItemRequestID(subscriptionId string, resourceGroupName string, resourceGuardName string, deleteProtectedItemRequestName string) DeleteProtectedItemRequestId {
	return DeleteProtectedItemRequestId{
		SubscriptionId:                 subscriptionId,
		ResourceGroupName:              resourceGroupName,
		ResourceGuardName:              resourceGuardName,
		DeleteProtectedItemRequestName: deleteProtectedItemRequestName,
	}
}

// ParseDeleteProtectedItemRequestID parses 'input' into a DeleteProtectedItemRequestId
func ParseDeleteProtectedItemRequestID(input string) (*DeleteProtectedItemRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(DeleteProtectedItemRequestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DeleteProtectedItemRequestId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ResourceGuardName, ok = parsed.Parsed["resourceGuardName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGuardName", *parsed)
	}

	if id.DeleteProtectedItemRequestName, ok = parsed.Parsed["deleteProtectedItemRequestName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "deleteProtectedItemRequestName", *parsed)
	}

	return &id, nil
}

// ParseDeleteProtectedItemRequestIDInsensitively parses 'input' case-insensitively into a DeleteProtectedItemRequestId
// note: this method should only be used for API response data and not user input
func ParseDeleteProtectedItemRequestIDInsensitively(input string) (*DeleteProtectedItemRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(DeleteProtectedItemRequestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DeleteProtectedItemRequestId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ResourceGuardName, ok = parsed.Parsed["resourceGuardName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGuardName", *parsed)
	}

	if id.DeleteProtectedItemRequestName, ok = parsed.Parsed["deleteProtectedItemRequestName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "deleteProtectedItemRequestName", *parsed)
	}

	return &id, nil
}

// ValidateDeleteProtectedItemRequestID checks that 'input' can be parsed as a Delete Protected Item Request ID
func ValidateDeleteProtectedItemRequestID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeleteProtectedItemRequestID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Delete Protected Item Request ID
func (id DeleteProtectedItemRequestId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/resourceGuards/%s/deleteProtectedItemRequests/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceGuardName, id.DeleteProtectedItemRequestName)
}

// Segments returns a slice of Resource ID Segments which comprise this Delete Protected Item Request ID
func (id DeleteProtectedItemRequestId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticResourceGuards", "resourceGuards", "resourceGuards"),
		resourceids.UserSpecifiedSegment("resourceGuardName", "resourceGuardValue"),
		resourceids.StaticSegment("staticDeleteProtectedItemRequests", "deleteProtectedItemRequests", "deleteProtectedItemRequests"),
		resourceids.UserSpecifiedSegment("deleteProtectedItemRequestName", "deleteProtectedItemRequestValue"),
	}
}

// String returns a human-readable description of this Delete Protected Item Request ID
func (id DeleteProtectedItemRequestId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Guard Name: %q", id.ResourceGuardName),
		fmt.Sprintf("Delete Protected Item Request Name: %q", id.DeleteProtectedItemRequestName),
	}
	return fmt.Sprintf("Delete Protected Item Request (%s)", strings.Join(components, "\n"))
}
