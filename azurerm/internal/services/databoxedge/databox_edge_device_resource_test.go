package databoxedge_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databoxedge/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DataboxEdgeDeviceResource struct {
}

func TestAccDataboxEdgeDevice_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_device", "test")
	r := DataboxEdgeDeviceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	},
	)
}

func TestAccDataboxEdgeDevice_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_device", "test")
	r := DataboxEdgeDeviceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_databox_edge_device"),
		},
	},
	)
}

func TestAccDataboxEdgeDevice_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_device", "test")
	r := DataboxEdgeDeviceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	},
	)
}

func TestAccDataboxEdgeDevice_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_device", "test")
	r := DataboxEdgeDeviceResource{}

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
	},
	)
}

func (DataboxEdgeDeviceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DataboxEdgeDeviceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataboxEdge.DeviceClient.Get(ctx, id.Name, id.ResourceGroup)
	if err != nil {
		return nil, fmt.Errorf("retrieving Databox Edge Device Name %q (resource group: %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.DeviceProperties != nil), nil
}

// Location has to be hard coded due to limited support of locations for this resource
func (DataboxEdgeDeviceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databoxedge-%d"
  location = "%s"
}
`, data.RandomInteger, "eastus")
}

func (r DataboxEdgeDeviceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_databox_edge_device" "test" {
  name                = "acctest-dd-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku_name = "EdgeP_Base-Standard"
}
`, r.template(data), data.RandomString)
}

func (r DataboxEdgeDeviceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_databox_edge_device" "import" {
  name                = azurerm_databox_edge_device.test.name
  resource_group_name = azurerm_databox_edge_device.test.resource_group_name
  location            = azurerm_databox_edge_device.test.location

  sku_name = "EdgeP_Base-Standard"
}
`, r.basic(data))
}

func (r DataboxEdgeDeviceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_databox_edge_device" "test" {
  name                = "acctest-dd-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku_name = "EdgeP_Base-Standard"

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomString)
}
