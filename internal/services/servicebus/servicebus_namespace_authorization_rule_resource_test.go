// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespacesauthorizationrule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServiceBusNamespaceAuthorizationRuleResource struct{}

func TestAccServiceBusNamespaceAuthorizationRule_listen(t *testing.T) {
	testAccServiceBusNamespaceAuthorizationRule(t, true, false, false)
}

func TestAccServiceBusNamespaceAuthorizationRule_send(t *testing.T) {
	testAccServiceBusNamespaceAuthorizationRule(t, false, true, false)
}

func TestAccServiceBusNamespaceAuthorizationRule_listensend(t *testing.T) {
	testAccServiceBusNamespaceAuthorizationRule(t, true, true, false)
}

func TestAccServiceBusNamespaceAuthorizationRule_manage(t *testing.T) {
	testAccServiceBusNamespaceAuthorizationRule(t, true, true, true)
}

func testAccServiceBusNamespaceAuthorizationRule(t *testing.T, listen, send, manage bool) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_authorization_rule", "test")
	r := ServiceBusNamespaceAuthorizationRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.base(data, listen, send, manage),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusNamespaceAuthorizationRule_rightsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_authorization_rule", "test")
	r := ServiceBusNamespaceAuthorizationRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.base(data, true, false, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("send").HasValue("false"),
				check.That(data.ResourceName).Key("manage").HasValue("false"),
			),
		},
		{
			Config: r.base(data, true, true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("namespace_id").Exists(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string_alias").HasValue(""),
				check.That(data.ResourceName).Key("secondary_connection_string_alias").HasValue(""),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("send").HasValue("true"),
				check.That(data.ResourceName).Key("manage").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusNamespaceAuthorizationRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_authorization_rule", "test")
	r := ServiceBusNamespaceAuthorizationRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.base(data, true, false, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("send").HasValue("false"),
				check.That(data.ResourceName).Key("manage").HasValue("false"),
			),
		},
		{
			Config:      r.requiresImport(data, true, false, false),
			ExpectError: acceptance.RequiresImportError("azurerm_servicebus_namespace_authorization_rule"),
		},
	})
}

func (t ServiceBusNamespaceAuthorizationRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := namespacesauthorizationrule.ParseAuthorizationRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceBus.NamespacesAuthClient.NamespacesGetAuthorizationRule(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (ServiceBusNamespaceAuthorizationRuleResource) base(data acceptance.TestData, listen, send, manage bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_namespace_authorization_rule" "test" {
  name         = "acctest-%[1]d"
  namespace_id = azurerm_servicebus_namespace.test.id

  listen = %[3]t
  send   = %[4]t
  manage = %[5]t
}
`, data.RandomInteger, data.Locations.Primary, listen, send, manage)
}

func (r ServiceBusNamespaceAuthorizationRuleResource) requiresImport(data acceptance.TestData, listen, send, manage bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace_authorization_rule" "import" {
  name         = azurerm_servicebus_namespace_authorization_rule.test.name
  namespace_id = azurerm_servicebus_namespace_authorization_rule.test.namespace_id

  listen = azurerm_servicebus_namespace_authorization_rule.test.listen
  send   = azurerm_servicebus_namespace_authorization_rule.test.send
  manage = azurerm_servicebus_namespace_authorization_rule.test.manage
}
`, r.base(data, listen, send, manage))
}

func (ServiceBusNamespaceAuthorizationRuleResource) withAliasConnectionString(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "primary" {
  name     = "acctest1RG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "secondary" {
  name     = "acctest2RG-%[1]d"
  location = "%[3]s"
}

resource "azurerm_servicebus_namespace" "primary_namespace_test" {
  name                         = "acctest1-%[1]d"
  location                     = azurerm_resource_group.primary.location
  resource_group_name          = azurerm_resource_group.primary.name
  sku                          = "Premium"
  premium_messaging_partitions = 1
  capacity                     = "1"
}

resource "azurerm_servicebus_namespace" "secondary_namespace_test" {
  name                         = "acctest2-%[1]d"
  location                     = azurerm_resource_group.secondary.location
  resource_group_name          = azurerm_resource_group.secondary.name
  sku                          = "Premium"
  premium_messaging_partitions = 1
  capacity                     = "1"
}

resource "azurerm_servicebus_namespace_disaster_recovery_config" "pairing_test" {
  name                 = "acctest-alias-%[1]d"
  primary_namespace_id = azurerm_servicebus_namespace.primary_namespace_test.id
  partner_namespace_id = azurerm_servicebus_namespace.secondary_namespace_test.id
}

resource "azurerm_servicebus_namespace_authorization_rule" "test" {
  name         = "namespace_rule_test"
  namespace_id = azurerm_servicebus_namespace.primary_namespace_test.id
  listen       = true
  send         = true
  manage       = true

  depends_on = [
    azurerm_servicebus_namespace_disaster_recovery_config.pairing_test
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}
