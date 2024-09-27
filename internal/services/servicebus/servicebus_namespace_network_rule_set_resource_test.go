// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2022-10-01-preview/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServiceBusNamespaceNetworkRuleSetResource struct{}

func TestAccServiceBusNamespaceNetworkRule_basic(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("azurerm_servicebus_namespace_network_rule_set is deprecated for 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_network_rule_set", "test")
	r := ServiceBusNamespaceNetworkRuleSetResource{}

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

func TestAccServiceBusNamespaceNetworkRule_complete(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("azurerm_servicebus_namespace_network_rule_set is deprecated for 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_network_rule_set", "test")
	r := ServiceBusNamespaceNetworkRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("trusted_services_allowed").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusNamespaceNetworkRule_update(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("azurerm_servicebus_namespace_network_rule_set is deprecated for 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_network_rule_set", "test")
	r := ServiceBusNamespaceNetworkRuleSetResource{}

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

func TestAccServiceBusNamespaceNetworkRule_requiresImport(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("azurerm_servicebus_namespace_network_rule_set is deprecated for 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_network_rule_set", "test")
	r := ServiceBusNamespaceNetworkRuleSetResource{}

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

func (t ServiceBusNamespaceNetworkRuleSetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := namespaces.ParseNamespaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceBus.NamespacesClient.GetNetworkRuleSet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ServiceBusNamespaceNetworkRuleSetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace_network_rule_set" "test" {
  namespace_id = azurerm_servicebus_namespace.test.id

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
  namespace_id = azurerm_servicebus_namespace.test.id

  default_action                = "Deny"
  trusted_services_allowed      = true
  public_network_access_enabled = true

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
  name                         = "acctest-sb-namespace-%[1]d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  sku                          = "Premium"
  premium_messaging_partitions = 1
  capacity                     = 1
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
  address_prefixes     = ["172.17.0.0/24"]

  service_endpoints = ["Microsoft.ServiceBus"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ServiceBusNamespaceNetworkRuleSetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace_network_rule_set" "import" {
  namespace_id = azurerm_servicebus_namespace_network_rule_set.test.namespace_id
}
`, r.basic(data))
}
