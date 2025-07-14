// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package servicebus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2022-10-01-preview/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ServicebusNamespaceCustomerManagedKeyResource struct{}

func TestAccServicebusNamespaceCustomerManagedKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_customer_managed_key", "test")
	r := ServicebusNamespaceCustomerManagedKeyResource{}
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

func TestAccServicebusNamespaceCustomerManagedKey_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_customer_managed_key", "test")
	r := ServicebusNamespaceCustomerManagedKeyResource{}
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

func TestAccServicebusNamespaceCustomerManagedKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_customer_managed_key", "test")
	r := ServicebusNamespaceCustomerManagedKeyResource{}

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

func TestAccServicebusNamespaceCustomerManagedKey_updated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_customer_managed_key", "test")
	r := ServicebusNamespaceCustomerManagedKeyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ServicebusNamespaceCustomerManagedKeyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := namespaces.ParseNamespaceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.ServiceBus.NamespacesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil || resp.Model.Properties.Encryption == nil {
		return pointer.To(false), nil
	}
	return pointer.To(true), nil
}

func (r ServicebusNamespaceCustomerManagedKeyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-sb-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                         = "acctestservicebusnamespace-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  sku                          = "Premium"
  premium_messaging_partitions = 1
  capacity                     = 1

  identity {
    type = "SystemAssigned"
  }

  lifecycle {
    ignore_changes = [customer_managed_key]
  }
}


data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                        = "acctest-key-vault-%s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Decrypt",
      "Encrypt",
      "Delete",
      "Get",
      "List",
      "Purge",
      "UnwrapKey",
      "WrapKey",
      "Verify",
      "GetRotationPolicy"
    ]
    secret_permissions = [
      "Set",
    ]
  }

  access_policy {
    tenant_id = azurerm_servicebus_namespace.test.identity[0].tenant_id
    object_id = azurerm_servicebus_namespace.test.identity[0].principal_id

    key_permissions = [
      "Create",
      "Decrypt",
      "Encrypt",
      "Delete",
      "Get",
      "List",
      "Purge",
      "UnwrapKey",
      "WrapKey",
      "Verify",
      "GetRotationPolicy"
    ]
    secret_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkey-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomString)
}

func (r ServicebusNamespaceCustomerManagedKeyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace_customer_managed_key" "test" {
  namespace_id     = azurerm_servicebus_namespace.test.id
  key_vault_key_id = azurerm_key_vault_key.test.id
}
`, r.template(data))
}

func (r ServicebusNamespaceCustomerManagedKeyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace_customer_managed_key" "test" {
  namespace_id                      = azurerm_servicebus_namespace.test.id
  key_vault_key_id                  = azurerm_key_vault_key.test.id
  infrastructure_encryption_enabled = false
}
`, r.template(data))
}

func (r ServicebusNamespaceCustomerManagedKeyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_servicebus_namespace_customer_managed_key" "import" {
  namespace_id     = azurerm_servicebus_namespace.test.id
  key_vault_key_id = azurerm_key_vault_key.test.id
}
`, r.template(data))
}

func (r ServicebusNamespaceCustomerManagedKeyResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "test2" {
  name                        = "acctest-key-vault-2%s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Decrypt",
      "Encrypt",
      "Delete",
      "Get",
      "List",
      "Purge",
      "UnwrapKey",
      "WrapKey",
      "Verify",
      "GetRotationPolicy"
    ]
    secret_permissions = [
      "Set",
    ]
  }

  access_policy {
    tenant_id = azurerm_servicebus_namespace.test.identity[0].tenant_id
    object_id = azurerm_servicebus_namespace.test.identity[0].principal_id

    key_permissions = [
      "Create",
      "Decrypt",
      "Encrypt",
      "Delete",
      "Get",
      "List",
      "Purge",
      "UnwrapKey",
      "WrapKey",
      "Verify",
      "GetRotationPolicy"
    ]
    secret_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_key_vault_key" "test2" {
  name         = "acctestkey2-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_servicebus_namespace_customer_managed_key" "test" {
  namespace_id                      = azurerm_servicebus_namespace.test.id
  key_vault_key_id                  = azurerm_key_vault_key.test2.id
  infrastructure_encryption_enabled = false
}
`, r.template(data), data.RandomString, data.RandomString)
}
