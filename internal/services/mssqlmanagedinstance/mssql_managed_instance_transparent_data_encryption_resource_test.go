// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlManagedInstanceTransparentDataEncryptionResource struct{}

func TestAccMsSqlManagedInstanceTransparentDataEncryption_keyVaultSystemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_transparent_data_encryption", "test")
	r := MsSqlManagedInstanceTransparentDataEncryptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlManagedInstanceTransparentDataEncryption_keyVaultUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_transparent_data_encryption", "test")
	r := MsSqlManagedInstanceTransparentDataEncryptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlManagedInstanceTransparentDataEncryption_autoRotate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_transparent_data_encryption", "test")
	r := MsSqlManagedInstanceTransparentDataEncryptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoRotate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_vault_key_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlManagedInstanceTransparentDataEncryption_systemManaged(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_transparent_data_encryption", "test")
	r := MsSqlManagedInstanceTransparentDataEncryptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemManaged(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_vault_key_id").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlManagedInstanceTransparentDataEncryption_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_transparent_data_encryption", "test")
	r := MsSqlManagedInstanceTransparentDataEncryptionResource{}

	// Test going from systemManaged to keyVault and back
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.systemManaged(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_vault_key_id").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}
func (MsSqlManagedInstanceTransparentDataEncryptionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ManagedInstanceEncryptionProtectorID(state.ID)
	if err != nil {
		return nil, err
	}

	instanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)
	resp, err := client.MSSQLManagedInstance.ManagedInstanceEncryptionProtectorClient.Get(ctx, instanceId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("encryption protector for managed instance %q (Resource Group %q) does not exist", id.ManagedInstanceName, id.ResourceGroup)
		}

		return nil, fmt.Errorf("reading Encryption Protector for managed instance %q (Resource Group %q): %v", id.ManagedInstanceName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r MsSqlManagedInstanceTransparentDataEncryptionResource) keyVaultSystemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "test" {
  name                        = "acctestsqlserver%[2]s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = false

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get", "List", "Create", "Delete", "Update", "Purge", "GetRotationPolicy", "SetRotationPolicy"
    ]
  }

  access_policy {
    tenant_id = azurerm_mssql_managed_instance.test.identity[0].tenant_id
    object_id = azurerm_mssql_managed_instance.test.identity[0].principal_id

    key_permissions = [
      "Get", "WrapKey", "UnwrapKey", "List", "Create", "GetRotationPolicy", "SetRotationPolicy"
    ]
  }
}

resource "azurerm_key_vault_key" "generated" {
  name         = "keyVault"
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

  depends_on = [
    azurerm_key_vault.test,
  ]
}

resource "azurerm_mssql_managed_instance_transparent_data_encryption" "test" {
  managed_instance_id = azurerm_mssql_managed_instance.test.id
  key_vault_key_id    = azurerm_key_vault_key.generated.id
}
`, r.serverSAMI(data), data.RandomStringOfLength(5))
}

func (r MsSqlManagedInstanceTransparentDataEncryptionResource) keyVaultUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "test" {
  name                        = "acctestsqlserver%[2]s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = false

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get", "List", "Create", "Delete", "Update", "Purge", "GetRotationPolicy", "SetRotationPolicy"
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    key_permissions = [
      "Get", "WrapKey", "UnwrapKey", "List", "Create", "GetRotationPolicy", "SetRotationPolicy"
    ]
  }
}

resource "azurerm_key_vault_key" "generated" {
  name         = "keyVault"
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

  depends_on = [
    azurerm_key_vault.test,
  ]
}

resource "azurerm_mssql_managed_instance_transparent_data_encryption" "test" {
  managed_instance_id = azurerm_mssql_managed_instance.test.id
  key_vault_key_id    = azurerm_key_vault_key.generated.id
}
`, r.serverUAMI(data), data.RandomStringOfLength(5))
}

func (r MsSqlManagedInstanceTransparentDataEncryptionResource) autoRotate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "test" {
  name                        = "acctestsqlserver%[2]s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = false

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get", "List", "Create", "Delete", "Update", "Purge", "GetRotationPolicy", "SetRotationPolicy"
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    key_permissions = [
      "Get", "WrapKey", "UnwrapKey", "List", "Create", "GetRotationPolicy", "SetRotationPolicy"
    ]
  }
}

resource "azurerm_key_vault_key" "generated" {
  name         = "keyVault"
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

  depends_on = [
    azurerm_key_vault.test,
  ]
}

resource "azurerm_mssql_managed_instance_transparent_data_encryption" "test" {
  managed_instance_id   = azurerm_mssql_managed_instance.test.id
  key_vault_key_id      = azurerm_key_vault_key.generated.id
  auto_rotation_enabled = true
}
`, r.serverUAMI(data), data.RandomStringOfLength(5))
}

func (r MsSqlManagedInstanceTransparentDataEncryptionResource) systemManaged(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_managed_instance_transparent_data_encryption" "test" {
  managed_instance_id = azurerm_mssql_managed_instance.test.id
}
`, r.serverSAMI(data))
}

func (MsSqlManagedInstanceTransparentDataEncryptionResource) serverSAMI(data acceptance.TestData) string {
	db := MsSqlManagedInstanceResource{}

	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }

    resource_group {
      /* Due to the creation of unmanaged Microsoft.Network/networkIntentPolicies in this service, 
      prevent_deletion_if_contains_resources has been added here to allow the test resources to be 
       deleted until this can be properly investigated
      */
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_client_config" "current" {}

%[1]s

resource "azurerm_mssql_managed_instance" "test" {
  name                = "acctestsqlserver%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.test.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "staging"
    database    = "test"
  }
}




`, db.template(data, data.Locations.Primary), data.RandomInteger)
}

func (MsSqlManagedInstanceTransparentDataEncryptionResource) serverUAMI(data acceptance.TestData) string {
	db := MsSqlManagedInstanceResource{}

	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
    resource_group {
      /* Due to the creation of unmanaged Microsoft.Network/networkIntentPolicies in this service,
      prevent_deletion_if_contains_resources has been added here to allow the test resources to be
       deleted until this can be properly investigated
      */
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctestuami%[2]d"
}

%[1]s

resource "azurerm_mssql_managed_instance" "test" {
  name                = "acctestsqlserver%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.test.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  tags = {
    environment = "staging"
    database    = "test"
  }
}








`, db.template(data, data.Locations.Primary), data.RandomInteger)
}
