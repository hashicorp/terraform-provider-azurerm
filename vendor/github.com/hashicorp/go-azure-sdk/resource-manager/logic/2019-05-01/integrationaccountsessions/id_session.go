package integrationaccountsessions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SessionId{}

// SessionId is a struct representing the Resource ID for a Session
type SessionId struct {
	SubscriptionId         string
	ResourceGroupName      string
	IntegrationAccountName string
	SessionName            string
}

// NewSessionID returns a new SessionId struct
func NewSessionID(subscriptionId string, resourceGroupName string, integrationAccountName string, sessionName string) SessionId {
	return SessionId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		IntegrationAccountName: integrationAccountName,
		SessionName:            sessionName,
	}
}

// ParseSessionID parses 'input' into a SessionId
func ParseSessionID(input string) (*SessionId, error) {
	parser := resourceids.NewParserFromResourceIdType(SessionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SessionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.IntegrationAccountName, ok = parsed.Parsed["integrationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "integrationAccountName", *parsed)
	}

	if id.SessionName, ok = parsed.Parsed["sessionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sessionName", *parsed)
	}

	return &id, nil
}

// ParseSessionIDInsensitively parses 'input' case-insensitively into a SessionId
// note: this method should only be used for API response data and not user input
func ParseSessionIDInsensitively(input string) (*SessionId, error) {
	parser := resourceids.NewParserFromResourceIdType(SessionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SessionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.IntegrationAccountName, ok = parsed.Parsed["integrationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "integrationAccountName", *parsed)
	}

	if id.SessionName, ok = parsed.Parsed["sessionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sessionName", *parsed)
	}

	return &id, nil
}

// ValidateSessionID checks that 'input' can be parsed as a Session ID
func ValidateSessionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSessionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Session ID
func (id SessionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/sessions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.IntegrationAccountName, id.SessionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Session ID
func (id SessionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLogic", "Microsoft.Logic", "Microsoft.Logic"),
		resourceids.StaticSegment("staticIntegrationAccounts", "integrationAccounts", "integrationAccounts"),
		resourceids.UserSpecifiedSegment("integrationAccountName", "integrationAccountValue"),
		resourceids.StaticSegment("staticSessions", "sessions", "sessions"),
		resourceids.UserSpecifiedSegment("sessionName", "sessionValue"),
	}
}

// String returns a human-readable description of this Session ID
func (id SessionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Integration Account Name: %q", id.IntegrationAccountName),
		fmt.Sprintf("Session Name: %q", id.SessionName),
	}
	return fmt.Sprintf("Session (%s)", strings.Join(components, "\n"))
}
