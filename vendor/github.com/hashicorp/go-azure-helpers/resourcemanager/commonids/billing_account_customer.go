// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &BillingAccountCustomerId{}

// BillingAccountCustomerId is a struct representing the Resource ID for a Billing Account Customer
type BillingAccountCustomerId struct {
	BillingAccountName string
	CustomerName       string
}

// NewBillingAccountCustomerID returns a new BillingAccountCustomerId struct
func NewBillingAccountCustomerID(billingAccountName string, customerName string) BillingAccountCustomerId {
	return BillingAccountCustomerId{
		BillingAccountName: billingAccountName,
		CustomerName:       customerName,
	}
}

// ParseBillingAccountCustomerID parses 'input' into a BillingAccountCustomerId
func ParseBillingAccountCustomerID(input string) (*BillingAccountCustomerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BillingAccountCustomerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BillingAccountCustomerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBillingAccountCustomerIDInsensitively parses 'input' case-insensitively into a BillingAccountCustomerId
// note: this method should only be used for API response data and not user input
func ParseBillingAccountCustomerIDInsensitively(input string) (*BillingAccountCustomerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BillingAccountCustomerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BillingAccountCustomerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BillingAccountCustomerId) FromParseResult(input resourceids.ParseResult) error {

	var ok bool

	if id.BillingAccountName, ok = input.Parsed["billingAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "billingAccountName", input)
	}

	if id.CustomerName, ok = input.Parsed["customerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "customerName", input)
	}

	return nil
}

// ValidateBillingAccountCustomerID checks that 'input' can be parsed as a Billing Account Customer ID
func ValidateBillingAccountCustomerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBillingAccountCustomerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Billing Account Customer ID
func (id BillingAccountCustomerId) ID() string {
	fmtString := "/providers/Microsoft.Billing/billingAccounts/%s/customers/%s"
	return fmt.Sprintf(fmtString, id.BillingAccountName, id.CustomerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Billing Account Customer ID
func (id BillingAccountCustomerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Billing", "Microsoft.Billing"),
		resourceids.StaticSegment("billingAccounts", "billingAccounts", "billingAccounts"),
		resourceids.UserSpecifiedSegment("billingAccountName", "billingAccountValue"),
		resourceids.StaticSegment("customers", "customers", "customers"),
		resourceids.UserSpecifiedSegment("customerName", "customerValue"),
	}
}

// String returns a human-readable description of this Billing Account Customer ID
func (id BillingAccountCustomerId) String() string {
	components := []string{
		fmt.Sprintf("Billing Account Name: %q", id.BillingAccountName),
		fmt.Sprintf("Customer Name: %q", id.CustomerName),
	}
	return fmt.Sprintf("Billing Account Customer (%s)", strings.Join(components, "\n"))
}
