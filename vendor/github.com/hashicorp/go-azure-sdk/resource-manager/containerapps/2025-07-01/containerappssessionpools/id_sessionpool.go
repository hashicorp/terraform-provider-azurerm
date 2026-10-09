package containerappssessionpools

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SessionPoolId{})
}

var _ resourceids.ResourceId = &SessionPoolId{}

// SessionPoolId is a struct representing the Resource ID for a Session Pool
type SessionPoolId struct {
	SubscriptionId    string
	ResourceGroupName string
	SessionPoolName   string
}

// NewSessionPoolID returns a new SessionPoolId struct
func NewSessionPoolID(subscriptionId string, resourceGroupName string, sessionPoolName string) SessionPoolId {
	return SessionPoolId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SessionPoolName:   sessionPoolName,
	}
}

// ParseSessionPoolID parses 'input' into a SessionPoolId
func ParseSessionPoolID(input string) (*SessionPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SessionPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SessionPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSessionPoolIDInsensitively parses 'input' case-insensitively into a SessionPoolId
// note: this method should only be used for API response data and not user input
func ParseSessionPoolIDInsensitively(input string) (*SessionPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SessionPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SessionPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SessionPoolId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SessionPoolName, ok = input.Parsed["sessionPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sessionPoolName", input)
	}

	return nil
}

// ValidateSessionPoolID checks that 'input' can be parsed as a Session Pool ID
func ValidateSessionPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSessionPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Session Pool ID
func (id SessionPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/sessionPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SessionPoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Session Pool ID
func (id SessionPoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticSessionPools", "sessionPools", "sessionPools"),
		resourceids.UserSpecifiedSegment("sessionPoolName", "sessionPoolName"),
	}
}

// String returns a human-readable description of this Session Pool ID
func (id SessionPoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Session Pool Name: %q", id.SessionPoolName),
	}
	return fmt.Sprintf("Session Pool (%s)", strings.Join(components, "\n"))
}
