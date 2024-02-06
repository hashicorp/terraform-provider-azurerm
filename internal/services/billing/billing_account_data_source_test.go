// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package billing_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type BillingAccountDataSource struct{}

func TestAccBillingAccountDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_billing_account", "test")
	r := BillingAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("/providers/Microsoft.Billing/billingAccounts/12345678"),
				check.That(data.ResourceName).Key("account_status").IsNotEmpty(),
				check.That(data.ResourceName).Key("account_type").IsNotEmpty(),
				check.That(data.ResourceName).Key("agreement_type").IsNotEmpty(),
				check.That(data.ResourceName).Key("display_name").IsNotEmpty(),
				check.That(data.ResourceName).Key("has_read_access").IsNotEmpty(),
				check.That(data.ResourceName).Key("sold_to.0.%").Exists(),
			),
		},
	})
}

func (BillingAccountDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_billing_account" "test" {
  name = "12345678"
}
`)
}
