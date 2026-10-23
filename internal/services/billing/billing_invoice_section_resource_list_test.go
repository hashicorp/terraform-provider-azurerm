// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package billing_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccBillingInvoiceSection_list_basic(t *testing.T) {
	skipBillingInvoiceSection(t)

	r := BillingInvoiceSectionResource{}
	listResourceAddress := "azurerm_billing_invoice_section.list"

	data := acceptance.BuildTestData(t, "azurerm_billing_invoice_section", "test")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQueryByBillingProfileId(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 2),
				},
			},
		},
	})
}

func (r BillingInvoiceSectionResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_billing_invoice_section" "test" {
  count = 2

  name               = "acctestis-${count.index}-%d"
  billing_profile_id = local.billing_profile_id
  display_name       = "acctestis-${count.index}-%d"
}
`, r.template(), data.RandomInteger, data.RandomInteger)
}

func (BillingInvoiceSectionResource) basicQueryByBillingProfileId() string {
	return fmt.Sprintf(`
list "azurerm_billing_invoice_section" "list" {
  provider = azurerm
  config {
    billing_profile_id = "/providers/Microsoft.Billing/billingAccounts/%s/billingProfiles/%s"
  }
}
`, os.Getenv("ARM_BILLING_ACCOUNT"), os.Getenv("ARM_BILLING_PROFILE"))
}
