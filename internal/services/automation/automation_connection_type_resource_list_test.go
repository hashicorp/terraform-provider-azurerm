package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccNetworkProfile_list_basic(t *testing.T) {
	r := AutomationConnectionTypeResource{}

	data := acceptance.BuildTestData(t, "azurerm_automation_connection_type", "test1")

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
				Query:             true,
				Config:            r.basicQueryByAutomationAccount(data),
				QueryResultChecks: []querycheck.QueryResultCheck{}, // TODO
			},
		},
	})
}

func (r AutomationConnectionTypeResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_connection_type" "test1" {
  name                    = "acctest1-%[1]d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  is_global               = false
  field {
    name = "my_def"
    type = "string"
  }
}

resource "azurerm_automation_connection_type" "test2" {
  name                    = "acctest2-%[1]d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  is_global               = false
  field {
    name = "my_def"
    type = "string"
  }
}

resource "azurerm_automation_connection_type" "test3" {
  name                    = "acctest3-%[1]d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  is_global               = false
  field {
    name = "my_def"
    type = "string"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r AutomationConnectionTypeResource) basicQueryByAutomationAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_automation_connection_type" "list" {
  provider = azurerm
  config {
    resource_group_name     = "acctestRG-%[1]d"
    automation_account_name = "acctest-%[1]d"
  }
}
`, data.RandomInteger)
}
