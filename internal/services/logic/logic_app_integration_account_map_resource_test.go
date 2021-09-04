package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogicAppIntegrationAccountMapResource struct{}

func TestAccLogicAppIntegrationAccountMap_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_map", "test")
	r := LogicAppIntegrationAccountMapResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content"), // not returned from the API
	})
}

func TestAccLogicAppIntegrationAccountMap_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_map", "test")
	r := LogicAppIntegrationAccountMapResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccLogicAppIntegrationAccountMap_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_map", "test")
	r := LogicAppIntegrationAccountMapResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content"), // not returned from the API
	})
}

func TestAccLogicAppIntegrationAccountMap_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_map", "test")
	r := LogicAppIntegrationAccountMapResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content"), // not returned from the API
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content"), // not returned from the API
	})
}

func (r LogicAppIntegrationAccountMapResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IntegrationAccountMapID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Logic.IntegrationAccountMapClient.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.MapName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q %+v", id, err)
	}

	return utils.Bool(resp.IntegrationAccountMapProperties != nil), nil
}

func (r LogicAppIntegrationAccountMapResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%d"
  location = "%s"
}

resource "azurerm_logic_app_integration_account" "test" {
  name                = "acctest-ia-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LogicAppIntegrationAccountMapResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_map" "test" {
  name                     = "acctest-iamap-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  map_type                 = "Xslt"
  content                  = file("testdata/integration_account_map_content.xsd")
}
`, r.template(data), data.RandomInteger)
}

func (r LogicAppIntegrationAccountMapResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_map" "import" {
  name                     = azurerm_logic_app_integration_account_map.test.name
  resource_group_name      = azurerm_logic_app_integration_account_map.test.resource_group_name
  integration_account_name = azurerm_logic_app_integration_account_map.test.integration_account_name
  map_type                 = azurerm_logic_app_integration_account_map.test.map_type
  content                  = azurerm_logic_app_integration_account_map.test.content
}
`, r.basic(data))
}

func (r LogicAppIntegrationAccountMapResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_map" "test" {
  name                     = "acctest-iamap-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  map_type                 = "Xslt"
  content                  = file("testdata/integration_account_map_content.xsd")

  metadata = {
    foo = "bar"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LogicAppIntegrationAccountMapResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_map" "test" {
  name                     = "acctest-iamap-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  map_type                 = "Xslt20"
  content                  = file("testdata/integration_account_map_content2.xsd")

  metadata = {
    foo = "bar2"
  }
}
`, r.template(data), data.RandomInteger)
}
