// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package billing_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2024-04-01/invoicesection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type BillingInvoiceSectionResource struct{}

// invoice sections require a pre-existing Microsoft Customer Agreement billing account and billing
// profile, neither of which can be created by Terraform
func skipBillingInvoiceSection(t *testing.T) {
	if os.Getenv("ARM_BILLING_ACCOUNT") == "" || os.Getenv("ARM_BILLING_PROFILE") == "" {
		t.Skip("skipping tests - no `ARM_BILLING_ACCOUNT` and `ARM_BILLING_PROFILE` provided")
	}
}

func TestAccBillingInvoiceSection_basic(t *testing.T) {
	skipBillingInvoiceSection(t)

	data := acceptance.BuildTestData(t, "azurerm_billing_invoice_section", "test")
	r := BillingInvoiceSectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBillingInvoiceSection_requiresImport(t *testing.T) {
	skipBillingInvoiceSection(t)

	data := acceptance.BuildTestData(t, "azurerm_billing_invoice_section", "test")
	r := BillingInvoiceSectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccBillingInvoiceSection_complete(t *testing.T) {
	skipBillingInvoiceSection(t)

	data := acceptance.BuildTestData(t, "azurerm_billing_invoice_section", "test")
	r := BillingInvoiceSectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBillingInvoiceSection_update(t *testing.T) {
	skipBillingInvoiceSection(t)

	data := acceptance.BuildTestData(t, "azurerm_billing_invoice_section", "test")
	r := BillingInvoiceSectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (BillingInvoiceSectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := invoicesection.ParseInvoiceSectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Billing.InvoiceSectionClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (BillingInvoiceSectionResource) template() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

locals {
  billing_profile_id = "/providers/Microsoft.Billing/billingAccounts/%s/billingProfiles/%s"
}
`, os.Getenv("ARM_BILLING_ACCOUNT"), os.Getenv("ARM_BILLING_PROFILE"))
}

func (r BillingInvoiceSectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_billing_invoice_section" "test" {
  name               = "acctestis-%d"
  billing_profile_id = local.billing_profile_id
  display_name       = "acctestis-%d"
}
`, r.template(), data.RandomInteger, data.RandomInteger)
}

func (r BillingInvoiceSectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_billing_invoice_section" "import" {
  name               = azurerm_billing_invoice_section.test.name
  billing_profile_id = azurerm_billing_invoice_section.test.billing_profile_id
  display_name       = azurerm_billing_invoice_section.test.display_name
}
`, r.basic(data))
}

func (r BillingInvoiceSectionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_billing_invoice_section" "test" {
  name               = "acctestis-%d"
  billing_profile_id = local.billing_profile_id
  display_name       = "acctestis-updated-%d"

  tags = {
    costCategory = "Support"
    pcCode       = "A123456"
  }
}
`, r.template(), data.RandomInteger, data.RandomInteger)
}
