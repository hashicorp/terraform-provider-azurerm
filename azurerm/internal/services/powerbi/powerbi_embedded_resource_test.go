package powerbi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/powerbi/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PowerBIEmbeddedResource struct {
}

func TestAccPowerBIEmbedded_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_powerbi_embedded", "test")
	r := PowerBIEmbeddedResource{}

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

func TestAccPowerBIEmbedded_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_powerbi_embedded", "test")
	r := PowerBIEmbeddedResource{}

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

func TestAccPowerBIEmbedded_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_powerbi_embedded", "test")
	r := PowerBIEmbeddedResource{}

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

func TestAccPowerBIEmbedded_gen2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_powerbi_embedded", "test")
	r := PowerBIEmbeddedResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gen2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mode").HasValue("Gen2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPowerBIEmbedded_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_powerbi_embedded", "test")
	r := PowerBIEmbeddedResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_powerbi_embedded"),
		},
	})
}

func (PowerBIEmbeddedResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.EmbeddedID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.PowerBI.CapacityClient.GetDetails(ctx, id.ResourceGroup, id.CapacityName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.DedicatedCapacityProperties != nil), nil
}

func (PowerBIEmbeddedResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-powerbi-%[1]d"
  location = "%[2]s"
}

data "azurerm_client_config" "test" {}
`, data.RandomInteger, data.Locations.Primary)
}

func (r PowerBIEmbeddedResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_powerbi_embedded" "test" {
  name                = "acctestpowerbi%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "A1"
  administrators      = [data.azurerm_client_config.test.object_id]
}
`, r.template(data), data.RandomInteger)
}

func (r PowerBIEmbeddedResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_powerbi_embedded" "test" {
  name                = "acctestpowerbi%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "A2"
  administrators      = [data.azurerm_client_config.test.object_id]

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PowerBIEmbeddedResource) gen2(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_powerbi_embedded" "test" {
  name                = "acctestpowerbi%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "A1"
  administrators      = [data.azurerm_client_config.test.object_id]
  mode                = "Gen2"
}
`, r.template(data), data.RandomInteger)
}

func (r PowerBIEmbeddedResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_powerbi_embedded" "import" {
  name                = azurerm_powerbi_embedded.test.name
  location            = azurerm_powerbi_embedded.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "A1"
  administrators      = [data.azurerm_client_config.test.object_id]
}
`, r.basic(data))
}
