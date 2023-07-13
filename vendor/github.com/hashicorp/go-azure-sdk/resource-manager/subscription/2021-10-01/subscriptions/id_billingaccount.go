package subscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = BillingAccountId{}

// BillingAccountId is a struct representing the Resource ID for a Billing Account
type BillingAccountId struct {
	BillingAccountId string
}

// NewBillingAccountID returns a new BillingAccountId struct
func NewBillingAccountID(billingAccountId string) BillingAccountId {
	return BillingAccountId{
		BillingAccountId: billingAccountId,
	}
}

// ParseBillingAccountID parses 'input' into a BillingAccountId
func ParseBillingAccountID(input string) (*BillingAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(BillingAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BillingAccountId{}

	if id.BillingAccountId, ok = parsed.Parsed["billingAccountId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "billingAccountId", *parsed)
	}

	return &id, nil
}

// ParseBillingAccountIDInsensitively parses 'input' case-insensitively into a BillingAccountId
// note: this method should only be used for API response data and not user input
func ParseBillingAccountIDInsensitively(input string) (*BillingAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(BillingAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BillingAccountId{}

	if id.BillingAccountId, ok = parsed.Parsed["billingAccountId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "billingAccountId", *parsed)
	}

	return &id, nil
}

// ValidateBillingAccountID checks that 'input' can be parsed as a Billing Account ID
func ValidateBillingAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBillingAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Billing Account ID
func (id BillingAccountId) ID() string {
	fmtString := "/providers/Microsoft.Billing/billingAccounts/%s"
	return fmt.Sprintf(fmtString, id.BillingAccountId)
}

// Segments returns a slice of Resource ID Segments which comprise this Billing Account ID
func (id BillingAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftBilling", "Microsoft.Billing", "Microsoft.Billing"),
		resourceids.StaticSegment("staticBillingAccounts", "billingAccounts", "billingAccounts"),
		resourceids.UserSpecifiedSegment("billingAccountId", "billingAccountIdValue"),
	}
}

// String returns a human-readable description of this Billing Account ID
func (id BillingAccountId) String() string {
	components := []string{
		fmt.Sprintf("Billing Account: %q", id.BillingAccountId),
	}
	return fmt.Sprintf("Billing Account (%s)", strings.Join(components, "\n"))
}
