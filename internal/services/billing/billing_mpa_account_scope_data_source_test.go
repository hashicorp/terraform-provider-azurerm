// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package billing_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type BillingMPAAccountDataSource struct{}

func TestAccBillingMPAAccountDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_billing_mpa_account_scope", "test")

	r := BillingMPAAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/customers/2281f543-7321-4cf9-1e23-edb4Oc31a31c"),
			),
		},
	})
}

func (BillingMPAAccountDataSource) basic() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_billing_mpa_account_scope" "test" {
  billing_account_name = "e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31"
  customer_name        = "2281f543-7321-4cf9-1e23-edb4Oc31a31c"
}
`
}
