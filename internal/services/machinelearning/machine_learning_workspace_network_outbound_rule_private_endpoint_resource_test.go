// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/managednetwork"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WorkspaceNetworkOutboundPrivateEndpointResource struct{}

func TestAccMachineLearningWorkspaceNetworkOutboundRulePrivateEndpoint_onlyApprovedOutbound(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint", "test")
	r := WorkspaceNetworkOutboundPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.onlyApprovedOutbound(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_resource_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningWorkspaceNetworkOutboundRulePrivateEndpoint_internetOutbound(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint", "test")
	r := WorkspaceNetworkOutboundPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.internetOutbound(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_resource_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningWorkspaceNetworkOutboundRulePrivateEndpoint_keyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint", "test")
	r := WorkspaceNetworkOutboundPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withKeyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_resource_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningWorkspaceNetworkOutboundRulePrivateEndpoint_workspace(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint", "test")
	r := WorkspaceNetworkOutboundPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withWorkspace(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_resource_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningWorkspaceNetworkOutboundRulePrivateEndpoint_redis(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint", "test")
	r := WorkspaceNetworkOutboundPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withRedis(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_resource_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningWorkspaceNetworkOutboundRulePrivateEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint", "test")
	r := WorkspaceNetworkOutboundPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.onlyApprovedOutbound(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r WorkspaceNetworkOutboundPrivateEndpointResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := managednetwork.ParseOutboundRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MachineLearning.ManagedNetwork.SettingsRuleGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving Machine Learning Workspace Outbound Rule Private Endpoint %q: %+v", state.ID, err)
	}

	return pointer.To(resp.Model.Properties != nil), nil
}

func (r WorkspaceNetworkOutboundPrivateEndpointResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ml-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestvault%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
    "Delete",
    "Purge",
    "GetRotationPolicy",
  ]
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[4]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomIntOfLength(10))
}

func (r WorkspaceNetworkOutboundPrivateEndpointResource) onlyApprovedOutbound(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%[1]s

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctest-MLW-%[2]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  managed_network {
    isolation_mode = "AllowOnlyApprovedOutbound"
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctestsa%[3]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint" "test" {
  name                = "acctest-MLW-outboundrule-%[3]s"
  workspace_id        = azurerm_machine_learning_workspace.test.id
  service_resource_id = azurerm_storage_account.test2.id
  sub_resource_target = "blob"
}
`, template, data.RandomInteger, data.RandomString)
}

func (r WorkspaceNetworkOutboundPrivateEndpointResource) internetOutbound(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%[1]s

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctest-MLW-%[2]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  managed_network {
    isolation_mode = "AllowInternetOutbound"
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctestsa%[3]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint" "test" {
  name                = "acctest-MLW-outboundrule-%[3]s"
  workspace_id        = azurerm_machine_learning_workspace.test.id
  service_resource_id = azurerm_storage_account.test2.id
  sub_resource_target = "blob"
}
`, template, data.RandomInteger, data.RandomStringOfLength(6))
}

func (r WorkspaceNetworkOutboundPrivateEndpointResource) withKeyVault(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%[1]s

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctest-MLW-%[2]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  managed_network {
    isolation_mode = "AllowOnlyApprovedOutbound"
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault" "test2" {
  name                = "acctestvault%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test2.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
    "Delete",
    "Purge",
    "GetRotationPolicy",
  ]
}

resource "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint" "test" {
  name                = "acctest-MLW-outboundrule-%[3]s"
  workspace_id        = azurerm_machine_learning_workspace.test.id
  service_resource_id = azurerm_key_vault.test2.id
  sub_resource_target = "vault"
}
`, template, data.RandomInteger, data.RandomStringOfLength(6))
}

func (r WorkspaceNetworkOutboundPrivateEndpointResource) withWorkspace(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%[1]s

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctest-MLW-%[2]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  managed_network {
    isolation_mode = "AllowOnlyApprovedOutbound"
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_machine_learning_workspace" "test2" {
  name                    = "acctest-MLW2-%[2]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint" "test" {
  name                = "acctest-MLW-outboundrule-%[3]s"
  workspace_id        = azurerm_machine_learning_workspace.test.id
  service_resource_id = azurerm_machine_learning_workspace.test2.id
  sub_resource_target = "amlworkspace"
}
`, template, data.RandomInteger, data.RandomStringOfLength(6))
}

func (r WorkspaceNetworkOutboundPrivateEndpointResource) withRedis(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%[1]s

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctest-MLW-%[2]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  managed_network {
    isolation_mode = "AllowOnlyApprovedOutbound"
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_redis_cache" "test" {
  name                 = "acctest-mlwcache-%[3]s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  capacity             = 2
  family               = "C"
  sku_name             = "Standard"
  non_ssl_port_enabled = false
  minimum_tls_version  = "1.2"

  redis_configuration {
  }
}

resource "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint" "test" {
  name                = "acctest-MLW-outboundrule-%[3]s"
  workspace_id        = azurerm_machine_learning_workspace.test.id
  service_resource_id = azurerm_redis_cache.test.id
  sub_resource_target = "redisCache"
}
`, template, data.RandomInteger, data.RandomStringOfLength(6))
}

func (r WorkspaceNetworkOutboundPrivateEndpointResource) requiresImport(data acceptance.TestData) string {
	template := r.onlyApprovedOutbound(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint" "import" {
  name                = azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint.test.name
  workspace_id        = azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint.test.workspace_id
  service_resource_id = azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint.test.service_resource_id
  sub_resource_target = azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint.test.sub_resource_target
}
`, template)
}
