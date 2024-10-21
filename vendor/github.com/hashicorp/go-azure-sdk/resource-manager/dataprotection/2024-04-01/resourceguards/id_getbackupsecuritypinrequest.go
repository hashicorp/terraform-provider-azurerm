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
	recaser.RegisterResourceId(&GetBackupSecurityPINRequestId{})
}

var _ resourceids.ResourceId = &GetBackupSecurityPINRequestId{}

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
	parser := resourceids.NewParserFromResourceIdType(&GetBackupSecurityPINRequestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GetBackupSecurityPINRequestId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGetBackupSecurityPINRequestIDInsensitively parses 'input' case-insensitively into a GetBackupSecurityPINRequestId
// note: this method should only be used for API response data and not user input
func ParseGetBackupSecurityPINRequestIDInsensitively(input string) (*GetBackupSecurityPINRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GetBackupSecurityPINRequestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GetBackupSecurityPINRequestId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GetBackupSecurityPINRequestId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.GetBackupSecurityPINRequestName, ok = input.Parsed["getBackupSecurityPINRequestName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "getBackupSecurityPINRequestName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("resourceGuardName", "resourceGuardName"),
		resourceids.StaticSegment("staticGetBackupSecurityPINRequests", "getBackupSecurityPINRequests", "getBackupSecurityPINRequests"),
		resourceids.UserSpecifiedSegment("getBackupSecurityPINRequestName", "getBackupSecurityPINRequestName"),
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
