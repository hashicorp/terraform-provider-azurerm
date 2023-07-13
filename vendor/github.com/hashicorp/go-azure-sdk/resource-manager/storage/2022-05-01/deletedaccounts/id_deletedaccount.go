package deletedaccounts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DeletedAccountId{}

// DeletedAccountId is a struct representing the Resource ID for a Deleted Account
type DeletedAccountId struct {
	SubscriptionId     string
	LocationName       string
	DeletedAccountName string
}

// NewDeletedAccountID returns a new DeletedAccountId struct
func NewDeletedAccountID(subscriptionId string, locationName string, deletedAccountName string) DeletedAccountId {
	return DeletedAccountId{
		SubscriptionId:     subscriptionId,
		LocationName:       locationName,
		DeletedAccountName: deletedAccountName,
	}
}

// ParseDeletedAccountID parses 'input' into a DeletedAccountId
func ParseDeletedAccountID(input string) (*DeletedAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(DeletedAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DeletedAccountId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.DeletedAccountName, ok = parsed.Parsed["deletedAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "deletedAccountName", *parsed)
	}

	return &id, nil
}

// ParseDeletedAccountIDInsensitively parses 'input' case-insensitively into a DeletedAccountId
// note: this method should only be used for API response data and not user input
func ParseDeletedAccountIDInsensitively(input string) (*DeletedAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(DeletedAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DeletedAccountId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.DeletedAccountName, ok = parsed.Parsed["deletedAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "deletedAccountName", *parsed)
	}

	return &id, nil
}

// ValidateDeletedAccountID checks that 'input' can be parsed as a Deleted Account ID
func ValidateDeletedAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeletedAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deleted Account ID
func (id DeletedAccountId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Storage/locations/%s/deletedAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.DeletedAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Deleted Account ID
func (id DeletedAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorage", "Microsoft.Storage", "Microsoft.Storage"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationValue"),
		resourceids.StaticSegment("staticDeletedAccounts", "deletedAccounts", "deletedAccounts"),
		resourceids.UserSpecifiedSegment("deletedAccountName", "deletedAccountValue"),
	}
}

// String returns a human-readable description of this Deleted Account ID
func (id DeletedAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Deleted Account Name: %q", id.DeletedAccountName),
	}
	return fmt.Sprintf("Deleted Account (%s)", strings.Join(components, "\n"))
}
