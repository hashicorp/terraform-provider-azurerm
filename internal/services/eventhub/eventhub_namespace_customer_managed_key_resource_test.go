// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventhub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type EventHubNamespaceCustomerManagedKeyResource struct{}

func TestAccEventHubNamespaceCustomerManagedKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_customer_managed_key", "test")
	r := EventHubNamespaceCustomerManagedKeyResource{}

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

func TestAccEventHubNamespaceCustomerManagedKey_withUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_customer_managed_key", "test")
	r := EventHubNamespaceCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespaceCustomerManagedKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_customer_managed_key", "test")
	r := EventHubNamespaceCustomerManagedKeyResource{}

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

func TestAccEventHubNamespaceCustomerManagedKey_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_customer_managed_key", "test")
	r := EventHubNamespaceCustomerManagedKeyResource{}

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

func TestAccEventHubNamespaceCustomerManagedKey_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_customer_managed_key", "test")
	r := EventHubNamespaceCustomerManagedKeyResource{}

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

func (r EventHubNamespaceCustomerManagedKeyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := namespaces.ParseNamespaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Eventhub.NamespacesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}
	if resp.Model == nil {
		return nil, fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	if resp.Model.Properties == nil || resp.Model.Properties.Encryption == nil {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (r EventHubNamespaceCustomerManagedKeyResource) requiresImport(data acceptance.TestData) string {
	template := EventHubNamespaceCustomerManagedKeyResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_namespace_customer_managed_key" "import" {
  eventhub_namespace_id = azurerm_eventhub_namespace_customer_managed_key.test.eventhub_namespace_id
  key_vault_key_ids     = azurerm_eventhub_namespace_customer_managed_key.test.key_vault_key_ids
}
`, template)
}

func (r EventHubNamespaceCustomerManagedKeyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_namespace_customer_managed_key" "test" {
  eventhub_namespace_id = azurerm_eventhub_namespace.test.id
  key_vault_key_ids     = [azurerm_key_vault_key.test.id]
}
`, r.template(data))
}

func (r EventHubNamespaceCustomerManagedKeyResource) withUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_namespace_customer_managed_key" "test" {
  eventhub_namespace_id     = azurerm_eventhub_namespace.test.id
  key_vault_key_ids         = [azurerm_key_vault_key.test.id]
  user_assigned_identity_id = azurerm_user_assigned_identity.test.id
}
`, r.templateWithUserAssignedIdentity(data))
}

func (r EventHubNamespaceCustomerManagedKeyResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "test2" {
  name         = "acctestkvkey2%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
  ]
}

resource "azurerm_eventhub_namespace_customer_managed_key" "test" {
  eventhub_namespace_id = azurerm_eventhub_namespace.test.id
  key_vault_key_ids     = [azurerm_key_vault_key.test2.id]
}
`, r.template(data), data.RandomString)
}

func (r EventHubNamespaceCustomerManagedKeyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "test2" {
  name         = "acctestkvkey2%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
  ]
}

resource "azurerm_eventhub_namespace_customer_managed_key" "test" {
  eventhub_namespace_id             = azurerm_eventhub_namespace.test.id
  key_vault_key_ids                 = [azurerm_key_vault_key.test.id, azurerm_key_vault_key.test2.id]
  infrastructure_encryption_enabled = true
}
`, r.template(data), data.RandomString)
}

func (r EventHubNamespaceCustomerManagedKeyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-namespacecmk-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-namespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Premium"

  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%[3]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_eventhub_namespace.test.identity.0.tenant_id
  object_id    = azurerm_eventhub_namespace.test.identity.0.principal_id

  key_permissions = ["Get", "UnwrapKey", "WrapKey", "GetRotationPolicy"]
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "List",
    "Purge",
    "Recover",
    "GetRotationPolicy"
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey%[3]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r EventHubNamespaceCustomerManagedKeyResource) templateWithUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-namespacecmk-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctest-identity-%[3]s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-namespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Premium"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%[3]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_user_assigned_identity.test.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions = ["Get", "UnwrapKey", "WrapKey", "GetRotationPolicy"]
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "List",
    "Purge",
    "Recover",
    "GetRotationPolicy"
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey%[3]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
