// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databricks_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2024-05-01/vnetpeering"
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
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksVirtualNetworkPeering_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_virtual_network_peering", "test")
	r := DatabricksVirtualNetworkPeeringResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
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
		data.ImportStep(),
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
		data.ImportStep(),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// This is needed for subscriptions that have AzSecPack applied on them by policy,
// tests will fail due to NSG's being automatically created for the resources
// without the tests knowledge causing the delete to fail...
func checkAzSecPackOverride() string {
	if os.Getenv("ARM_TEST_AZSECPACK_OVERRIDE") != "" {
		return (`
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  `)
	}

	return ""
}

func (DatabricksVirtualNetworkPeeringResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := vnetpeering.ParseVirtualNetworkPeeringID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataBricks.VnetPeeringClient.Get(ctx, *id)

	if response.WasNotFound(resp.HttpResponse) {
		return pointer.To(false), nil
	}

	if err != nil {
		return nil, fmt.Errorf("making Read request on Databricks %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (DatabricksVirtualNetworkPeeringResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
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
  name                = "acctest-ws-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DatabricksVirtualNetworkPeeringResource) update(data acceptance.TestData) string {
	features := "features {}"
	if override := checkAzSecPackOverride(); override != "" {
		features = fmt.Sprintf("features {%s}", override)
	}

	return fmt.Sprintf(`
provider "azurerm" {
  %s
}

%s

resource "azurerm_databricks_virtual_network_peering" "test" {
  name                = "acctest-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_databricks_workspace.test.id

  remote_address_space_prefixes = azurerm_virtual_network.remote.address_space
  remote_virtual_network_id     = azurerm_virtual_network.remote.id

  allow_virtual_network_access = false
}

resource "azurerm_virtual_network_peering" "remote" {
  name                      = "to-acctest-%[3]d"
  resource_group_name       = azurerm_resource_group.test.name
  virtual_network_name      = azurerm_virtual_network.remote.name
  remote_virtual_network_id = azurerm_databricks_virtual_network_peering.test.virtual_network_id
}
`, features, r.template(data), data.RandomInteger)
}

func (r DatabricksVirtualNetworkPeeringResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
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

func (r DatabricksVirtualNetworkPeeringResource) basic(data acceptance.TestData) string {
	features := "features {}"
	if override := checkAzSecPackOverride(); override != "" {
		features = fmt.Sprintf("features {%s}", override)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  %s
}

%s

resource "azurerm_databricks_virtual_network_peering" "test" {
  name                = "acctest-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_databricks_workspace.test.id

  remote_address_space_prefixes = azurerm_virtual_network.remote.address_space
  remote_virtual_network_id     = azurerm_virtual_network.remote.id
}

resource "azurerm_virtual_network_peering" "remote" {
  name                      = "to-acctest-%[3]d"
  resource_group_name       = azurerm_resource_group.test.name
  virtual_network_name      = azurerm_virtual_network.remote.name
  remote_virtual_network_id = azurerm_databricks_virtual_network_peering.test.virtual_network_id
}
`, features, r.template(data), data.RandomInteger)
}

func (r DatabricksVirtualNetworkPeeringResource) complete(data acceptance.TestData) string {
	features := "features {}"
	if override := checkAzSecPackOverride(); override != "" {
		features = fmt.Sprintf("features {%s}", override)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  %s
}

%s

resource "azurerm_databricks_virtual_network_peering" "test" {
  name                = "acctest-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_databricks_workspace.test.id

  remote_address_space_prefixes = azurerm_virtual_network.remote.address_space
  remote_virtual_network_id     = azurerm_virtual_network.remote.id

  allow_virtual_network_access = true
  allow_forwarded_traffic      = true
  allow_gateway_transit        = false
  use_remote_gateways          = false
}

resource "azurerm_virtual_network_peering" "remote" {
  name                      = "to-acctest-%[3]d"
  resource_group_name       = azurerm_resource_group.test.name
  virtual_network_name      = azurerm_virtual_network.remote.name
  remote_virtual_network_id = azurerm_databricks_virtual_network_peering.test.virtual_network_id
}
`, features, r.template(data), data.RandomInteger)
}

func (r DatabricksVirtualNetworkPeeringResource) completeUpdate(data acceptance.TestData) string {
	features := "features {}"
	if override := checkAzSecPackOverride(); override != "" {
		features = fmt.Sprintf("features {%s}", override)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  %s
}

%s

resource "azurerm_databricks_virtual_network_peering" "test" {
  name                = "acctest-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_databricks_workspace.test.id

  remote_address_space_prefixes = azurerm_virtual_network.remote.address_space
  remote_virtual_network_id     = azurerm_virtual_network.remote.id

  allow_virtual_network_access = false
  allow_forwarded_traffic      = false
  allow_gateway_transit        = false
  use_remote_gateways          = false
}

resource "azurerm_virtual_network_peering" "remote" {
  name                      = "to-acctest-%[3]d"
  resource_group_name       = azurerm_resource_group.test.name
  virtual_network_name      = azurerm_virtual_network.remote.name
  remote_virtual_network_id = azurerm_databricks_virtual_network_peering.test.virtual_network_id
}
`, features, r.template(data), data.RandomInteger)
}
