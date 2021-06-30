package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type AutomationVariableDateTimeResource struct {
}

func TestAccAutomationVariableDateTime_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_datetime", "test")
	r := AutomationVariableDateTimeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("value").HasValue("2019-04-24T21:40:54.074Z"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationVariableDateTime_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_datetime", "test")
	r := AutomationVariableDateTimeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("This variable is created by Terraform acceptance test."),
				check.That(data.ResourceName).Key("value").HasValue("2019-04-20T08:40:04.02Z"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationVariableDateTime_basicCompleteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_datetime", "test")
	r := AutomationVariableDateTimeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("value").HasValue("2019-04-24T21:40:54.074Z"),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("This variable is created by Terraform acceptance test."),
				check.That(data.ResourceName).Key("value").HasValue("2019-04-20T08:40:04.02Z"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("value").HasValue("2019-04-24T21:40:54.074Z"),
			),
		},
	})
}

func (t AutomationVariableDateTimeResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	return testCheckAzureRMAutomationVariableExists(ctx, clients, state, "Datetime")
}

func (AutomationVariableDateTimeResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_variable_datetime" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  value                   = "2019-04-24T21:40:54.074Z"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (AutomationVariableDateTimeResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_variable_datetime" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  description             = "This variable is created by Terraform acceptance test."
  value                   = "2019-04-20T08:40:04.02Z"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
