// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &BillingAccountInvoiceSectionId{}

// BillingAccountInvoiceSectionId is a struct representing the Resource ID for a Billing Account Invoice Section
type BillingAccountInvoiceSectionId struct {
	BillingAccountName string
	BillingProfileName string
	InvoiceSectionName string
}

// NewBillingAccountInvoiceSectionID returns a new BillingAccountInvoiceSectionId struct
func NewBillingAccountInvoiceSectionID(billingAccountName string, billingProfileName string, invoiceSectionName string) BillingAccountInvoiceSectionId {
	return BillingAccountInvoiceSectionId{
		BillingAccountName: billingAccountName,
		BillingProfileName: billingProfileName,
		InvoiceSectionName: invoiceSectionName,
	}
}

// ParseBillingAccountInvoiceSectionID parses 'input' into a BillingAccountInvoiceSectionId
func ParseBillingAccountInvoiceSectionID(input string) (*BillingAccountInvoiceSectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BillingAccountInvoiceSectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BillingAccountInvoiceSectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBillingAccountInvoiceSectionIDInsensitively parses 'input' case-insensitively into a BillingAccountInvoiceSectionId
// note: this method should only be used for API response data and not user input
func ParseBillingAccountInvoiceSectionIDInsensitively(input string) (*BillingAccountInvoiceSectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BillingAccountInvoiceSectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BillingAccountInvoiceSectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BillingAccountInvoiceSectionId) FromParseResult(input resourceids.ParseResult) error {

	var ok bool

	if id.BillingAccountName, ok = input.Parsed["billingAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "billingAccountName", input)
	}

	if id.BillingProfileName, ok = input.Parsed["billingProfileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "billingProfileName", input)
	}

	if id.InvoiceSectionName, ok = input.Parsed["invoiceSectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "invoiceSectionName", input)
	}

	return nil
}

// ValidateBillingAccountInvoiceSectionID checks that 'input' can be parsed as a Billing Account Invoice Section ID
func ValidateBillingAccountInvoiceSectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBillingAccountInvoiceSectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Billing Account Invoice Section ID
func (id BillingAccountInvoiceSectionId) ID() string {
	fmtString := "/providers/Microsoft.Billing/billingAccounts/%s/billingProfiles/%s/invoiceSections/%s"
	return fmt.Sprintf(fmtString, id.BillingAccountName, id.BillingProfileName, id.InvoiceSectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Billing Account Invoice Section ID
func (id BillingAccountInvoiceSectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Billing", "Microsoft.Billing"),
		resourceids.StaticSegment("billingAccounts", "billingAccounts", "billingAccounts"),
		resourceids.UserSpecifiedSegment("billingAccountName", "billingAccountValue"),
		resourceids.StaticSegment("billingProfiles", "billingProfiles", "billingProfiles"),
		resourceids.UserSpecifiedSegment("billingProfileName", "billingProfileValue"),
		resourceids.StaticSegment("invoiceSections", "invoiceSections", "invoiceSections"),
		resourceids.UserSpecifiedSegment("invoiceSectionName", "invoiceSectionValue"),
	}
}

// String returns a human-readable description of this Billing Account Invoice Section ID
func (id BillingAccountInvoiceSectionId) String() string {
	components := []string{
		fmt.Sprintf("Billing Account Name: %q", id.BillingAccountName),
		fmt.Sprintf("Billing Profile Name: %q", id.BillingProfileName),
		fmt.Sprintf("Invoice Section Name: %q", id.InvoiceSectionName),
	}
	return fmt.Sprintf("Billing Account Invoice Section (%s)", strings.Join(components, "\n"))
}
