package invoicesection

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BillingProfileId{})
}

var _ resourceids.ResourceId = &BillingProfileId{}

// BillingProfileId is a struct representing the Resource ID for a Billing Profile
type BillingProfileId struct {
	BillingAccountName string
	BillingProfileName string
}

// NewBillingProfileID returns a new BillingProfileId struct
func NewBillingProfileID(billingAccountName string, billingProfileName string) BillingProfileId {
	return BillingProfileId{
		BillingAccountName: billingAccountName,
		BillingProfileName: billingProfileName,
	}
}

// ParseBillingProfileID parses 'input' into a BillingProfileId
func ParseBillingProfileID(input string) (*BillingProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BillingProfileId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BillingProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBillingProfileIDInsensitively parses 'input' case-insensitively into a BillingProfileId
// note: this method should only be used for API response data and not user input
func ParseBillingProfileIDInsensitively(input string) (*BillingProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BillingProfileId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BillingProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BillingProfileId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BillingAccountName, ok = input.Parsed["billingAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "billingAccountName", input)
	}

	if id.BillingProfileName, ok = input.Parsed["billingProfileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "billingProfileName", input)
	}

	return nil
}

// ValidateBillingProfileID checks that 'input' can be parsed as a Billing Profile ID
func ValidateBillingProfileID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBillingProfileID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Billing Profile ID
func (id BillingProfileId) ID() string {
	fmtString := "/providers/Microsoft.Billing/billingAccounts/%s/billingProfiles/%s"
	return fmt.Sprintf(fmtString, id.BillingAccountName, id.BillingProfileName)
}

// Segments returns a slice of Resource ID Segments which comprise this Billing Profile ID
func (id BillingProfileId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftBilling", "Microsoft.Billing", "Microsoft.Billing"),
		resourceids.StaticSegment("staticBillingAccounts", "billingAccounts", "billingAccounts"),
		resourceids.UserSpecifiedSegment("billingAccountName", "billingAccountName"),
		resourceids.StaticSegment("staticBillingProfiles", "billingProfiles", "billingProfiles"),
		resourceids.UserSpecifiedSegment("billingProfileName", "billingProfileName"),
	}
}

// String returns a human-readable description of this Billing Profile ID
func (id BillingProfileId) String() string {
	components := []string{
		fmt.Sprintf("Billing Account Name: %q", id.BillingAccountName),
		fmt.Sprintf("Billing Profile Name: %q", id.BillingProfileName),
	}
	return fmt.Sprintf("Billing Profile (%s)", strings.Join(components, "\n"))
}
