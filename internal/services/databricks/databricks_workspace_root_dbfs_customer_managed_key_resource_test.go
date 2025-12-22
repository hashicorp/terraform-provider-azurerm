// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2024-05-01/workspaces"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DatabricksWorkspaceRootDbfsCustomerManagedKeyResource struct{}

func TestAccDatabricksWorkspaceRootDbfsCustomerManagedKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_root_dbfs_customer_managed_key", "test")
	r := DatabricksWorkspaceRootDbfsCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicVersionless(data),
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

func TestAccDatabricksWorkspaceRootDbfsCustomerManagedKey_remove(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_root_dbfs_customer_managed_key", "test")
	parent := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceRootDbfsCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.remove(data),
			Check: acceptance.ComposeTestCheckFunc(
				// Then ensure the encryption settings on the Databricks Workspace
				// have been reverted to their default state
				check.That(parent.ResourceName).DoesNotExistInAzure(r),
			),
		},
		parent.ImportStep(),
	})
}

func TestAccDatabricksWorkspaceRootDbfsCustomerManagedKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_root_dbfs_customer_managed_key", "test")
	r := DatabricksWorkspaceRootDbfsCustomerManagedKeyResource{}

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

func TestAccDatabricksWorkspaceRootDbfsCustomerManagedKey_basicManagedHSM(t *testing.T) {
	if os.Getenv("ARM_TEST_HSM_KEY") == "" {
		t.Skip("skipping as ARM_TEST_HSM_KEY is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_root_dbfs_customer_managed_key", "test")
	r := DatabricksWorkspaceRootDbfsCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicManagedHSM(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicManagedHSMVersionless(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicManagedHSM(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksWorkspaceRootDbfsCustomerManagedKey_basicAltSubscription(t *testing.T) {
	altSubscription := altSubscriptionCheck()
	if altSubscription == nil {
		t.Skip("Skipping: Test requires `ARM_SUBSCRIPTION_ID_ALT` and `ARM_TENANT_ID` environment variables to be specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_root_dbfs_customer_managed_key", "test")
	r := DatabricksWorkspaceRootDbfsCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAltSubscription(data, altSubscription),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// `key_vault_id` is always set into state based on config, so it will be missing during imports
		data.ImportStep("key_vault_id"),
	})
}

func (DatabricksWorkspaceRootDbfsCustomerManagedKeyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := workspaces.ParseWorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataBricks.WorkspacesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	// This is the only way we can tell if the CMK has actually been provisioned or not...
	if resp.Model != nil && resp.Model.Properties.Parameters != nil && resp.Model.Properties.Parameters.Encryption != nil && resp.Model.Properties.Parameters.Encryption.Value != nil && resp.Model.Properties.Parameters.Encryption.Value.KeySource != nil {
		if *resp.Model.Properties.Parameters.Encryption.Value.KeySource == workspaces.KeySourceMicrosoftPointKeyvault {
			return pointer.To(true), nil
		}
	}

	return pointer.To(false), nil
}

func (r DatabricksWorkspaceRootDbfsCustomerManagedKeyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_databricks_workspace_root_dbfs_customer_managed_key" "test" {
  workspace_id     = azurerm_databricks_workspace.test.id
  key_vault_key_id = azurerm_key_vault_key.test.id

  depends_on = [azurerm_key_vault_access_policy.databricks]
}
`, r.keyVaultTemplate(data))
}

func (r DatabricksWorkspaceRootDbfsCustomerManagedKeyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_databricks_workspace_root_dbfs_customer_managed_key" "import" {
  workspace_id     = azurerm_databricks_workspace.test.id
  key_vault_key_id = azurerm_key_vault_key.test.id
}
`, r.basic(data))
}

func (r DatabricksWorkspaceRootDbfsCustomerManagedKeyResource) basicVersionless(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_databricks_workspace_root_dbfs_customer_managed_key" "test" {
  workspace_id     = azurerm_databricks_workspace.test.id
  key_vault_key_id = azurerm_key_vault_key.test.versionless_id

  depends_on = [azurerm_key_vault_access_policy.databricks]
}
`, r.keyVaultTemplate(data))
}

func (r DatabricksWorkspaceRootDbfsCustomerManagedKeyResource) remove(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s
`, r.keyVaultTemplate(data))
}

func (r DatabricksWorkspaceRootDbfsCustomerManagedKeyResource) basicManagedHSM(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_databricks_workspace_root_dbfs_customer_managed_key" "test" {
  workspace_id     = azurerm_databricks_workspace.test.id
  key_vault_key_id = azurerm_key_vault_managed_hardware_security_module_key.test.versioned_id

  depends_on = [azurerm_key_vault_managed_hardware_security_module_role_assignment.user]
}
`, r.managedHSMVaultTemplate(data))
}

func (r DatabricksWorkspaceRootDbfsCustomerManagedKeyResource) basicManagedHSMVersionless(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_databricks_workspace_root_dbfs_customer_managed_key" "test" {
  workspace_id     = azurerm_databricks_workspace.test.id
  key_vault_key_id = azurerm_key_vault_managed_hardware_security_module_key.test.id

  depends_on = [azurerm_key_vault_managed_hardware_security_module_role_assignment.user]
}
`, r.managedHSMVaultTemplate(data))
}

