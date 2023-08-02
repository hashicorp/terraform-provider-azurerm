package quotas

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = QuotaId{}

// QuotaId is a struct representing the Resource ID for a Quota
type QuotaId struct {
	SubscriptionId string
	LocationName   string
	QuotaName      string
}

// NewQuotaID returns a new QuotaId struct
func NewQuotaID(subscriptionId string, locationName string, quotaName string) QuotaId {
	return QuotaId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		QuotaName:      quotaName,
	}
}

// ParseQuotaID parses 'input' into a QuotaId
func ParseQuotaID(input string) (*QuotaId, error) {
	parser := resourceids.NewParserFromResourceIdType(QuotaId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := QuotaId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.QuotaName, ok = parsed.Parsed["quotaName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "quotaName", *parsed)
	}

	return &id, nil
}

// ParseQuotaIDInsensitively parses 'input' case-insensitively into a QuotaId
// note: this method should only be used for API response data and not user input
func ParseQuotaIDInsensitively(input string) (*QuotaId, error) {
	parser := resourceids.NewParserFromResourceIdType(QuotaId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := QuotaId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.QuotaName, ok = parsed.Parsed["quotaName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "quotaName", *parsed)
	}

	return &id, nil
}

// ValidateQuotaID checks that 'input' can be parsed as a Quota ID
func ValidateQuotaID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseQuotaID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Quota ID
func (id QuotaId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.LoadTestService/locations/%s/quotas/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.QuotaName)
}

// Segments returns a slice of Resource ID Segments which comprise this Quota ID
func (id QuotaId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLoadTestService", "Microsoft.LoadTestService", "Microsoft.LoadTestService"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationValue"),
		resourceids.StaticSegment("staticQuotas", "quotas", "quotas"),
		resourceids.UserSpecifiedSegment("quotaName", "quotaValue"),
	}
}

// String returns a human-readable description of this Quota ID
func (id QuotaId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Quota Name: %q", id.QuotaName),
	}
	return fmt.Sprintf("Quota (%s)", strings.Join(components, "\n"))
}
