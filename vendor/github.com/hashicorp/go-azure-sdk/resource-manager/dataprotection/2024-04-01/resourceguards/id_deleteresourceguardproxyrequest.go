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
	recaser.RegisterResourceId(&DeleteResourceGuardProxyRequestId{})
}

var _ resourceids.ResourceId = &DeleteResourceGuardProxyRequestId{}

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
	parser := resourceids.NewParserFromResourceIdType(&DeleteResourceGuardProxyRequestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeleteResourceGuardProxyRequestId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDeleteResourceGuardProxyRequestIDInsensitively parses 'input' case-insensitively into a DeleteResourceGuardProxyRequestId
// note: this method should only be used for API response data and not user input
func ParseDeleteResourceGuardProxyRequestIDInsensitively(input string) (*DeleteResourceGuardProxyRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeleteResourceGuardProxyRequestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeleteResourceGuardProxyRequestId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DeleteResourceGuardProxyRequestId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DeleteResourceGuardProxyRequestName, ok = input.Parsed["deleteResourceGuardProxyRequestName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deleteResourceGuardProxyRequestName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("resourceGuardName", "resourceGuardName"),
		resourceids.StaticSegment("staticDeleteResourceGuardProxyRequests", "deleteResourceGuardProxyRequests", "deleteResourceGuardProxyRequests"),
		resourceids.UserSpecifiedSegment("deleteResourceGuardProxyRequestName", "deleteResourceGuardProxyRequestName"),
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
