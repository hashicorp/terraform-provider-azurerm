// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DatabricksWorkspaceServerlessResource struct{}

func TestAccDatabricksWorkspaceServerless_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_serverless", "test")
	r := DatabricksWorkspaceServerlessResource{}

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

func TestAccDatabricksWorkspaceServerless_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_serverless", "test")
	r := DatabricksWorkspaceServerlessResource{}

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

func TestAccDatabricksWorkspaceServerless_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_serverless", "test")
	r := DatabricksWorkspaceServerlessResource{}

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

func TestAccDatabricksWorkspaceServerless_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_serverless", "test")
	r := DatabricksWorkspaceServerlessResource{}

	data.ResourceTestIgnoreRecreate(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksWorkspaceServerless_privateLink(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_serverless", "test")
	r := DatabricksWorkspaceServerlessResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateLink(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksWorkspaceServerless_altSubscriptionCmkServicesOnly(t *testing.T) {
	altSubscription := altSubscriptionCheck()

	if altSubscription == nil {
		t.Skip("Skipping: Test requires `ARM_SUBSCRIPTION_ID_ALT` and `ARM_TENANT_ID` environment variables to be specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_serverless", "test")
	r := DatabricksWorkspaceServerlessResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.altSubscriptionCmkServicesOnly(data, altSubscription),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("managed_services_cmk_key_vault_id"),
	})
}

func TestAccDatabricksWorkspaceServerless_enhancedSecurityComplianceWithoutAutomaticClusterUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_serverless", "test")
	r := DatabricksWorkspaceServerlessResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.enhancedSecurityCompliance(data, false, true, true),
			ExpectError: regexp.MustCompile("`automatic_cluster_update_enabled` .* must be set to `true` when `compliance_security_profile_enabled` is set to `true`"),
		},
	})
}

func TestAccDatabricksWorkspaceServerless_enhancedSecurityComplianceWithoutEnhancedSecurityMonitoring(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_serverless", "test")
	r := DatabricksWorkspaceServerlessResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.enhancedSecurityCompliance(data, true, true, false),
			ExpectError: regexp.MustCompile("`enhanced_security_monitoring_enabled` must be set to `true` when `compliance_security_profile_enabled` is set to `true`"),
		},
	})
}

func TestAccDatabricksWorkspaceServerless_enhancedSecurityComplianceStandardsWithoutProfileEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_serverless", "test")
	r := DatabricksWorkspaceServerlessResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.enhancedSecurityCompliance(data, false, false, false),
			ExpectError: regexp.MustCompile("`compliance_security_profile_standards` cannot be set when `compliance_security_profile_enabled` is `false`"),
		},
	})
}

func (DatabricksWorkspaceServerlessResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := workspaces.ParseWorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.DataBricks.WorkspacesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (DatabricksWorkspaceServerlessResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-databricks-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DatabricksWorkspaceServerlessResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_databricks_workspace_serverless" "test" {
  name                = "acctest-dbws-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomInteger)
}

func (r DatabricksWorkspaceServerlessResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_databricks_workspace_serverless" "import" {
  name                = azurerm_databricks_workspace_serverless.test.name
  resource_group_name = azurerm_databricks_workspace_serverless.test.resource_group_name
  location            = azurerm_databricks_workspace_serverless.test.location
}
`, r.basic(data))
}

func (r DatabricksWorkspaceServerlessResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctest-kv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  purge_protection_enabled   = true
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctest-certificate"
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

  depends_on = [azurerm_key_vault_access_policy.terraform]
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

resource "azurerm_key_vault_access_policy" "managed" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault.test.tenant_id
  object_id    = "%s"

  key_permissions = [
    "Get",
    "GetRotationPolicy",
    "UnwrapKey",
    "WrapKey",
  ]
}

resource "azurerm_databricks_workspace_serverless" "test" {
  name                                  = "acctest-dbws-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  location                              = azurerm_resource_group.test.location
  managed_services_cmk_key_vault_key_id = azurerm_key_vault_key.test.id

  enhanced_security_compliance {
    automatic_cluster_update_enabled      = true
    compliance_security_profile_enabled   = true
    compliance_security_profile_standards = ["HIPAA"]
    enhanced_security_monitoring_enabled  = true
  }

  tags = {
    Environment = "Sandbox"
    Label       = "Test"
  }

  depends_on = [azurerm_key_vault_access_policy.managed]
}
`, r.template(data), data.RandomString, getDatabricksPrincipalId(data.Client().SubscriptionID), data.RandomInteger)
}

