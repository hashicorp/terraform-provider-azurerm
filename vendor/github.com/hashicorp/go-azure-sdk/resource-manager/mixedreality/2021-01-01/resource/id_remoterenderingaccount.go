package resource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RemoteRenderingAccountId{})
}

var _ resourceids.ResourceId = &RemoteRenderingAccountId{}

// RemoteRenderingAccountId is a struct representing the Resource ID for a Remote Rendering Account
type RemoteRenderingAccountId struct {
	SubscriptionId             string
	ResourceGroupName          string
	RemoteRenderingAccountName string
}

// NewRemoteRenderingAccountID returns a new RemoteRenderingAccountId struct
func NewRemoteRenderingAccountID(subscriptionId string, resourceGroupName string, remoteRenderingAccountName string) RemoteRenderingAccountId {
	return RemoteRenderingAccountId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		RemoteRenderingAccountName: remoteRenderingAccountName,
	}
}

// ParseRemoteRenderingAccountID parses 'input' into a RemoteRenderingAccountId
func ParseRemoteRenderingAccountID(input string) (*RemoteRenderingAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RemoteRenderingAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RemoteRenderingAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRemoteRenderingAccountIDInsensitively parses 'input' case-insensitively into a RemoteRenderingAccountId
// note: this method should only be used for API response data and not user input
func ParseRemoteRenderingAccountIDInsensitively(input string) (*RemoteRenderingAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RemoteRenderingAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RemoteRenderingAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RemoteRenderingAccountId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RemoteRenderingAccountName, ok = input.Parsed["remoteRenderingAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "remoteRenderingAccountName", input)
	}

	return nil
}

// ValidateRemoteRenderingAccountID checks that 'input' can be parsed as a Remote Rendering Account ID
func ValidateRemoteRenderingAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRemoteRenderingAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Remote Rendering Account ID
func (id RemoteRenderingAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MixedReality/remoteRenderingAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RemoteRenderingAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Remote Rendering Account ID
func (id RemoteRenderingAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMixedReality", "Microsoft.MixedReality", "Microsoft.MixedReality"),
		resourceids.StaticSegment("staticRemoteRenderingAccounts", "remoteRenderingAccounts", "remoteRenderingAccounts"),
		resourceids.UserSpecifiedSegment("remoteRenderingAccountName", "remoteRenderingAccountName"),
	}
}

// String returns a human-readable description of this Remote Rendering Account ID
func (id RemoteRenderingAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Remote Rendering Account Name: %q", id.RemoteRenderingAccountName),
	}
	return fmt.Sprintf("Remote Rendering Account (%s)", strings.Join(components, "\n"))
}
