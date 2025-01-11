package resourceguards

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&UpdateProtectedItemRequestId{})
}

var _ resourceids.ResourceId = &UpdateProtectedItemRequestId{}

// UpdateProtectedItemRequestId is a struct representing the Resource ID for a Update Protected Item Request
type UpdateProtectedItemRequestId struct {
	SubscriptionId                 string
	ResourceGroupName              string
	ResourceGuardName              string
	UpdateProtectedItemRequestName string
}

// NewUpdateProtectedItemRequestID returns a new UpdateProtectedItemRequestId struct
func NewUpdateProtectedItemRequestID(subscriptionId string, resourceGroupName string, resourceGuardName string, updateProtectedItemRequestName string) UpdateProtectedItemRequestId {
	return UpdateProtectedItemRequestId{
		SubscriptionId:                 subscriptionId,
		ResourceGroupName:              resourceGroupName,
		ResourceGuardName:              resourceGuardName,
		UpdateProtectedItemRequestName: updateProtectedItemRequestName,
	}
}

// ParseUpdateProtectedItemRequestID parses 'input' into a UpdateProtectedItemRequestId
func ParseUpdateProtectedItemRequestID(input string) (*UpdateProtectedItemRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(&UpdateProtectedItemRequestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := UpdateProtectedItemRequestId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseUpdateProtectedItemRequestIDInsensitively parses 'input' case-insensitively into a UpdateProtectedItemRequestId
// note: this method should only be used for API response data and not user input
func ParseUpdateProtectedItemRequestIDInsensitively(input string) (*UpdateProtectedItemRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(&UpdateProtectedItemRequestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := UpdateProtectedItemRequestId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *UpdateProtectedItemRequestId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ResourceGuardName, ok = input.Parsed["resourceGuardName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGuardName", input)
	}

	if id.UpdateProtectedItemRequestName, ok = input.Parsed["updateProtectedItemRequestName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "updateProtectedItemRequestName", input)
	}

	return nil
}

// ValidateUpdateProtectedItemRequestID checks that 'input' can be parsed as a Update Protected Item Request ID
func ValidateUpdateProtectedItemRequestID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseUpdateProtectedItemRequestID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Update Protected Item Request ID
func (id UpdateProtectedItemRequestId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/resourceGuards/%s/updateProtectedItemRequests/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceGuardName, id.UpdateProtectedItemRequestName)
}

// Segments returns a slice of Resource ID Segments which comprise this Update Protected Item Request ID
func (id UpdateProtectedItemRequestId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticResourceGuards", "resourceGuards", "resourceGuards"),
		resourceids.UserSpecifiedSegment("resourceGuardName", "resourceGuardName"),
		resourceids.StaticSegment("staticUpdateProtectedItemRequests", "updateProtectedItemRequests", "updateProtectedItemRequests"),
		resourceids.UserSpecifiedSegment("updateProtectedItemRequestName", "updateProtectedItemRequestName"),
	}
}

// String returns a human-readable description of this Update Protected Item Request ID
func (id UpdateProtectedItemRequestId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Guard Name: %q", id.ResourceGuardName),
		fmt.Sprintf("Update Protected Item Request Name: %q", id.UpdateProtectedItemRequestName),
	}
	return fmt.Sprintf("Update Protected Item Request (%s)", strings.Join(components, "\n"))
}
