package synchronizationsetting

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SynchronizationSettingId{})
}

var _ resourceids.ResourceId = &SynchronizationSettingId{}

// SynchronizationSettingId is a struct representing the Resource ID for a Synchronization Setting
type SynchronizationSettingId struct {
	SubscriptionId             string
	ResourceGroupName          string
	AccountName                string
	ShareName                  string
	SynchronizationSettingName string
}

// NewSynchronizationSettingID returns a new SynchronizationSettingId struct
func NewSynchronizationSettingID(subscriptionId string, resourceGroupName string, accountName string, shareName string, synchronizationSettingName string) SynchronizationSettingId {
	return SynchronizationSettingId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		AccountName:                accountName,
		ShareName:                  shareName,
		SynchronizationSettingName: synchronizationSettingName,
	}
}

// ParseSynchronizationSettingID parses 'input' into a SynchronizationSettingId
func ParseSynchronizationSettingID(input string) (*SynchronizationSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SynchronizationSettingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SynchronizationSettingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSynchronizationSettingIDInsensitively parses 'input' case-insensitively into a SynchronizationSettingId
// note: this method should only be used for API response data and not user input
func ParseSynchronizationSettingIDInsensitively(input string) (*SynchronizationSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SynchronizationSettingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SynchronizationSettingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SynchronizationSettingId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AccountName, ok = input.Parsed["accountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "accountName", input)
	}

	if id.ShareName, ok = input.Parsed["shareName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "shareName", input)
	}

	if id.SynchronizationSettingName, ok = input.Parsed["synchronizationSettingName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "synchronizationSettingName", input)
	}

	return nil
}

// ValidateSynchronizationSettingID checks that 'input' can be parsed as a Synchronization Setting ID
func ValidateSynchronizationSettingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSynchronizationSettingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Synchronization Setting ID
func (id SynchronizationSettingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataShare/accounts/%s/shares/%s/synchronizationSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ShareName, id.SynchronizationSettingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Synchronization Setting ID
func (id SynchronizationSettingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataShare", "Microsoft.DataShare", "Microsoft.DataShare"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountName"),
		resourceids.StaticSegment("staticShares", "shares", "shares"),
		resourceids.UserSpecifiedSegment("shareName", "shareName"),
		resourceids.StaticSegment("staticSynchronizationSettings", "synchronizationSettings", "synchronizationSettings"),
		resourceids.UserSpecifiedSegment("synchronizationSettingName", "synchronizationSettingName"),
	}
}

// String returns a human-readable description of this Synchronization Setting ID
func (id SynchronizationSettingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Share Name: %q", id.ShareName),
		fmt.Sprintf("Synchronization Setting Name: %q", id.SynchronizationSettingName),
	}
	return fmt.Sprintf("Synchronization Setting (%s)", strings.Join(components, "\n"))
}
