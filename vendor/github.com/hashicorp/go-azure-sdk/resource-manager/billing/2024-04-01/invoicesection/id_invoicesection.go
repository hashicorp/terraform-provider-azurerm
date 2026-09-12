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
	recaser.RegisterResourceId(&InvoiceSectionId{})
}

var _ resourceids.ResourceId = &InvoiceSectionId{}

// InvoiceSectionId is a struct representing the Resource ID for a Invoice Section
type InvoiceSectionId struct {
	BillingAccountName string
	BillingProfileName string
	InvoiceSectionName string
}

// NewInvoiceSectionID returns a new InvoiceSectionId struct
func NewInvoiceSectionID(billingAccountName string, billingProfileName string, invoiceSectionName string) InvoiceSectionId {
	return InvoiceSectionId{
		BillingAccountName: billingAccountName,
		BillingProfileName: billingProfileName,
		InvoiceSectionName: invoiceSectionName,
	}
}

// ParseInvoiceSectionID parses 'input' into a InvoiceSectionId
func ParseInvoiceSectionID(input string) (*InvoiceSectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InvoiceSectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InvoiceSectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseInvoiceSectionIDInsensitively parses 'input' case-insensitively into a InvoiceSectionId
// note: this method should only be used for API response data and not user input
func ParseInvoiceSectionIDInsensitively(input string) (*InvoiceSectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InvoiceSectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InvoiceSectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *InvoiceSectionId) FromParseResult(input resourceids.ParseResult) error {
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

// ValidateInvoiceSectionID checks that 'input' can be parsed as a Invoice Section ID
func ValidateInvoiceSectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseInvoiceSectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Invoice Section ID
func (id InvoiceSectionId) ID() string {
	fmtString := "/providers/Microsoft.Billing/billingAccounts/%s/billingProfiles/%s/invoiceSections/%s"
	return fmt.Sprintf(fmtString, id.BillingAccountName, id.BillingProfileName, id.InvoiceSectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Invoice Section ID
func (id InvoiceSectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftBilling", "Microsoft.Billing", "Microsoft.Billing"),
		resourceids.StaticSegment("staticBillingAccounts", "billingAccounts", "billingAccounts"),
		resourceids.UserSpecifiedSegment("billingAccountName", "billingAccountName"),
		resourceids.StaticSegment("staticBillingProfiles", "billingProfiles", "billingProfiles"),
		resourceids.UserSpecifiedSegment("billingProfileName", "billingProfileName"),
		resourceids.StaticSegment("staticInvoiceSections", "invoiceSections", "invoiceSections"),
		resourceids.UserSpecifiedSegment("invoiceSectionName", "invoiceSectionName"),
	}
}

// String returns a human-readable description of this Invoice Section ID
func (id InvoiceSectionId) String() string {
	components := []string{
		fmt.Sprintf("Billing Account Name: %q", id.BillingAccountName),
		fmt.Sprintf("Billing Profile Name: %q", id.BillingProfileName),
		fmt.Sprintf("Invoice Section Name: %q", id.InvoiceSectionName),
	}
	return fmt.Sprintf("Invoice Section (%s)", strings.Join(components, "\n"))
}
