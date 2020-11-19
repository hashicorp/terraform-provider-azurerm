package automation_test

import (
	`context`
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	`github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure`
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	`github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appconfiguration/parse`
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AutomationAccountResource struct {
}

func TestAccAzureRMAutomationAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_account", "test")
	r := AutomationAccountResource {}

	data.ResourceTest(t, r, []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key( "sku_name").HasValue( "Basic"),
					check.That(data.ResourceName).Key( "dsc_server_endpoint").Exists(),
					check.That(data.ResourceName).Key( "dsc_primary_access_key").Exists(),
					check.That(data.ResourceName).Key( "dsc_secondary_access_key").Exists(),
				),
			},
			data.ImportStep(),
	})
}

func TestAccAzureRMAutomationAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_account", "test")
	r := AutomationAccountResource {}

	data.ResourceTest(t, r, []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			{
				Config: r.requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_automation_account"),
			},
	})
}

func TestAccAzureRMAutomationAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_account", "test")
	r := AutomationAccountResource {}

	data.ResourceTest(t, r, []resource.TestStep{
			{
				Config: r.complete(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key( "sku_name").HasValue( "Basic"),
					check.That(data.ResourceName).Key( "tags.hello").HasValue( "world"),
				),
			},
			data.ImportStep(),
	})
}


func (t AutomationAccountResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	name := id.Path["automationAccounts"]

	resp, err := clients.Automation.AccountClient.Get(ctx, id.ResourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Automation Account %q (resource group: %q): %+v", name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.AccountProperties != nil), nil
}

func (AutomationAccountResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AutomationAccountResource) requiresImport(data acceptance.TestData) string {
	template := AutomationAccountResource{}.basic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_automation_account" "import" {
  name                = azurerm_automation_account.test.name
  location            = azurerm_automation_account.test.location
  resource_group_name = azurerm_automation_account.test.resource_group_name

  sku_name = "Basic"
}
`, template)
}

func (AutomationAccountResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Basic"

  tags = {
    "hello" = "world"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
