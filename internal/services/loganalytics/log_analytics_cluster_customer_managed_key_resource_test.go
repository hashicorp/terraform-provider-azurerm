// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogAnalyticsClusterCustomerManagedKeyResource struct{}

func TestAccLogAnalyticsClusterCustomerManagedKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster_customer_managed_key", "test")
	r := LogAnalyticsClusterCustomerManagedKeyResource{}

	if os.Getenv("ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS is not specified")
		return
	}

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

func TestAccLogAnalyticsClusterCustomerManagedKey_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster_customer_managed_key", "test")
	r := LogAnalyticsClusterCustomerManagedKeyResource{}

	if os.Getenv("ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS is not specified")
		return
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
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
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	enabled := false
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if kv := props.KeyVaultProperties; kv != nil {
				enabled = kv.KeyVaultUri != nil && kv.KeyVersion != nil && kv.KeyName != nil
			}
		}
	}

	return utils.Bool(enabled), nil
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
  name                = "vault%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  soft_delete_retention_days = 7
  purge_protection_enabled   = false
}


resource "azurerm_key_vault_access_policy" "terraform" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "List",
    "Purge",
    "Update",
    "GetRotationPolicy",
  ]

  secret_permissions = [
    "Get",
    "Delete",
    "Set",
  ]

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
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

  depends_on = [azurerm_key_vault_access_policy.terraform]
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "Get",
    "UnwrapKey",
    "WrapKey",
    "GetRotationPolicy",
  ]

  tenant_id = azurerm_log_analytics_cluster.test.identity.0.tenant_id
  object_id = azurerm_log_analytics_cluster.test.identity.0.principal_id

  depends_on = [azurerm_key_vault_access_policy.terraform]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r LogAnalyticsClusterCustomerManagedKeyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster_customer_managed_key" "test" {
  log_analytics_cluster_id = azurerm_log_analytics_cluster.test.id
  key_vault_key_id         = azurerm_key_vault_key.test.id

  depends_on = [azurerm_key_vault_access_policy.test]
}
`, r.template(data))
}

func (r LogAnalyticsClusterCustomerManagedKeyResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

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

  depends_on = [azurerm_key_vault_access_policy.terraform]
}

resource "azurerm_log_analytics_cluster_customer_managed_key" "test" {
  log_analytics_cluster_id = azurerm_log_analytics_cluster.test.id
  key_vault_key_id         = azurerm_key_vault_key.test2.id

  depends_on = [azurerm_key_vault_access_policy.test]
}

`, r.template(data), data.RandomString)
}
