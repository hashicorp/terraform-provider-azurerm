package enrollmentaccounts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &EnrollmentAccountId{}

// EnrollmentAccountId is a struct representing the Resource ID for a Enrollment Account
type EnrollmentAccountId struct {
	BillingAccountName    string
	EnrollmentAccountName string
}

// NewEnrollmentAccountID returns a new EnrollmentAccountId struct
func NewEnrollmentAccountID(billingAccountName string, enrollmentAccountName string) EnrollmentAccountId {
	return EnrollmentAccountId{
		BillingAccountName:    billingAccountName,
		EnrollmentAccountName: enrollmentAccountName,
	}
}

// ParseEnrollmentAccountID parses 'input' into a EnrollmentAccountId
func ParseEnrollmentAccountID(input string) (*EnrollmentAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EnrollmentAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EnrollmentAccountId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseEnrollmentAccountIDInsensitively parses 'input' case-insensitively into a EnrollmentAccountId
// note: this method should only be used for API response data and not user input
func ParseEnrollmentAccountIDInsensitively(input string) (*EnrollmentAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EnrollmentAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EnrollmentAccountId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *EnrollmentAccountId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BillingAccountName, ok = input.Parsed["billingAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "billingAccountName", input)
	}

	if id.EnrollmentAccountName, ok = input.Parsed["enrollmentAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "enrollmentAccountName", input)
	}

	return nil
}

// ValidateEnrollmentAccountID checks that 'input' can be parsed as a Enrollment Account ID
func ValidateEnrollmentAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEnrollmentAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Enrollment Account ID
func (id EnrollmentAccountId) ID() string {
	fmtString := "/providers/Microsoft.Billing/billingAccounts/%s/enrollmentAccounts/%s"
	return fmt.Sprintf(fmtString, id.BillingAccountName, id.EnrollmentAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Enrollment Account ID
func (id EnrollmentAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftBilling", "Microsoft.Billing", "Microsoft.Billing"),
		resourceids.StaticSegment("staticBillingAccounts", "billingAccounts", "billingAccounts"),
		resourceids.UserSpecifiedSegment("billingAccountName", "billingAccountValue"),
		resourceids.StaticSegment("staticEnrollmentAccounts", "enrollmentAccounts", "enrollmentAccounts"),
		resourceids.UserSpecifiedSegment("enrollmentAccountName", "enrollmentAccountValue"),
	}
}

// String returns a human-readable description of this Enrollment Account ID
func (id EnrollmentAccountId) String() string {
	components := []string{
		fmt.Sprintf("Billing Account Name: %q", id.BillingAccountName),
		fmt.Sprintf("Enrollment Account Name: %q", id.EnrollmentAccountName),
	}
	return fmt.Sprintf("Enrollment Account (%s)", strings.Join(components, "\n"))
}
