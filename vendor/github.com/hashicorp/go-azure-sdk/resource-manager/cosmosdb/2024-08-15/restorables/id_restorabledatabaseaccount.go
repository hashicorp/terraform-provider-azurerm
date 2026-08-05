package restorables

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RestorableDatabaseAccountId{})
}

var _ resourceids.ResourceId = &RestorableDatabaseAccountId{}

// RestorableDatabaseAccountId is a struct representing the Resource ID for a Restorable Database Account
type RestorableDatabaseAccountId struct {
	SubscriptionId string
	LocationName   string
	InstanceId     string
}

// NewRestorableDatabaseAccountID returns a new RestorableDatabaseAccountId struct
func NewRestorableDatabaseAccountID(subscriptionId string, locationName string, instanceId string) RestorableDatabaseAccountId {
	return RestorableDatabaseAccountId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		InstanceId:     instanceId,
	}
}

// ParseRestorableDatabaseAccountID parses 'input' into a RestorableDatabaseAccountId
func ParseRestorableDatabaseAccountID(input string) (*RestorableDatabaseAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RestorableDatabaseAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RestorableDatabaseAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRestorableDatabaseAccountIDInsensitively parses 'input' case-insensitively into a RestorableDatabaseAccountId
// note: this method should only be used for API response data and not user input
func ParseRestorableDatabaseAccountIDInsensitively(input string) (*RestorableDatabaseAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RestorableDatabaseAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RestorableDatabaseAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RestorableDatabaseAccountId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.InstanceId, ok = input.Parsed["instanceId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "instanceId", input)
	}

	return nil
}

// ValidateRestorableDatabaseAccountID checks that 'input' can be parsed as a Restorable Database Account ID
func ValidateRestorableDatabaseAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRestorableDatabaseAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Restorable Database Account ID
func (id RestorableDatabaseAccountId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.DocumentDB/locations/%s/restorableDatabaseAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.InstanceId)
}

// Segments returns a slice of Resource ID Segments which comprise this Restorable Database Account ID
func (id RestorableDatabaseAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticRestorableDatabaseAccounts", "restorableDatabaseAccounts", "restorableDatabaseAccounts"),
		resourceids.UserSpecifiedSegment("instanceId", "instanceId"),
	}
}

// String returns a human-readable description of this Restorable Database Account ID
func (id RestorableDatabaseAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Instance: %q", id.InstanceId),
	}
	return fmt.Sprintf("Restorable Database Account (%s)", strings.Join(components, "\n"))
}
