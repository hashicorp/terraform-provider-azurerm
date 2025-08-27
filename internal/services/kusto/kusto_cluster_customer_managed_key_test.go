// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-uuid"
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

func TestAccKustoClusterCustomerManagedKey_managedHsmKey(t *testing.T) {
	// Skipping this test by default, as the managed HSM is costly.
	if os.Getenv("ARM_TEST_HSM_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_HSM_KEY is not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster_customer_managed_key", "test")
	r := KustoClusterCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedHsmKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_hsm_key_id").Exists(),
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

func (KustoClusterCustomerManagedKeyResource) managedHsmKey(data acceptance.TestData) string {
	template := KustoClusterCustomerManagedKeyResource{}.hsmKeyTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_cluster_customer_managed_key" "test" {
  cluster_id         = azurerm_kusto_cluster.test.id
  managed_hsm_key_id = azurerm_key_vault_managed_hardware_security_module_key.test.id
}
`, template)
}

func (KustoClusterCustomerManagedKeyResource) hsmKeyTemplate(data acceptance.TestData) string {
	raName1, _ := uuid.GenerateUUID()
	raName2, _ := uuid.GenerateUUID()
	raName3, _ := uuid.GenerateUUID()
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
  name     = "acctestRG-%[3]d"
  location = "%[2]s"
}

resource "azurerm_key_vault" "test" {
  name                       = "acc%[3]d"
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
  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_certificate" "cert" {
  count        = 3
  name         = "acchsmcert${count.index}"
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

// Purge Protection must be enabled to configure Managed HSM Key
// https://learn.microsoft.com/en-us/azure/data-explorer/security#store-customer-managed-keys-in-azure-key-vault
resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                       = "kvHsm%[3]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  sku_name                   = "Standard_B1"
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  admin_object_ids           = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled   = true
  soft_delete_retention_days = 7

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

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "client1" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[4]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.crypto-officer.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "client2" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[5]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.crypto-user.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id

  // Avoid role assignments conflicts.
  depends_on = [azurerm_key_vault_managed_hardware_security_module_role_assignment.client1]
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "kusto" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[6]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.encrypt-user.resource_manager_id
  principal_id       = azurerm_kusto_cluster.test.identity[0].principal_id

  // Avoid role assignments conflicts.
  depends_on = [azurerm_key_vault_managed_hardware_security_module_role_assignment.client2]
}

resource "azurerm_key_vault_managed_hardware_security_module_key" "test" {
  name           = "acctestHSMK-%[1]s"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "RSA-HSM"
  key_size       = 2048
  key_opts       = ["unwrapKey", "wrapKey"]

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.client1,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.client2,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.kusto,
  ]
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%[1]s"
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

`, data.RandomString, data.Locations.Primary, data.RandomInteger, raName1, raName2, raName3)
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