func (r DatabricksWorkspaceServerlessResource) privateLink(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "public" {
  name                 = "acctest-sn-public-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "private" {
  name                 = "acctest-sn-private-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "privatelink" {
  name                 = "acctest-snpl-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]

  private_endpoint_network_policies = "Enabled"
}

resource "azurerm_network_security_group" "nsg" {
  name                = "acctest-nsg-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.public.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = azurerm_subnet.private.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_databricks_workspace_serverless" "test" {
  name                          = "acctest-dbws-%[2]d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  public_network_access_enabled = false

  tags = {
    Environment = "Sandbox"
    Label       = "Test"
  }
}

resource "azurerm_private_endpoint" "databricks" {
  name                = "acctest-endpoint-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.privatelink.id

  private_service_connection {
    name                           = "acctest-psc-%[2]d"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_databricks_workspace_serverless.test.id
    subresource_names              = ["databricks_ui_api"]
  }
}

resource "azurerm_private_dns_zone" "test" {
  name                = "privatelink.azuredatabricks.net"
  resource_group_name = azurerm_resource_group.test.name

  depends_on = [azurerm_private_endpoint.databricks]
}

resource "azurerm_private_dns_cname_record" "test" {
  name                = azurerm_databricks_workspace_serverless.test.workspace_url
  zone_name           = azurerm_private_dns_zone.test.name
  resource_group_name = azurerm_resource_group.test.name
  ttl                 = 300
  record              = "eastus2-c2.azuredatabricks.net"
}
`, r.template(data), data.RandomInteger)
}

func (DatabricksWorkspaceServerlessResource) altSubscriptionCmkServicesOnly(data acceptance.TestData, alt *DatabricksWorkspaceAlternateSubscription) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

provider "azurerm-alt" {
  features {}

  tenant_id       = "%[5]s"
  subscription_id = "%[6]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-pri-sub-services-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "keyVault" {
  provider = azurerm-alt

  name     = "acctestRG-databricks-alt-sub-services-%[1]d"
  location = "%[2]s"
}

resource "azurerm_databricks_workspace_serverless" "test" {
  name                                  = "acctest-databricks-pri-sub-%[1]d"
  resource_group_name                   = azurerm_resource_group.test.name
  location                              = azurerm_resource_group.test.location
  managed_services_cmk_key_vault_id     = azurerm_key_vault.keyVault.id
  managed_services_cmk_key_vault_key_id = azurerm_key_vault_key.services.id

  tags = {
    Environment = "Sandbox"
    Label       = "Test"
  }

  depends_on = [azurerm_key_vault_access_policy.managed]
}

# Create this in a different subscription...
resource "azurerm_key_vault" "keyVault" {
  provider = azurerm-alt

  name                = "kv-altsub-%[3]s"
  location            = azurerm_resource_group.keyVault.location
  resource_group_name = azurerm_resource_group.keyVault.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "services" {
  provider = azurerm-alt

  name         = "acctest-services-certificate"
  key_vault_id = azurerm_key_vault.keyVault.id
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

  depends_on = [azurerm_key_vault_access_policy.terraform]
}

resource "azurerm_key_vault_access_policy" "terraform" {
  provider = azurerm-alt

  key_vault_id = azurerm_key_vault.keyVault.id
  tenant_id    = azurerm_key_vault.keyVault.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "List",
    "Create",
    "Decrypt",
    "Encrypt",
    "Sign",
    "UnwrapKey",
    "Verify",
    "WrapKey",
    "Delete",
    "Restore",
    "Recover",
    "Update",
    "Purge",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}

resource "azurerm_key_vault_access_policy" "managed" {
  provider = azurerm-alt

  key_vault_id = azurerm_key_vault.keyVault.id
  tenant_id    = azurerm_key_vault.keyVault.tenant_id
  object_id    = "%[4]s"

  key_permissions = [
    "Get",
    "UnwrapKey",
    "WrapKey",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, getDatabricksPrincipalId(data.Client().SubscriptionID), alt.tenantID, alt.subscriptionID)
}

func (r DatabricksWorkspaceServerlessResource) enhancedSecurityCompliance(data acceptance.TestData, automaticClusterUpdateEnabled bool, complianceSecurityProfileEnabled bool, enhancedSecurityMonitoringEnabled bool) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_databricks_workspace_serverless" "test" {
  name                = "acctest-dbws-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  enhanced_security_compliance {
    automatic_cluster_update_enabled      = %[3]t
    compliance_security_profile_enabled   = %[4]t
    compliance_security_profile_standards = ["HIPAA"]
    enhanced_security_monitoring_enabled  = %[5]t
  }
}
`, r.template(data), data.RandomInteger, automaticClusterUpdateEnabled, complianceSecurityProfileEnabled, enhancedSecurityMonitoringEnabled)
}
