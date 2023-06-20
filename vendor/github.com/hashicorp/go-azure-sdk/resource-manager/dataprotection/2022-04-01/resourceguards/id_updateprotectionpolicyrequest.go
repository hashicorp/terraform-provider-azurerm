package resourceguards

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = UpdateProtectionPolicyRequestId{}

// UpdateProtectionPolicyRequestId is a struct representing the Resource ID for a Update Protection Policy Request
type UpdateProtectionPolicyRequestId struct {
	SubscriptionId                    string
	ResourceGroupName                 string
	ResourceGuardName                 string
	UpdateProtectionPolicyRequestName string
}

// NewUpdateProtectionPolicyRequestID returns a new UpdateProtectionPolicyRequestId struct
func NewUpdateProtectionPolicyRequestID(subscriptionId string, resourceGroupName string, resourceGuardName string, updateProtectionPolicyRequestName string) UpdateProtectionPolicyRequestId {
	return UpdateProtectionPolicyRequestId{
		SubscriptionId:                    subscriptionId,
		ResourceGroupName:                 resourceGroupName,
		ResourceGuardName:                 resourceGuardName,
		UpdateProtectionPolicyRequestName: updateProtectionPolicyRequestName,
	}
}

// ParseUpdateProtectionPolicyRequestID parses 'input' into a UpdateProtectionPolicyRequestId
func ParseUpdateProtectionPolicyRequestID(input string) (*UpdateProtectionPolicyRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(UpdateProtectionPolicyRequestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UpdateProtectionPolicyRequestId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ResourceGuardName, ok = parsed.Parsed["resourceGuardName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGuardName", *parsed)
	}

	if id.UpdateProtectionPolicyRequestName, ok = parsed.Parsed["updateProtectionPolicyRequestName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "updateProtectionPolicyRequestName", *parsed)
	}

	return &id, nil
}

// ParseUpdateProtectionPolicyRequestIDInsensitively parses 'input' case-insensitively into a UpdateProtectionPolicyRequestId
// note: this method should only be used for API response data and not user input
func ParseUpdateProtectionPolicyRequestIDInsensitively(input string) (*UpdateProtectionPolicyRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(UpdateProtectionPolicyRequestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UpdateProtectionPolicyRequestId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ResourceGuardName, ok = parsed.Parsed["resourceGuardName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGuardName", *parsed)
	}

	if id.UpdateProtectionPolicyRequestName, ok = parsed.Parsed["updateProtectionPolicyRequestName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "updateProtectionPolicyRequestName", *parsed)
	}

	return &id, nil
}

// ValidateUpdateProtectionPolicyRequestID checks that 'input' can be parsed as a Update Protection Policy Request ID
func ValidateUpdateProtectionPolicyRequestID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseUpdateProtectionPolicyRequestID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Update Protection Policy Request ID
func (id UpdateProtectionPolicyRequestId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/resourceGuards/%s/updateProtectionPolicyRequests/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceGuardName, id.UpdateProtectionPolicyRequestName)
}

// Segments returns a slice of Resource ID Segments which comprise this Update Protection Policy Request ID
func (id UpdateProtectionPolicyRequestId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticResourceGuards", "resourceGuards", "resourceGuards"),
		resourceids.UserSpecifiedSegment("resourceGuardName", "resourceGuardValue"),
		resourceids.StaticSegment("staticUpdateProtectionPolicyRequests", "updateProtectionPolicyRequests", "updateProtectionPolicyRequests"),
		resourceids.UserSpecifiedSegment("updateProtectionPolicyRequestName", "updateProtectionPolicyRequestValue"),
	}
}

// String returns a human-readable description of this Update Protection Policy Request ID
func (id UpdateProtectionPolicyRequestId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Guard Name: %q", id.ResourceGuardName),
		fmt.Sprintf("Update Protection Policy Request Name: %q", id.UpdateProtectionPolicyRequestName),
	}
	return fmt.Sprintf("Update Protection Policy Request (%s)", strings.Join(components, "\n"))
}
