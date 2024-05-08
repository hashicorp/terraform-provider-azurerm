// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KustoClusterCustomerManagedKeyResource struct{}

func TestAccKustoClusterCustomerManagedKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster_customer_managed_key", "test")
	r := KustoClusterCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_vault_id").Exists(),
				check.That(data.ResourceName).Key("key_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			// Delete the encryption settings resource and verify it is gone
			Config: r.template(data),
			Check: acceptance.ComposeTestCheckFunc(
				// Then ensure the encryption settings on the Kusto cluster
				// have been reverted to their default state
				check.That("azurerm_kusto_cluster.test").DoesNotExistInAzure(r),
			),
		},
	})
}

func TestAccKustoClusterCustomerManagedKey_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster_customer_managed_key", "test")
	r := KustoClusterCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_vault_id").Exists(),
				check.That(data.ResourceName).Key("key_name").Exists(),
				check.That(data.ResourceName).Key("key_version").Exists(),
			),
		},
		data.ImportStep(),
		{
			// Delete the encryption settings resource and verify it is gone
			Config: r.template(data),
			Check: acceptance.ComposeTestCheckFunc(
				// Then ensure the encryption settings on the Kusto cluster
				// have been reverted to their default state
				check.That("azurerm_kusto_cluster.test").DoesNotExistInAzure(r),
			),
		},
	})
}

func TestAccKustoClusterCustomerManagedKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster_customer_managed_key", "test")
	r := KustoClusterCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_vault_id").Exists(),
				check.That(data.ResourceName).Key("key_name").Exists(),
				check.That(data.ResourceName).Key("key_version").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccKustoClusterCustomerManagedKey_updateKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster_customer_managed_key", "test")
	r := KustoClusterCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_vault_id").Exists(),
				check.That(data.ResourceName).Key("key_name").Exists(),
				check.That(data.ResourceName).Key("key_version").Exists(),
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

func TestAccKustoClusterCustomerManagedKey_userIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster_customer_managed_key", "test")
	r := KustoClusterCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (KustoClusterCustomerManagedKeyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseKustoClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Kusto.ClustersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	if resp.Model == nil {
		return nil, fmt.Errorf("response model is empty")
	}

	if resp.Model.Properties == nil || resp.Model.Properties.KeyVaultProperties == nil {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (KustoClusterCustomerManagedKeyResource) basic(data acceptance.TestData) string {
	template := KustoClusterCustomerManagedKeyResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_cluster_customer_managed_key" "test" {
  cluster_id   = azurerm_kusto_cluster.test.id
  key_vault_id = azurerm_key_vault.test.id
  key_name     = azurerm_key_vault_key.first.name
}
`, template)
}

func (KustoClusterCustomerManagedKeyResource) complete(data acceptance.TestData) string {
	template := KustoClusterCustomerManagedKeyResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_cluster_customer_managed_key" "test" {
  cluster_id   = azurerm_kusto_cluster.test.id
  key_vault_id = azurerm_key_vault.test.id
  key_name     = azurerm_key_vault_key.first.name
  key_version  = azurerm_key_vault_key.first.version
}
`, template)
}

func (KustoClusterCustomerManagedKeyResource) requiresImport(data acceptance.TestData) string {
	template := KustoClusterCustomerManagedKeyResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_cluster_customer_managed_key" "import" {
  cluster_id   = azurerm_kusto_cluster_customer_managed_key.test.cluster_id
  key_vault_id = azurerm_kusto_cluster_customer_managed_key.test.key_vault_id
  key_name     = azurerm_kusto_cluster_customer_managed_key.test.key_name
}
`, template)
}

func (KustoClusterCustomerManagedKeyResource) updated(data acceptance.TestData) string {
	template := KustoClusterCustomerManagedKeyResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "second" {
  name         = "second"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.cluster,
  ]
}

resource "azurerm_kusto_cluster_customer_managed_key" "test" {
  cluster_id   = azurerm_kusto_cluster.test.id
  key_vault_id = azurerm_key_vault.test.id
  key_name     = azurerm_key_vault_key.second.name
  key_version  = azurerm_key_vault_key.second.version
}
`, template)
}

func (KustoClusterCustomerManagedKeyResource) userIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "cluster" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions = ["Get", "UnwrapKey", "WrapKey", "GetRotationPolicy"]
}

resource "azurerm_key_vault_access_policy" "client" {
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
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "test"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.cluster,
  ]
}

resource "azurerm_kusto_cluster_customer_managed_key" "test" {
  cluster_id    = azurerm_kusto_cluster.test.id
  key_vault_id  = azurerm_key_vault.test.id
  key_name      = azurerm_key_vault_key.test.name
  key_version   = azurerm_key_vault_key.test.version
  user_identity = azurerm_user_assigned_identity.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}

func (KustoClusterCustomerManagedKeyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "cluster" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_kusto_cluster.test.identity.0.tenant_id
  object_id    = azurerm_kusto_cluster.test.identity.0.principal_id

  key_permissions = ["Get", "UnwrapKey", "WrapKey"]
}

resource "azurerm_key_vault_access_policy" "client" {
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
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}

resource "azurerm_key_vault_key" "first" {
  name         = "test"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.cluster,
  ]
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
