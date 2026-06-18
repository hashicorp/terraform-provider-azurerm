package privatelinkhubs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PrivateLinkHubId{})
}

var _ resourceids.ResourceId = &PrivateLinkHubId{}

// PrivateLinkHubId is a struct representing the Resource ID for a Private Link Hub
type PrivateLinkHubId struct {
	SubscriptionId     string
	ResourceGroupName  string
	PrivateLinkHubName string
}

// NewPrivateLinkHubID returns a new PrivateLinkHubId struct
func NewPrivateLinkHubID(subscriptionId string, resourceGroupName string, privateLinkHubName string) PrivateLinkHubId {
	return PrivateLinkHubId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		PrivateLinkHubName: privateLinkHubName,
	}
}

// ParsePrivateLinkHubID parses 'input' into a PrivateLinkHubId
func ParsePrivateLinkHubID(input string) (*PrivateLinkHubId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateLinkHubId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateLinkHubId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePrivateLinkHubIDInsensitively parses 'input' case-insensitively into a PrivateLinkHubId
// note: this method should only be used for API response data and not user input
func ParsePrivateLinkHubIDInsensitively(input string) (*PrivateLinkHubId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateLinkHubId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateLinkHubId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PrivateLinkHubId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.PrivateLinkHubName, ok = input.Parsed["privateLinkHubName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateLinkHubName", input)
	}

	return nil
}

// ValidatePrivateLinkHubID checks that 'input' can be parsed as a Private Link Hub ID
func ValidatePrivateLinkHubID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateLinkHubID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Link Hub ID
func (id PrivateLinkHubId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Synapse/privateLinkHubs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateLinkHubName)
}

// Segments returns a slice of Resource ID Segments which comprise this Private Link Hub ID
func (id PrivateLinkHubId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSynapse", "Microsoft.Synapse", "Microsoft.Synapse"),
		resourceids.StaticSegment("staticPrivateLinkHubs", "privateLinkHubs", "privateLinkHubs"),
		resourceids.UserSpecifiedSegment("privateLinkHubName", "privateLinkHubName"),
	}
}

// String returns a human-readable description of this Private Link Hub ID
func (id PrivateLinkHubId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Link Hub Name: %q", id.PrivateLinkHubName),
	}
	return fmt.Sprintf("Private Link Hub (%s)", strings.Join(components, "\n"))
}
