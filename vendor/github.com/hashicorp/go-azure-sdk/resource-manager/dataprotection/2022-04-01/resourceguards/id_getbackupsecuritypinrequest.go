package resourceguards

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = GetBackupSecurityPINRequestId{}

// GetBackupSecurityPINRequestId is a struct representing the Resource ID for a Get Backup Security P I N Request
type GetBackupSecurityPINRequestId struct {
	SubscriptionId                  string
	ResourceGroupName               string
	ResourceGuardName               string
	GetBackupSecurityPINRequestName string
}

// NewGetBackupSecurityPINRequestID returns a new GetBackupSecurityPINRequestId struct
func NewGetBackupSecurityPINRequestID(subscriptionId string, resourceGroupName string, resourceGuardName string, getBackupSecurityPINRequestName string) GetBackupSecurityPINRequestId {
	return GetBackupSecurityPINRequestId{
		SubscriptionId:                  subscriptionId,
		ResourceGroupName:               resourceGroupName,
		ResourceGuardName:               resourceGuardName,
		GetBackupSecurityPINRequestName: getBackupSecurityPINRequestName,
	}
}

// ParseGetBackupSecurityPINRequestID parses 'input' into a GetBackupSecurityPINRequestId
func ParseGetBackupSecurityPINRequestID(input string) (*GetBackupSecurityPINRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(GetBackupSecurityPINRequestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GetBackupSecurityPINRequestId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ResourceGuardName, ok = parsed.Parsed["resourceGuardName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGuardName", *parsed)
	}

	if id.GetBackupSecurityPINRequestName, ok = parsed.Parsed["getBackupSecurityPINRequestName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "getBackupSecurityPINRequestName", *parsed)
	}

	return &id, nil
}

// ParseGetBackupSecurityPINRequestIDInsensitively parses 'input' case-insensitively into a GetBackupSecurityPINRequestId
// note: this method should only be used for API response data and not user input
func ParseGetBackupSecurityPINRequestIDInsensitively(input string) (*GetBackupSecurityPINRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(GetBackupSecurityPINRequestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GetBackupSecurityPINRequestId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ResourceGuardName, ok = parsed.Parsed["resourceGuardName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGuardName", *parsed)
	}

	if id.GetBackupSecurityPINRequestName, ok = parsed.Parsed["getBackupSecurityPINRequestName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "getBackupSecurityPINRequestName", *parsed)
	}

	return &id, nil
}

// ValidateGetBackupSecurityPINRequestID checks that 'input' can be parsed as a Get Backup Security P I N Request ID
func ValidateGetBackupSecurityPINRequestID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGetBackupSecurityPINRequestID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Get Backup Security P I N Request ID
func (id GetBackupSecurityPINRequestId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/resourceGuards/%s/getBackupSecurityPINRequests/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceGuardName, id.GetBackupSecurityPINRequestName)
}

// Segments returns a slice of Resource ID Segments which comprise this Get Backup Security P I N Request ID
func (id GetBackupSecurityPINRequestId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticResourceGuards", "resourceGuards", "resourceGuards"),
		resourceids.UserSpecifiedSegment("resourceGuardName", "resourceGuardValue"),
		resourceids.StaticSegment("staticGetBackupSecurityPINRequests", "getBackupSecurityPINRequests", "getBackupSecurityPINRequests"),
		resourceids.UserSpecifiedSegment("getBackupSecurityPINRequestName", "getBackupSecurityPINRequestValue"),
	}
}

// String returns a human-readable description of this Get Backup Security P I N Request ID
func (id GetBackupSecurityPINRequestId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Guard Name: %q", id.ResourceGuardName),
		fmt.Sprintf("Get Backup Security P I N Request Name: %q", id.GetBackupSecurityPINRequestName),
	}
	return fmt.Sprintf("Get Backup Security P I N Request (%s)", strings.Join(components, "\n"))
}
