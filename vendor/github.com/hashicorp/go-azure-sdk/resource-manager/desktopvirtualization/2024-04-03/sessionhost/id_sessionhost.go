package sessionhost

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SessionHostId{})
}

var _ resourceids.ResourceId = &SessionHostId{}

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
	parser := resourceids.NewParserFromResourceIdType(&SessionHostId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SessionHostId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSessionHostIDInsensitively parses 'input' case-insensitively into a SessionHostId
// note: this method should only be used for API response data and not user input
func ParseSessionHostIDInsensitively(input string) (*SessionHostId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SessionHostId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SessionHostId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SessionHostId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.HostPoolName, ok = input.Parsed["hostPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hostPoolName", input)
	}

	if id.SessionHostName, ok = input.Parsed["sessionHostName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sessionHostName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("hostPoolName", "hostPoolName"),
		resourceids.StaticSegment("staticSessionHosts", "sessionHosts", "sessionHosts"),
		resourceids.UserSpecifiedSegment("sessionHostName", "sessionHostName"),
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
