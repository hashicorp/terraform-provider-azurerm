// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/sentinelonboardingstates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SecurityInsightsSentinelOnboardingStateResource struct{}

func TestAccSecurityInsightsSentinelOnboardingState_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_log_analytics_workspace_onboarding", "test")
	r := SecurityInsightsSentinelOnboardingStateResource{}
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

func TestAccSecurityInsightsSentinelOnboardingState_basicWithName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_log_analytics_workspace_onboarding", "test")
	r := SecurityInsightsSentinelOnboardingStateResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityInsightsSentinelOnboardingState_ToggleCmkEnabled(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS is not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_sentinel_log_analytics_workspace_onboarding", "test")
	r := SecurityInsightsSentinelOnboardingStateResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCmk(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityInsightsSentinelOnboardingState_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_log_analytics_workspace_onboarding", "test")
	r := SecurityInsightsSentinelOnboardingStateResource{}
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

func (r SecurityInsightsSentinelOnboardingStateResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := sentinelonboardingstates.ParseOnboardingStateID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Sentinel.OnboardingStatesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r SecurityInsightsSentinelOnboardingStateResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "test" {
  workspace_id = azurerm_log_analytics_workspace.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r SecurityInsightsSentinelOnboardingStateResource) basicWithName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "test" {
  workspace_id = azurerm_log_analytics_workspace.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r SecurityInsightsSentinelOnboardingStateResource) withCmk(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_cluster" "test" {
  name                = "acctest-LA-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}

data "azuread_service_principal" "cosmos" {
  display_name = "Azure Cosmos DB"
}


resource "azurerm_key_vault" "test" {
  name                = "vault%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  soft_delete_retention_days = 7
  purge_protection_enabled   = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    key_permissions = [
      "Create",
      "Delete",
      "Get",
      "List",
      "Purge",
      "Update",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "Set",
    ]
  }

  access_policy {
    tenant_id = azurerm_log_analytics_cluster.test.identity.0.tenant_id
    object_id = azurerm_log_analytics_cluster.test.identity.0.principal_id
    key_permissions = [
      "Get",
      "UnwrapKey",
      "WrapKey"
    ]
  }

  access_policy {
    tenant_id = azurerm_log_analytics_cluster.test.identity.0.tenant_id
    object_id = data.azuread_service_principal.cosmos.object_id
    key_permissions = [
      "Get",
      "UnwrapKey",
      "WrapKey"
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%[3]s"
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

resource "azurerm_log_analytics_cluster_customer_managed_key" "test" {
  log_analytics_cluster_id = azurerm_log_analytics_cluster.test.id
  key_vault_key_id         = azurerm_key_vault_key.test.id

}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-law-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
  lifecycle {
    ignore_changes = [sku]
  }
}

resource "azurerm_log_analytics_linked_service" "test" {
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id
  write_access_id     = azurerm_log_analytics_cluster.test.id

  depends_on = [azurerm_log_analytics_cluster_customer_managed_key.test]
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "test" {
  workspace_id                 = azurerm_log_analytics_workspace.test.id
  customer_managed_key_enabled = true

  depends_on = [azurerm_log_analytics_linked_service.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r SecurityInsightsSentinelOnboardingStateResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "import" {
  workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
}
`, config)
}
