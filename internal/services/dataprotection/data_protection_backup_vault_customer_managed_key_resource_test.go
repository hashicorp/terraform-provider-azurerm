// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupvaults"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupVaultCustomerManagedKeyResource struct{}

func TestAccDataProtectionBackupVaultCustomerManagedKey_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault_customer_managed_key", "test")
	r := DataProtectionBackupVaultCustomerManagedKeyResource{}
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

func TestAccDataProtectionBackupVaultCustomerManagedKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault_customer_managed_key", "test")
	r := DataProtectionBackupVaultCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDataProtectionBackupVaultCustomerManagedKey_updated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault_customer_managed_key", "test")
	r := DataProtectionBackupVaultCustomerManagedKeyResource{}
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

func TestAccDataProtectionBackupVaultCustomerManagedKey_managedHSM(t *testing.T) {
	if os.Getenv("ARM_TEST_HSM_KEY") == "" {
		t.Skip("skipping as ARM_TEST_HSM_KEY is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault_customer_managed_key", "test")
	r := DataProtectionBackupVaultCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedHSM(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.managedHSMVersionless(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.managedHSM(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupVaultCustomerManagedKey_conflictedEncryptionSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault_customer_managed_key", "test")
	r := DataProtectionBackupVaultCustomerManagedKeyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.conflictedEncryptionSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectError: regexp.MustCompile("customer managed keys settings have been specified in `user_assigned_identity_encryption_settings` block of `azurerm_data_protection_backup_vault` resource. `azurerm_data_protection_backup_vault_customer_managed_key` resource is not required and should be removed"),
		},
	})
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backupvaults.ParseBackupVaultID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupVaultClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil || resp.Model.Properties.SecuritySettings == nil || resp.Model.Properties.SecuritySettings.EncryptionSettings == nil {
		return pointer.To(false), nil
	}
	return pointer.To(true), nil
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dataprotection-%d"
  location = "%s"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"

  identity {
    type = "SystemAssigned"
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
    tenant_id = azurerm_data_protection_backup_vault.test.identity[0].tenant_id
    object_id = azurerm_data_protection_backup_vault.test.identity[0].principal_id

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

func (r DataProtectionBackupVaultCustomerManagedKeyResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_vault_customer_managed_key" "test" {
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  key_vault_key_id                = azurerm_key_vault_key.test.id
}
`, template)
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) requiresImport(data acceptance.TestData) string {
	template := r.complete(data)
	return fmt.Sprintf(`
%s
resource "azurerm_data_protection_backup_vault_customer_managed_key" "import" {
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault_customer_managed_key.test.data_protection_backup_vault_id
  key_vault_key_id                = azurerm_data_protection_backup_vault_customer_managed_key.test.key_vault_key_id
}
`, template)
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) updated(data acceptance.TestData) string {
	template := r.template(data)
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
    tenant_id = azurerm_data_protection_backup_vault.test.identity[0].tenant_id
    object_id = azurerm_data_protection_backup_vault.test.identity[0].principal_id

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

resource "azurerm_data_protection_backup_vault_customer_managed_key" "test" {
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  key_vault_key_id                = azurerm_key_vault_key.test2.id
}
`, template, data.RandomString, data.RandomString)
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) managedHSM(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

%[1]s

resource "azurerm_data_protection_backup_vault_customer_managed_key" "test" {
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  key_vault_key_id                = azurerm_key_vault_managed_hardware_security_module_key.test.versioned_id

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.backup-vault,
  ]
}
`, r.templateManagedHSM(data))
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) managedHSMVersionless(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

%[1]s

resource "azurerm_data_protection_backup_vault_customer_managed_key" "test" {
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  key_vault_key_id                = azurerm_key_vault_managed_hardware_security_module_key.test.id

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.backup-vault,
  ]
}
`, r.templateManagedHSM(data))
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) templateManagedHSM(data acceptance.TestData) string {
	uuid1, _ := uuid.GenerateUUID()
	uuid2, _ := uuid.GenerateUUID()
	uuid3, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%[1]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Recover",
      "Update",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Delete",
      "Get",
      "Set",
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "DeleteIssuers",
      "Get",
      "Purge",
      "Update"
    ]
  }
}

resource "azurerm_key_vault_certificate" "cert" {
  count        = 3
  name         = "acctesthsmcert${count.index}"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      extended_key_usage = []

      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }
}

resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                                      = "acctestkvHsm-%[1]s"
  resource_group_name                       = azurerm_resource_group.test.name
  location                                  = azurerm_resource_group.test.location
  sku_name                                  = "Standard_B1"
  tenant_id                                 = data.azurerm_client_config.current.tenant_id
  admin_object_ids                          = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled                  = true
  soft_delete_retention_days                = 7
  security_domain_key_vault_certificate_ids = [for cert in azurerm_key_vault_certificate.cert : cert.id]
  security_domain_quorum                    = 3
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "crypto-officer" {
  name           = "515eb02d-2335-4d2d-92f2-b1cbdf9c3778"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "crypto-user" {
  name           = "21dbd100-6940-42c2-9190-5d6cb909625b"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "encrypt-user" {
  name           = "33413926-3206-4cdd-b39a-83574fe37a17"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "test" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[2]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.crypto-officer.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "test1" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[3]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.crypto-user.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id

  depends_on = [azurerm_key_vault_managed_hardware_security_module_role_assignment.test]
}

resource "azurerm_key_vault_managed_hardware_security_module_key" "test" {
  name           = "acctestHSMK-%[1]s"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "RSA-HSM"
  key_size       = 2048
  key_opts       = ["unwrapKey", "wrapKey"]

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.test,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.test1
  ]
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[4]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "backup-vault" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[5]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.encrypt-user.resource_manager_id
  principal_id       = azurerm_data_protection_backup_vault.test.identity.0.principal_id

  depends_on = [azurerm_key_vault_managed_hardware_security_module_role_assignment.test1]
}
`, data.RandomString, uuid1, uuid2, data.RandomInteger, uuid3)
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) conflictedEncryptionSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dataprotection-%d"
  location = "%s"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  user_assigned_identity_encryption_settings {
    identity_id      = azurerm_user_assigned_identity.test.id
    key_vault_key_id = azurerm_key_vault_key.test.id
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
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

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

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctestBV-%d"
}

resource "azurerm_data_protection_backup_vault_customer_managed_key" "test" {
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  key_vault_key_id                = azurerm_key_vault_key.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomString, data.RandomInteger)
}
