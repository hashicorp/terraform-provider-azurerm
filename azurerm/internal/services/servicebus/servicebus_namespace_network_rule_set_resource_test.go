package servicebus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ServiceBusNamespaceNetworkRuleSetResource struct {
}

func TestAccServiceBusNamespaceNetworkRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_network_rule_set", "test")
	r := ServiceBusNamespaceNetworkRuleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusNamespaceNetworkRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_network_rule_set", "test")
	r := ServiceBusNamespaceNetworkRuleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusNamespaceNetworkRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_network_rule_set", "test")
	r := ServiceBusNamespaceNetworkRuleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusNamespaceNetworkRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_network_rule_set", "test")
	r := ServiceBusNamespaceNetworkRuleSetResource{}

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

func (t ServiceBusNamespaceNetworkRuleSetResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.NamespaceNetworkRuleSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceBus.NamespacesClientPreview.GetNetworkRuleSet(ctx, id.ResourceGroup, id.NamespaceName)
	if err != nil {
		return nil, fmt.Errorf("reading Service Bus NameSpace Network Rule Set (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r ServiceBusNamespaceNetworkRuleSetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace_network_rule_set" "test" {
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  default_action = "Deny"

  network_rules {
    subnet_id                            = azurerm_subnet.test.id
    ignore_missing_vnet_service_endpoint = false
  }
}
`, r.template(data))
}

func (r ServiceBusNamespaceNetworkRuleSetResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace_network_rule_set" "test" {
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  default_action = "Deny"

  network_rules {
    subnet_id                            = azurerm_subnet.test.id
    ignore_missing_vnet_service_endpoint = false
  }

  ip_rules = ["1.1.1.1"]
}
`, r.template(data))
}

func (ServiceBusNamespaceNetworkRuleSetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sb-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-sb-namespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Premium"

  capacity = 1
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["172.17.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}

resource "azurerm_subnet" "test" {
  name                 = "${azurerm_virtual_network.test.name}-default"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "172.17.0.0/24"

  service_endpoints = ["Microsoft.ServiceBus"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ServiceBusNamespaceNetworkRuleSetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace_network_rule_set" "import" {
  namespace_name      = azurerm_servicebus_namespace_network_rule_set.test.namespace_name
  resource_group_name = azurerm_servicebus_namespace_network_rule_set.test.resource_group_name
}
`, r.basic(data))
}
