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
	recaser.RegisterResourceId(&DisableSoftDeleteRequestId{})
}

var _ resourceids.ResourceId = &DisableSoftDeleteRequestId{}

// DisableSoftDeleteRequestId is a struct representing the Resource ID for a Disable Soft Delete Request
type DisableSoftDeleteRequestId struct {
	SubscriptionId               string
	ResourceGroupName            string
	ResourceGuardName            string
	DisableSoftDeleteRequestName string
}

// NewDisableSoftDeleteRequestID returns a new DisableSoftDeleteRequestId struct
func NewDisableSoftDeleteRequestID(subscriptionId string, resourceGroupName string, resourceGuardName string, disableSoftDeleteRequestName string) DisableSoftDeleteRequestId {
	return DisableSoftDeleteRequestId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		ResourceGuardName:            resourceGuardName,
		DisableSoftDeleteRequestName: disableSoftDeleteRequestName,
	}
}

// ParseDisableSoftDeleteRequestID parses 'input' into a DisableSoftDeleteRequestId
func ParseDisableSoftDeleteRequestID(input string) (*DisableSoftDeleteRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DisableSoftDeleteRequestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DisableSoftDeleteRequestId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDisableSoftDeleteRequestIDInsensitively parses 'input' case-insensitively into a DisableSoftDeleteRequestId
// note: this method should only be used for API response data and not user input
func ParseDisableSoftDeleteRequestIDInsensitively(input string) (*DisableSoftDeleteRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DisableSoftDeleteRequestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DisableSoftDeleteRequestId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DisableSoftDeleteRequestId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DisableSoftDeleteRequestName, ok = input.Parsed["disableSoftDeleteRequestName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "disableSoftDeleteRequestName", input)
	}

	return nil
}

// ValidateDisableSoftDeleteRequestID checks that 'input' can be parsed as a Disable Soft Delete Request ID
func ValidateDisableSoftDeleteRequestID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDisableSoftDeleteRequestID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Disable Soft Delete Request ID
func (id DisableSoftDeleteRequestId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/resourceGuards/%s/disableSoftDeleteRequests/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceGuardName, id.DisableSoftDeleteRequestName)
}

// Segments returns a slice of Resource ID Segments which comprise this Disable Soft Delete Request ID
func (id DisableSoftDeleteRequestId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticResourceGuards", "resourceGuards", "resourceGuards"),
		resourceids.UserSpecifiedSegment("resourceGuardName", "resourceGuardName"),
		resourceids.StaticSegment("staticDisableSoftDeleteRequests", "disableSoftDeleteRequests", "disableSoftDeleteRequests"),
		resourceids.UserSpecifiedSegment("disableSoftDeleteRequestName", "disableSoftDeleteRequestName"),
	}
}

// String returns a human-readable description of this Disable Soft Delete Request ID
func (id DisableSoftDeleteRequestId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Guard Name: %q", id.ResourceGuardName),
		fmt.Sprintf("Disable Soft Delete Request Name: %q", id.DisableSoftDeleteRequestName),
	}
	return fmt.Sprintf("Disable Soft Delete Request (%s)", strings.Join(components, "\n"))
}
