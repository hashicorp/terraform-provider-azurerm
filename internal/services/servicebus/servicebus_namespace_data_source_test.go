// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ServiceBusNamespaceDataSource struct{}

func TestAccDataSourceServiceBusNamespace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("sku").Exists(),
				check.That(data.ResourceName).Key("capacity").Exists(),
				check.That(data.ResourceName).Key("local_auth_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("minimum_tls_version").HasValue("1.2"),
				check.That(data.ResourceName).Key("default_primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("default_secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("default_primary_key").Exists(),
				check.That(data.ResourceName).Key("default_secondary_key").Exists(),
				check.That(data.ResourceName).Key("endpoint").Exists(),
			),
		},
	})
}

func TestAccDataSourceServiceBusNamespace_premium(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.premium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("sku").Exists(),
				check.That(data.ResourceName).Key("capacity").Exists(),
				check.That(data.ResourceName).Key("default_primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("default_secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("default_primary_key").Exists(),
				check.That(data.ResourceName).Key("default_secondary_key").Exists(),
				check.That(data.ResourceName).Key("endpoint").Exists(),
			),
		},
	})
}

func TestAccDataSourceServiceBusNamespace_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
	})
}

func TestAccDataSourceServiceBusNamespace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("local_auth_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccDataSourceServiceBusNamespace_publicNetworkAccess(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.publicNetworkAccess(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccDataSourceServiceBusNamespace_networkRuleSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.networkRuleSet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("network_rule_set.#").HasValue("1"),
			),
		},
	})
}

func TestAccDataSourceServiceBusNamespace_customerManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.customerManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("customer_managed_key.#").HasValue("1"),
			),
		},
	})
}

func (ServiceBusNamespaceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace" "test" {
  name                = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, ServicebusNamespaceResource{}.basic(data))
}

func (ServiceBusNamespaceDataSource) premium(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace" "test" {
  name                = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, ServicebusNamespaceResource{}.premium(data))
}

func (ServiceBusNamespaceDataSource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace" "test" {
  name                = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, ServicebusNamespaceResource{}.identitySystemAssigned(data))
}

func (ServiceBusNamespaceDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace" "test" {
  name                = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, ServicebusNamespaceResource{}.complete(data))
}

func (ServiceBusNamespaceDataSource) publicNetworkAccess(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace" "test" {
  name                = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, ServicebusNamespaceResource{}.publicNetworkAccessUpdate(data))
}

func (ServiceBusNamespaceDataSource) networkRuleSet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace" "test" {
  name                = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, ServicebusNamespaceResource{}.networkRuleSet(data))
}

func (ServiceBusNamespaceDataSource) customerManagedKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace" "test" {
  name                = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, ServicebusNamespaceResource{}.customerManagedKey(data))
}
