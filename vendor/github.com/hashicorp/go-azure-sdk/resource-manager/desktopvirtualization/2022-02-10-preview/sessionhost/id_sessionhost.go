package sessionhost

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SessionHostId{}

// SessionHostId is a struct representing the Resource ID for a Session Host
type SessionHostId struct {
	SubscriptionId    string
	ResourceGroupName string
	HostPoolName      string
	SessionHostName   string
}

// NewSessionHostID returns a new SessionHostId struct
func NewSessionHostID(subscriptionId string, resourceGroupName string, hostPoolName string, sessionHostName string) SessionHostId {
	return SessionHostId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		HostPoolName:      hostPoolName,
		SessionHostName:   sessionHostName,
	}
}

// ParseSessionHostID parses 'input' into a SessionHostId
func ParseSessionHostID(input string) (*SessionHostId, error) {
	parser := resourceids.NewParserFromResourceIdType(SessionHostId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SessionHostId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.HostPoolName, ok = parsed.Parsed["hostPoolName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "hostPoolName", *parsed)
	}

	if id.SessionHostName, ok = parsed.Parsed["sessionHostName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sessionHostName", *parsed)
	}

	return &id, nil
}

// ParseSessionHostIDInsensitively parses 'input' case-insensitively into a SessionHostId
// note: this method should only be used for API response data and not user input
func ParseSessionHostIDInsensitively(input string) (*SessionHostId, error) {
	parser := resourceids.NewParserFromResourceIdType(SessionHostId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SessionHostId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.HostPoolName, ok = parsed.Parsed["hostPoolName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "hostPoolName", *parsed)
	}

	if id.SessionHostName, ok = parsed.Parsed["sessionHostName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sessionHostName", *parsed)
	}

	return &id, nil
}

// ValidateSessionHostID checks that 'input' can be parsed as a Session Host ID
func ValidateSessionHostID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSessionHostID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Session Host ID
func (id SessionHostId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/hostPools/%s/sessionHosts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HostPoolName, id.SessionHostName)
}

// Segments returns a slice of Resource ID Segments which comprise this Session Host ID
func (id SessionHostId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDesktopVirtualization", "Microsoft.DesktopVirtualization", "Microsoft.DesktopVirtualization"),
		resourceids.StaticSegment("staticHostPools", "hostPools", "hostPools"),
		resourceids.UserSpecifiedSegment("hostPoolName", "hostPoolValue"),
		resourceids.StaticSegment("staticSessionHosts", "sessionHosts", "sessionHosts"),
		resourceids.UserSpecifiedSegment("sessionHostName", "sessionHostValue"),
	}
}

// String returns a human-readable description of this Session Host ID
func (id SessionHostId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Host Pool Name: %q", id.HostPoolName),
		fmt.Sprintf("Session Host Name: %q", id.SessionHostName),
	}
	return fmt.Sprintf("Session Host (%s)", strings.Join(components, "\n"))
}
