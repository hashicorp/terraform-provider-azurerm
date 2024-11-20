// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/managednetwork"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkspaceNetworkOutboundRuleServiceTagResource struct{}

func TestAccMachineLearningWorkspaceNetworkOutboundRuleServiceTag_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_network_outbound_rule_service_tag", "test")
	r := WorkspaceNetworkOutboundRuleServiceTagResource{}

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

func TestAccMachineLearningWorkspaceNetworkOutboundRuleServiceTag_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_network_outbound_rule_service_tag", "test")
	r := WorkspaceNetworkOutboundRuleServiceTagResource{}

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
				check.That(data.ResourceName).Key("service_tag").HasValue("AppService"),
				check.That(data.ResourceName).Key("protocol").HasValue("UDP"),
				check.That(data.ResourceName).Key("port_ranges").HasValue("445"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningWorkspaceNetworkOutboundRuleServiceTag_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_network_outbound_rule_service_tag", "test")
	r := WorkspaceNetworkOutboundRuleServiceTagResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_tag").Exists(),
				check.That(data.ResourceName).Key("protocol").Exists(),
				check.That(data.ResourceName).Key("port_ranges").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r WorkspaceNetworkOutboundRuleServiceTagResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := managednetwork.ParseOutboundRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MachineLearning.ManagedNetwork.SettingsRuleGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Machine Learning Workspace Outbound Rule Service Tag %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.Model.Properties != nil), nil
}

func (r WorkspaceNetworkOutboundRuleServiceTagResource) template(data acceptance.TestData) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomIntOfLength(10))
}

func (r WorkspaceNetworkOutboundRuleServiceTagResource) basic(data acceptance.TestData) string {
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

resource "azurerm_machine_learning_workspace_network_outbound_rule_service_tag" "test" {
  name         = "acctest-MLW-outboundrule-%[3]s"
  workspace_id = azurerm_machine_learning_workspace.test.id
  service_tag  = "AppConfiguration"
  protocol     = "TCP"
  port_ranges  = "443"
}
`, template, data.RandomInteger, data.RandomString)
}

func (r WorkspaceNetworkOutboundRuleServiceTagResource) update(data acceptance.TestData) string {
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

resource "azurerm_machine_learning_workspace_network_outbound_rule_service_tag" "test" {
  name         = "acctest-MLW-outboundrule-%[3]s"
  workspace_id = azurerm_machine_learning_workspace.test.id
  service_tag  = "AppService"
  protocol     = "UDP"
  port_ranges  = "445"
}
`, template, data.RandomInteger, data.RandomString)
}

func (r WorkspaceNetworkOutboundRuleServiceTagResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_workspace_network_outbound_rule_service_tag" "import" {
  name         = azurerm_machine_learning_workspace_network_outbound_rule_service_tag.test.name
  workspace_id = azurerm_machine_learning_workspace_network_outbound_rule_service_tag.test.workspace_id
  service_tag  = azurerm_machine_learning_workspace_network_outbound_rule_service_tag.test.service_tag
  protocol     = azurerm_machine_learning_workspace_network_outbound_rule_service_tag.test.protocol
  port_ranges  = azurerm_machine_learning_workspace_network_outbound_rule_service_tag.test.port_ranges
}
`, template)
}
