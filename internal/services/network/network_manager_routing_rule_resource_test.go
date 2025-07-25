package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/routingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagerRoutingRuleResource struct{}

func testAccNetworkManagerRoutingRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_routing_rule", "test")
	r := ManagerRoutingRuleResource{}

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

func testAccNetworkManagerRoutingRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_routing_rule", "test")
	r := ManagerRoutingRuleResource{}

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

func testAccNetworkManagerRoutingRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_routing_rule", "test")
	r := ManagerRoutingRuleResource{}

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

func testAccNetworkManagerRoutingRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_routing_rule", "test")
	r := ManagerRoutingRuleResource{}

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

func (r ManagerRoutingRuleResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := routingrules.ParseRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.RoutingRules.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagerRoutingRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_network_manager_routing_rule" "test" {
  name               = "acctest-nmrr-%[2]d"
  rule_collection_id = azurerm_network_manager_routing_rule_collection.test.id
  destination {
    type    = "AddressPrefix"
    address = "10.0.0.0/16"
  }

  next_hop {
    type = "VirtualNetworkGateway"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerRoutingRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_routing_rule" "import" {
  name               = azurerm_network_manager_routing_rule.test.name
  rule_collection_id = azurerm_network_manager_routing_rule_collection.test.id
  destination {
    type    = azurerm_network_manager_routing_rule.test.destination[0].type
    address = azurerm_network_manager_routing_rule.test.destination[0].address
  }

  next_hop {
    type = azurerm_network_manager_routing_rule.test.next_hop[0].type
  }
}
`, r.basic(data))
}

func (r ManagerRoutingRuleResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_network_manager_routing_rule" "test" {
  name               = "acctest-nmrr-%[2]d"
  rule_collection_id = azurerm_network_manager_routing_rule_collection.test.id
  description        = "This is a test Routing Rule"
  destination {
    type    = "ServiceTag"
    address = "10.0.0.0/16"
  }

  next_hop {
    type    = "VirtualAppliance"
    address = "10.0.10.1"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerRoutingRuleResource) template(data acceptance.TestData) string {
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
  name               = "acctest-nmng-%[1]d"
  network_manager_id = azurerm_network_manager.test.id
}

resource "azurerm_network_manager_routing_configuration" "test" {
  name               = "acctest-nmrc-%[1]d"
  network_manager_id = azurerm_network_manager.test.id
}

resource "azurerm_network_manager_routing_rule_collection" "test" {
  name                     = "acctest-nmrrc-%[1]d"
  routing_configuration_id = azurerm_network_manager_routing_configuration.test.id
  network_group_ids        = ["azurerm_network_manager_network_group.test.id"]
  description              = "test routing rule collection"
}
`, data.RandomInteger, data.Locations.Primary)
}
