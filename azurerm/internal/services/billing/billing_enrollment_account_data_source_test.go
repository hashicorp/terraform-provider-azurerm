package billing_test

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type BillingEnrollmentAccountDataSource struct{}

func TestAccBillingEnrollmentAccountDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_billing_enrollment_account_scope", "test")

	r := BillingEnrollmentAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("/providers/Microsoft.Billing/billingAccounts/12345678/enrollmentAccounts/123456"),
			),
		},
	})
}

func (BillingEnrollmentAccountDataSource) basic() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_billing_enrollment_account_scope" "test" {
  billing_account_name    = "12345678"
  enrollment_account_name = "123456"
}
`
}
