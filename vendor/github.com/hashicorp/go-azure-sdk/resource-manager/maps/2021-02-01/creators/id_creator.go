package creators

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CreatorId{}

// CreatorId is a struct representing the Resource ID for a Creator
type CreatorId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	CreatorName       string
}

// NewCreatorID returns a new CreatorId struct
func NewCreatorID(subscriptionId string, resourceGroupName string, accountName string, creatorName string) CreatorId {
	return CreatorId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		CreatorName:       creatorName,
	}
}

// ParseCreatorID parses 'input' into a CreatorId
func ParseCreatorID(input string) (*CreatorId, error) {
	parser := resourceids.NewParserFromResourceIdType(CreatorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CreatorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "accountName", *parsed)
	}

	if id.CreatorName, ok = parsed.Parsed["creatorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "creatorName", *parsed)
	}

	return &id, nil
}

// ParseCreatorIDInsensitively parses 'input' case-insensitively into a CreatorId
// note: this method should only be used for API response data and not user input
func ParseCreatorIDInsensitively(input string) (*CreatorId, error) {
	parser := resourceids.NewParserFromResourceIdType(CreatorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CreatorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "accountName", *parsed)
	}

	if id.CreatorName, ok = parsed.Parsed["creatorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "creatorName", *parsed)
	}

	return &id, nil
}

// ValidateCreatorID checks that 'input' can be parsed as a Creator ID
func ValidateCreatorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCreatorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Creator ID
func (id CreatorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Maps/accounts/%s/creators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.CreatorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Creator ID
func (id CreatorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMaps", "Microsoft.Maps", "Microsoft.Maps"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticCreators", "creators", "creators"),
		resourceids.UserSpecifiedSegment("creatorName", "creatorValue"),
	}
}

// String returns a human-readable description of this Creator ID
func (id CreatorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Creator Name: %q", id.CreatorName),
	}
	return fmt.Sprintf("Creator (%s)", strings.Join(components, "\n"))
}
