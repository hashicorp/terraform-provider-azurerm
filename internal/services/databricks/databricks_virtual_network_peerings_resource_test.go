package databricks_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2023-02-01/vnetpeering"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DatabricksVirtualNetworkPeeringResource struct{}

func TestAccDatabricksVirtualNetworkPeering_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_virtual_network_peering", "test")
	r := DatabricksVirtualNetworkPeeringResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("workspace_id"),
	})
}

func TestAccDatabricksVirtualNetworkPeering_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_virtual_network_peering", "test")
	r := DatabricksVirtualNetworkPeeringResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDatabricksVirtualNetworkPeering_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_virtual_network_peering", "test")
	r := DatabricksVirtualNetworkPeeringResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("workspace_id"),
		{
			Config: r.update(data, "standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("workspace_id"),
	})
}

func TestAccDatabricksVirtualNetworkPeering_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_virtual_network_peering", "test")
	r := DatabricksVirtualNetworkPeeringResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("workspace_id"),
	})
}

func TestAccDatabricksVirtualNetworkPeering_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_virtual_network_peering", "test")
	r := DatabricksVirtualNetworkPeeringResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("workspace_id"),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("workspace_id"),
	})
}

func (DatabricksVirtualNetworkPeeringResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := vnetpeering.ParseVirtualNetworkPeeringID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataBricks.VnetPeeringsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("making Read request on Databricks %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (DatabricksVirtualNetworkPeeringResource) update(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "remote" {
  name                = "acctest-vnet-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.1.0/24"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctest-workspace-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "%[3]s"
}

resource "azurerm_databricks_virtual_network_peering" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_databricks_workspace.test.id

  remote_address_space_prefixes = azurerm_virtual_network.remote.address_space
  remote_virtual_network_id     = azurerm_virtual_network.remote.id

  virtual_network_access_enabled = false
}

resource "azurerm_virtual_network_peering" "remote" {
  name                      = "to-acctest-%[1]d"
  resource_group_name       = azurerm_resource_group.test.name
  virtual_network_name      = azurerm_virtual_network.remote.name
  remote_virtual_network_id = azurerm_databricks_virtual_network_peering.test.virtual_network_id
}
`, data.RandomInteger, data.Locations.Primary, sku)
}

func (DatabricksVirtualNetworkPeeringResource) requiresImport(data acceptance.TestData) string {
	template := DatabricksVirtualNetworkPeeringResource{}.basic(data, "standard")
	return fmt.Sprintf(`
%s

resource "azurerm_databricks_virtual_network_peering" "import" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_databricks_workspace.test.id

  remote_address_space_prefixes = azurerm_virtual_network.remote.address_space
  remote_virtual_network_id     = azurerm_virtual_network.remote.id
}
`, template, data.RandomInteger)
}

func (DatabricksVirtualNetworkPeeringResource) basic(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "remote" {
  name                = "acctest-vnet-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.1.0/24"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctest-workspace-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "%[3]s"
}

resource "azurerm_databricks_virtual_network_peering" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_databricks_workspace.test.id

  remote_address_space_prefixes = azurerm_virtual_network.remote.address_space
  remote_virtual_network_id     = azurerm_virtual_network.remote.id
}

resource "azurerm_virtual_network_peering" "remote" {
  name                      = "to-acctest-%[1]d"
  resource_group_name       = azurerm_resource_group.test.name
  virtual_network_name      = azurerm_virtual_network.remote.name
  remote_virtual_network_id = azurerm_databricks_virtual_network_peering.test.virtual_network_id
}
`, data.RandomInteger, data.Locations.Primary, sku)
}

func (DatabricksVirtualNetworkPeeringResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "remote" {
  name                = "acctest-vnet-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.1.0/24"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctest-workspace-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"
}

resource "azurerm_databricks_virtual_network_peering" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_databricks_workspace.test.id

  remote_address_space_prefixes = azurerm_virtual_network.remote.address_space
  remote_virtual_network_id     = azurerm_virtual_network.remote.id

  virtual_network_access_enabled = true
  forwarded_traffic_enabled      = true
  gateway_transit_enabled        = false
  remote_gateways_enabled        = false
}

resource "azurerm_virtual_network_peering" "remote" {
  name                      = "to-acctest-%[1]d"
  resource_group_name       = azurerm_resource_group.test.name
  virtual_network_name      = azurerm_virtual_network.remote.name
  remote_virtual_network_id = azurerm_databricks_virtual_network_peering.test.virtual_network_id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (DatabricksVirtualNetworkPeeringResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "remote" {
  name                = "acctest-vnet-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.1.0/24"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctest-workspace-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"
}

resource "azurerm_databricks_virtual_network_peering" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_databricks_workspace.test.id

  remote_address_space_prefixes = azurerm_virtual_network.remote.address_space
  remote_virtual_network_id     = azurerm_virtual_network.remote.id

  virtual_network_access_enabled = false
  forwarded_traffic_enabled      = false
  gateway_transit_enabled        = false
  remote_gateways_enabled        = false
}

resource "azurerm_virtual_network_peering" "remote" {
  name                      = "to-acctest-%[1]d"
  resource_group_name       = azurerm_resource_group.test.name
  virtual_network_name      = azurerm_virtual_network.remote.name
  remote_virtual_network_id = azurerm_databricks_virtual_network_peering.test.virtual_network_id
}
`, data.RandomInteger, data.Locations.Primary)
}
