package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/routingrulecollections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagerRoutingRuleCollectionResource struct{}

func testAccNetworkManagerRoutingRuleCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_routing_rule_collection", "test")
	r := ManagerRoutingRuleCollectionResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkManagerRoutingRuleCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_routing_rule_collection", "test")
	r := ManagerRoutingRuleCollectionResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
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

func testAccNetworkManagerRoutingRuleCollection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_routing_rule_collection", "test")
	r := ManagerRoutingRuleCollectionResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccNetworkManagerRoutingRuleCollection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_routing_rule_collection", "test")
	r := ManagerRoutingRuleCollectionResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ManagerRoutingRuleCollectionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := routingrulecollections.ParseRuleCollectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.RoutingRuleCollections.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagerRoutingRuleCollectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_network_manager_routing_rule_collection" "test" {
  name                     = "acctest-nmrrc-%[2]d"
  routing_configuration_id = azurerm_network_manager_routing_configuration.test.id
  network_group_ids        = [azurerm_network_manager_network_group.test.id]
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerRoutingRuleCollectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_routing_rule_collection" "import" {
  name                     = azurerm_network_manager_routing_rule_collection.test.name
  routing_configuration_id = azurerm_network_manager_routing_configuration.test.id
  network_group_ids        = [azurerm_network_manager_network_group.test.id]
}
`, r.basic(data))
}

func (r ManagerRoutingRuleCollectionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_network_manager_routing_rule_collection" "test" {
  name                          = "acctest-nmrrc-%[2]d"
  routing_configuration_id      = azurerm_network_manager_routing_configuration.test.id
  bgp_route_propagation_enabled = true
  network_group_ids             = [azurerm_network_manager_network_group.test.id, azurerm_network_manager_network_group.test2.id]
  description                   = "This is a test Routing Configuration"
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerRoutingRuleCollectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-manager-rc-%d"
  location = "%s"
}

data "azurerm_subscription" "current" {}

resource "azurerm_network_manager" "test" {
  name                = "acctest-nm-rc-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Routing"]
}

resource "azurerm_network_manager_network_group" "test" {
  name               = "acctest-nmng-1-%[1]d"
  network_manager_id = azurerm_network_manager.test.id
}

resource "azurerm_network_manager_network_group" "test2" {
  name               = "acctest-nmng-2-%[1]d"
  network_manager_id = azurerm_network_manager.test.id
}

resource "azurerm_network_manager_routing_configuration" "test" {
  name               = "acctest-nmrc-%[1]d"
  network_manager_id = azurerm_network_manager.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}
