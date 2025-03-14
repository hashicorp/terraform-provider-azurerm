// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogAnalyticsClusterCustomerManagedKeyResource struct{}

func TestAccLogAnalyticsClusterCustomerManagedKey_keyVaultKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster_customer_managed_key", "test")
	r := LogAnalyticsClusterCustomerManagedKeyResource{}

	if os.Getenv("ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS is not specified")
		return
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultKey(data, "azurerm_key_vault_key.test1.id"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultKey(data, "azurerm_key_vault_key.test1.versionless_id"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultKey(data, "azurerm_key_vault_key.test2.id"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsClusterCustomerManagedKey_managedHsmKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster_customer_managed_key", "test")
	r := LogAnalyticsClusterCustomerManagedKeyResource{}

	if os.Getenv("ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS is not specified")
		return
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedHsmKey(data, "azurerm_key_vault_managed_hardware_security_module_key.test1.versioned_id"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.managedHsmKey(data, "azurerm_key_vault_managed_hardware_security_module_key.test1.id"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.managedHsmKey(data, "azurerm_key_vault_managed_hardware_security_module_key.test2.versioned_id"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t LogAnalyticsClusterCustomerManagedKeyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := clusters.ParseClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LogAnalytics.ClusterClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	exists := false
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if kv := props.KeyVaultProperties; kv != nil {
				exists = kv.KeyVaultUri != nil && kv.KeyVersion != nil && kv.KeyName != nil
			}
		}
	}

	return pointer.To(exists), nil
}

func (r LogAnalyticsClusterCustomerManagedKeyResource) keyVaultKey(data acceptance.TestData, keyId string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster_customer_managed_key" "test" {
  log_analytics_cluster_id = azurerm_log_analytics_cluster.test.id
  key_vault_key_id         = %s

  depends_on = [azurerm_key_vault_access_policy.system_identity]
}
`, r.keyVaultKeyTemplate(data), keyId)
}

func (r LogAnalyticsClusterCustomerManagedKeyResource) managedHsmKey(data acceptance.TestData, keyId string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster_customer_managed_key" "test" {
  log_analytics_cluster_id = azurerm_log_analytics_cluster.test.id
  managed_hsm_key_id       = %s

  depends_on = [azurerm_key_vault_access_policy.system_identity]
}
`, r.managedHsmKeyTemplate(data), keyId)
}

func (r LogAnalyticsClusterCustomerManagedKeyResource) keyVaultKeyTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "test1" {
  name         = "key1-%[2]s"
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

  depends_on = [azurerm_key_vault_access_policy.principal]
}

resource "azurerm_key_vault_key" "test2" {
  name         = "key2-%[2]s"
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

  depends_on = [azurerm_key_vault_access_policy.principal]
}
`, r.template(data), data.RandomString)
}

func (r LogAnalyticsClusterCustomerManagedKeyResource) managedHsmKeyTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_certificate" "test" {
  count        = 3
  name         = "acctest%[2]d${count.index}"
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

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
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

  depends_on = [azurerm_key_vault_access_policy.principal]
}

resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                     = "acctest-%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  sku_name                 = "Standard_B1"
  tenant_id                = data.azurerm_client_config.current.tenant_id
  admin_object_ids         = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled = false

  security_domain_key_vault_certificate_ids = [for cert in azurerm_key_vault_certificate.test : cert.id]
  security_domain_quorum                    = 3

  depends_on = [azurerm_key_vault_access_policy.principal]
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "crypto_user_system" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[4]s"
  scope              = "/"
  role_definition_id = "/Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/21dbd100-6940-42c2-9190-5d6cb909625b"
  principal_id       = azurerm_log_analytics_cluster.test.identity[0].principal_id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "crypto_user_principal" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[5]s"
  scope              = "/"
  role_definition_id = "/Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/21dbd100-6940-42c2-9190-5d6cb909625b"
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "crypto_officer_principal" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[6]s"
  scope              = "/"
  role_definition_id = "/Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/515eb02d-2335-4d2d-92f2-b1cbdf9c3778"
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_key" "test1" {
  name           = "acctest1%[2]d"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "RSA-HSM"
  key_size       = 2048
  key_opts       = ["sign", "unwrapKey", "wrapKey"]

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.crypto_user_principal,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.crypto_officer_principal
  ]
}

resource "azurerm_key_vault_managed_hardware_security_module_key" "test2" {
  name           = "acctest2%[2]d"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "RSA-HSM"
  key_size       = 2048
  key_opts       = ["sign", "unwrapKey", "wrapKey"]

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.crypto_user_principal,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.crypto_officer_principal
  ]
}
`, r.template(data), data.RandomInteger, data.RandomString, uuid.New(), uuid.New(), uuid.New())
}

func (LogAnalyticsClusterCustomerManagedKeyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-la-%[1]d"
  location = "%[2]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_log_analytics_cluster" "test" {
  name                = "acctest-LA-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault" "test" {
  name                = "acctestvault%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  soft_delete_retention_days = 7
  purge_protection_enabled   = false
}


resource "azurerm_key_vault_access_policy" "principal" {
  key_vault_id = azurerm_key_vault.test.id

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "List",
    "Purge",
    "Update",
    "GetRotationPolicy",
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

resource "azurerm_key_vault_access_policy" "system_identity" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "Get",
    "UnwrapKey",
    "WrapKey",
    "GetRotationPolicy",
  ]

  tenant_id = azurerm_log_analytics_cluster.test.identity.0.tenant_id
  object_id = azurerm_log_analytics_cluster.test.identity.0.principal_id

  depends_on = [azurerm_key_vault_access_policy.principal]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
