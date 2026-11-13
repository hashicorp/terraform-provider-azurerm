// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package billing_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type BillingInvoiceSectionDataSource struct{}

func TestAccBillingInvoiceSectionDataSource_basic(t *testing.T) {
	skipBillingInvoiceSection(t)

	data := acceptance.BuildTestData(t, "data.azurerm_billing_invoice_section", "test")
	d := BillingInvoiceSectionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("system_id").Exists(),
				check.That(data.ResourceName).Key("state").Exists(),
			),
		},
	})
}

func (BillingInvoiceSectionDataSource) basic(data acceptance.TestData) string {
	r := BillingInvoiceSectionResource{}

	return fmt.Sprintf(`
%s

data "azurerm_billing_invoice_section" "test" {
  name               = azurerm_billing_invoice_section.test.name
  billing_profile_id = azurerm_billing_invoice_section.test.billing_profile_id
}
`, r.basic(data))
}
