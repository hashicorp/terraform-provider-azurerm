package managedhsms

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DeletedManagedHSMId{})
}

var _ resourceids.ResourceId = &DeletedManagedHSMId{}

// DeletedManagedHSMId is a struct representing the Resource ID for a Deleted Managed H S M
type DeletedManagedHSMId struct {
	SubscriptionId        string
	LocationName          string
	DeletedManagedHSMName string
}

// NewDeletedManagedHSMID returns a new DeletedManagedHSMId struct
func NewDeletedManagedHSMID(subscriptionId string, locationName string, deletedManagedHSMName string) DeletedManagedHSMId {
	return DeletedManagedHSMId{
		SubscriptionId:        subscriptionId,
		LocationName:          locationName,
		DeletedManagedHSMName: deletedManagedHSMName,
	}
}

// ParseDeletedManagedHSMID parses 'input' into a DeletedManagedHSMId
func ParseDeletedManagedHSMID(input string) (*DeletedManagedHSMId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedManagedHSMId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedManagedHSMId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDeletedManagedHSMIDInsensitively parses 'input' case-insensitively into a DeletedManagedHSMId
// note: this method should only be used for API response data and not user input
func ParseDeletedManagedHSMIDInsensitively(input string) (*DeletedManagedHSMId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedManagedHSMId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedManagedHSMId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DeletedManagedHSMId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.DeletedManagedHSMName, ok = input.Parsed["deletedManagedHSMName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deletedManagedHSMName", input)
	}

	return nil
}

// ValidateDeletedManagedHSMID checks that 'input' can be parsed as a Deleted Managed H S M ID
func ValidateDeletedManagedHSMID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeletedManagedHSMID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deleted Managed H S M ID
func (id DeletedManagedHSMId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.KeyVault/locations/%s/deletedManagedHSMs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.DeletedManagedHSMName)
}

// Segments returns a slice of Resource ID Segments which comprise this Deleted Managed H S M ID
func (id DeletedManagedHSMId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKeyVault", "Microsoft.KeyVault", "Microsoft.KeyVault"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticDeletedManagedHSMs", "deletedManagedHSMs", "deletedManagedHSMs"),
		resourceids.UserSpecifiedSegment("deletedManagedHSMName", "deletedManagedHSMName"),
	}
}

// String returns a human-readable description of this Deleted Managed H S M ID
func (id DeletedManagedHSMId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Deleted Managed H S M Name: %q", id.DeletedManagedHSMName),
	}
	return fmt.Sprintf("Deleted Managed H S M (%s)", strings.Join(components, "\n"))
}
