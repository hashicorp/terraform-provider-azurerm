// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &BillingEnrollmentAccountId{}

// BillingEnrollmentAccountId is a struct representing the Resource ID for a Billing Enrollment Account
type BillingEnrollmentAccountId struct {
	EnrollmentAccountName string
}

// NewBillingEnrollmentAccountID returns a new BillingEnrollmentAccountId struct
func NewBillingEnrollmentAccountID(enrollmentAccountName string) BillingEnrollmentAccountId {
	return BillingEnrollmentAccountId{
		EnrollmentAccountName: enrollmentAccountName,
	}
}

// ParseBillingEnrollmentAccountID parses 'input' into a BillingEnrollmentAccountId
func ParseBillingEnrollmentAccountID(input string) (*BillingEnrollmentAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BillingEnrollmentAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BillingEnrollmentAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBillingEnrollmentAccountIDInsensitively parses 'input' case-insensitively into a BillingEnrollmentAccountId
// note: this method should only be used for API response data and not user input
func ParseBillingEnrollmentAccountIDInsensitively(input string) (*BillingEnrollmentAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BillingEnrollmentAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BillingEnrollmentAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BillingEnrollmentAccountId) FromParseResult(input resourceids.ParseResult) error {

	var ok bool

	if id.EnrollmentAccountName, ok = input.Parsed["enrollmentAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "enrollmentAccountName", input)
	}

	return nil
}

// ValidateBillingEnrollmentAccountID checks that 'input' can be parsed as a Billing Enrollment Account ID
func ValidateBillingEnrollmentAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBillingEnrollmentAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Billing Enrollment Account ID
func (id BillingEnrollmentAccountId) ID() string {
	fmtString := "/providers/Microsoft.Billing/enrollmentAccounts/%s"
	return fmt.Sprintf(fmtString, id.EnrollmentAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Billing Enrollment Account ID
func (id BillingEnrollmentAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Billing", "Microsoft.Billing"),
		resourceids.StaticSegment("enrollmentAccounts", "enrollmentAccounts", "enrollmentAccounts"),
		resourceids.UserSpecifiedSegment("enrollmentAccountName", "enrollmentAccountValue"),
	}
}

// String returns a human-readable description of this Billing Enrollment Account ID
func (id BillingEnrollmentAccountId) String() string {
	components := []string{
		fmt.Sprintf("Enrollment Account Name: %q", id.EnrollmentAccountName),
	}
	return fmt.Sprintf("Billing Enrollment Account (%s)", strings.Join(components, "\n"))
}
