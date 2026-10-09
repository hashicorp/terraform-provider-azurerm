// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2024-04-01/invoicesection"
)

func TestInvoiceSectionNamePattern(t *testing.T) {
	testData := []struct {
		Input    string
		Expected bool
	}{
		{Input: "", Expected: false},
		{Input: "a", Expected: true},
		{Input: "invoice-section_1", Expected: true},
		{Input: "12345", Expected: true},
		{Input: "invoice section", Expected: false},
		{Input: "invoice.section", Expected: false},
		{Input: strings.Repeat("a", 128), Expected: true},
		{Input: strings.Repeat("a", 129), Expected: false},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		if actual := invoiceSectionNamePattern.MatchString(v.Input); actual != v.Expected {
			t.Fatalf("expected %t but got %t for %q", v.Expected, actual, v.Input)
		}
	}
}

func TestExpandBillingInvoiceSection(t *testing.T) {
	tags := map[string]string{"costCategory": "Support"}

	testData := []struct {
		Name         string
		Input        BillingInvoiceSectionResourceModel
		ExpectedTags map[string]string
	}{
		{
			// an explicit empty map is required so removing `tags` from the config clears them
			Name:         "nil tags are sent as an empty map",
			Input:        BillingInvoiceSectionResourceModel{DisplayName: "example"},
			ExpectedTags: map[string]string{},
		},
		{
			Name:         "tags are sent to both locations",
			Input:        BillingInvoiceSectionResourceModel{DisplayName: "example", Tags: tags},
			ExpectedTags: tags,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual := expandBillingInvoiceSection(v.Input)

		if pointer.From(actual.Properties.DisplayName) != v.Input.DisplayName {
			t.Fatalf("expected `displayName` %q but got %q", v.Input.DisplayName, pointer.From(actual.Properties.DisplayName))
		}

		for name, got := range map[string]*map[string]string{
			"envelope":   actual.Tags,
			"properties": actual.Properties.Tags,
		} {
			if got == nil {
				t.Fatalf("expected the %s `tags` to be set but it was nil", name)
			}
			if len(*got) != len(v.ExpectedTags) {
				t.Fatalf("expected the %s `tags` to be %v but got %v", name, v.ExpectedTags, *got)
			}
			for k, expected := range v.ExpectedTags {
				if (*got)[k] != expected {
					t.Fatalf("expected the %s `tags` to contain %q=%q but got %q", name, k, expected, (*got)[k])
				}
			}
		}
	}
}

func TestFlattenBillingInvoiceSectionTags(t *testing.T) {
	populated := map[string]string{"costCategory": "Support"}
	envelopeOnly := map[string]string{"pcCode": "A123456"}

	testData := []struct {
		Name     string
		Input    invoicesection.InvoiceSection
		Expected map[string]string
	}{
		{
			Name:     "no tags anywhere",
			Input:    invoicesection.InvoiceSection{},
			Expected: nil,
		},
		{
			Name: "tags only within properties",
			Input: invoicesection.InvoiceSection{
				Properties: &invoicesection.InvoiceSectionProperties{Tags: pointer.To(populated)},
			},
			Expected: populated,
		},
		{
			Name: "tags only in the envelope",
			Input: invoicesection.InvoiceSection{
				Tags:       pointer.To(envelopeOnly),
				Properties: &invoicesection.InvoiceSectionProperties{},
			},
			Expected: envelopeOnly,
		},
		{
			// the service echoing an empty `properties.tags` must not mask a populated envelope
			Name: "empty properties tags falls back to the envelope",
			Input: invoicesection.InvoiceSection{
				Tags:       pointer.To(envelopeOnly),
				Properties: &invoicesection.InvoiceSectionProperties{Tags: pointer.To(map[string]string{})},
			},
			Expected: envelopeOnly,
		},
		{
			Name: "properties wins when both are populated",
			Input: invoicesection.InvoiceSection{
				Tags:       pointer.To(envelopeOnly),
				Properties: &invoicesection.InvoiceSectionProperties{Tags: pointer.To(populated)},
			},
			Expected: populated,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual := flattenBillingInvoiceSectionTags(v.Input)
		if len(actual) != len(v.Expected) {
			t.Fatalf("expected %v but got %v", v.Expected, actual)
		}
		for k, expected := range v.Expected {
			if actual[k] != expected {
				t.Fatalf("expected %q for key %q but got %q", expected, k, actual[k])
			}
		}
	}
}

func TestFlattenBillingInvoiceSectionEligibilityDetails(t *testing.T) {
	testData := []struct {
		Name     string
		Input    *[]invoicesection.DeleteInvoiceSectionEligibilityDetail
		Expected string
	}{
		{
			Name:     "nil",
			Input:    nil,
			Expected: "the Billing API didn't provide a reason",
		},
		{
			Name:     "empty",
			Input:    pointer.To([]invoicesection.DeleteInvoiceSectionEligibilityDetail{}),
			Expected: "the Billing API didn't provide a reason",
		},
		{
			Name: "single reason",
			Input: pointer.To([]invoicesection.DeleteInvoiceSectionEligibilityDetail{
				{
					Code:    pointer.To(invoicesection.DeleteInvoiceSectionEligibilityCodeActiveBillingSubscriptions),
					Message: pointer.To("There are active billing subscriptions."),
				},
			}),
			Expected: "ActiveBillingSubscriptions: There are active billing subscriptions.",
		},
		{
			Name: "multiple reasons",
			Input: pointer.To([]invoicesection.DeleteInvoiceSectionEligibilityDetail{
				{
					Code:    pointer.To(invoicesection.DeleteInvoiceSectionEligibilityCodeActiveBillingSubscriptions),
					Message: pointer.To("There are active billing subscriptions."),
				},
				{
					Code:    pointer.To(invoicesection.DeleteInvoiceSectionEligibilityCodeReservedInstances),
					Message: pointer.To("There are reserved instances."),
				},
			}),
			Expected: "ActiveBillingSubscriptions: There are active billing subscriptions., ReservedInstances: There are reserved instances.",
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		if actual := flattenBillingInvoiceSectionEligibilityDetails(v.Input); actual != v.Expected {
			t.Fatalf("expected %q but got %q", v.Expected, actual)
		}
	}
}