func (r DatabricksWorkspaceRootDbfsCustomerManagedKeyResource) basicAltSubscription(data acceptance.TestData, alternate *DatabricksWorkspaceAlternateSubscription) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_databricks_workspace_root_dbfs_customer_managed_key" "test" {
  workspace_id     = azurerm_databricks_workspace.test.id
  key_vault_key_id = azurerm_key_vault_key.test.id
  key_vault_id     = azurerm_key_vault.test.id

  depends_on = [azurerm_key_vault_access_policy.databricks]
}
`, r.keyVaultAltSubscriptionTemplate(data, alternate), data.RandomInteger)
}

func (r DatabricksWorkspaceRootDbfsCustomerManagedKeyResource) keyVaultTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault" "test" {
  name                = "acctest-kv-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "test" {
  depends_on = [azurerm_key_vault_access_policy.terraform]

  name         = "acctest-key-%[3]d"
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

resource "azurerm_key_vault_access_policy" "terraform" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault.test.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "List",
    "Create",
    "Decrypt",
    "Encrypt",
    "GetRotationPolicy",
    "Sign",
    "UnwrapKey",
    "Verify",
    "WrapKey",
    "Delete",
    "Restore",
    "Recover",
    "Update",
    "Purge",
  ]
}

resource "azurerm_key_vault_access_policy" "databricks" {
  depends_on = [azurerm_databricks_workspace.test]

  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_databricks_workspace.test.storage_account_identity.0.tenant_id
  object_id    = azurerm_databricks_workspace.test.storage_account_identity.0.principal_id

  key_permissions = [
    "Get",
    "GetRotationPolicy",
    "UnwrapKey",
    "WrapKey",
    "Delete",
  ]
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r DatabricksWorkspaceRootDbfsCustomerManagedKeyResource) keyVaultAltSubscriptionTemplate(data acceptance.TestData, alternate *DatabricksWorkspaceAlternateSubscription) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm-alt" {
  features {}

  tenant_id       = "%[2]s"
  subscription_id = "%[3]s"
}

resource "azurerm_resource_group" "keyVault" {
  provider = azurerm-alt

  name     = "acctestRG-databricks-dbfs-alt-sub-%[4]d"
  location = "%[5]s"
}

resource "azurerm_key_vault" "test" {
  provider = azurerm-alt

  name                = "acctestkv-alt-%[6]s"
  location            = azurerm_resource_group.keyVault.location
  resource_group_name = azurerm_resource_group.keyVault.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "test" {
  depends_on = [azurerm_key_vault_access_policy.terraform]
  provider   = azurerm-alt

  name         = "acctest-key-%[4]d"
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

resource "azurerm_key_vault_access_policy" "terraform" {
  provider = azurerm-alt

  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault.test.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "List",
    "Create",
    "Decrypt",
    "Encrypt",
    "GetRotationPolicy",
    "Sign",
    "UnwrapKey",
    "Verify",
    "WrapKey",
    "Delete",
    "Restore",
    "Recover",
    "Update",
    "Purge",
  ]
}

resource "azurerm_key_vault_access_policy" "databricks" {
  depends_on = [azurerm_databricks_workspace.test]
  provider   = azurerm-alt

  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_databricks_workspace.test.storage_account_identity.0.tenant_id
  object_id    = azurerm_databricks_workspace.test.storage_account_identity.0.principal_id

  key_permissions = [
    "Get",
    "GetRotationPolicy",
    "UnwrapKey",
    "WrapKey",
    "Delete",
  ]
}
`, r.template(data), alternate.tenantID, alternate.subscriptionID, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r DatabricksWorkspaceRootDbfsCustomerManagedKeyResource) managedHSMVaultTemplate(data acceptance.TestData) string {
	roleAssignmentName1, _ := uuid.GenerateUUID()
	roleAssignmentName2, _ := uuid.GenerateUUID()
	roleAssignmentName3, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault" "test" {
  name                       = "acctest%[2]s"
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
  name                                      = "acctestkvHsm%[2]s"
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
  name               = "%[3]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.crypto-officer.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "test1" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[4]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.crypto-user.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id

  depends_on = [azurerm_key_vault_managed_hardware_security_module_role_assignment.test]
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "user" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[5]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.encrypt-user.resource_manager_id
  principal_id       = azurerm_databricks_workspace.test.storage_account_identity.0.principal_id

  depends_on = [azurerm_key_vault_managed_hardware_security_module_role_assignment.test1]
}

resource "azurerm_key_vault_managed_hardware_security_module_key" "test" {
  name           = "acctestHSMK-%[2]s"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "RSA-HSM"
  key_size       = 2048
  key_opts       = ["unwrapKey", "wrapKey"]

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.test,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.test1
  ]
}
`, r.template(data), data.RandomString, roleAssignmentName1, roleAssignmentName2, roleAssignmentName3)
}

func (DatabricksWorkspaceRootDbfsCustomerManagedKeyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-db-%[1]d"
  location = "%[2]s"
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestDBW-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "premium"

  customer_managed_key_enabled      = true
  infrastructure_encryption_enabled = true

  custom_parameters {
    no_public_ip = true
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
