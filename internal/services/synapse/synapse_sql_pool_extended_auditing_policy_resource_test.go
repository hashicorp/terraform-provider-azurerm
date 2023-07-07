// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SynapseSqlPoolExtendedAuditingPolicyResource struct{}

func TestAccSynapseSqlPoolExtendedAuditingPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool_extended_auditing_policy", "test")
	r := SynapseSqlPoolExtendedAuditingPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
	})
}

func TestAccSynapseSqlPoolExtendedAuditingPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool_extended_auditing_policy", "test")
	r := SynapseSqlPoolExtendedAuditingPolicyResource{}

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

func TestAccSynapseSqlPoolExtendedAuditingPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool_extended_auditing_policy", "test")
	r := SynapseSqlPoolExtendedAuditingPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
	})
}

func TestAccSynapseSqlPoolExtendedAuditingPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool_extended_auditing_policy", "test")
	r := SynapseSqlPoolExtendedAuditingPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
	})
}

func TestAccSynapseSqlPoolExtendedAuditingPolicy_storageAccBehindFireWall(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool_extended_auditing_policy", "test")
	r := SynapseSqlPoolExtendedAuditingPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountBehindFireWall(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
	})
}

func TestAccSynapseSqlPoolExtendedAuditingPolicy_logAnalytics(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool_extended_auditing_policy", "test")
	r := SynapseSqlPoolExtendedAuditingPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
		{
			Config: r.logAnalytics(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
		{
			Config: r.logAnalyticsAndStorageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
	})
}

func TestAccSynapseSqlPoolExtendedAuditingPolicy_eventhub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool_extended_auditing_policy", "test")
	r := SynapseSqlPoolExtendedAuditingPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
		{
			Config: r.eventhub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
		{
			Config: r.eventhubAndStorageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
	})
}

func (SynapseSqlPoolExtendedAuditingPolicyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SqlPoolExtendedAuditingPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Synapse.SqlPoolExtendedBlobAuditingPoliciesClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (SynapseSqlPoolExtendedAuditingPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestsw%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "sw" {
  name                     = "acctestsw%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%[1]d"
  storage_account_id = azurerm_storage_account.sw.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%[1]d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_synapse_sql_pool" "test" {
  name                 = "acctestSP%s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  sku_name             = "DW100c"
  create_mode          = "Default"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r SynapseSqlPoolExtendedAuditingPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_sql_pool_extended_auditing_policy" "test" {
  sql_pool_id                = azurerm_synapse_sql_pool.test.id
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, r.template(data))
}

func (r SynapseSqlPoolExtendedAuditingPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_sql_pool_extended_auditing_policy" "import" {
  sql_pool_id                = azurerm_synapse_sql_pool.test.id
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, r.template(data))
}

func (r SynapseSqlPoolExtendedAuditingPolicyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_sql_pool_extended_auditing_policy" "test" {
  sql_pool_id                             = azurerm_synapse_sql_pool.test.id
  storage_endpoint                        = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.test.primary_access_key
  storage_account_access_key_is_secondary = false
  retention_in_days                       = 6
}
`, r.template(data))
}

func (r SynapseSqlPoolExtendedAuditingPolicyResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test2" {
  name                     = "unlikely23exst2acc2%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_synapse_sql_pool_extended_auditing_policy" "test" {
  sql_pool_id                             = azurerm_synapse_sql_pool.test.id
  storage_endpoint                        = azurerm_storage_account.test2.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.test2.primary_access_key
  storage_account_access_key_is_secondary = true
  retention_in_days                       = 3
}
`, r.template(data), data.RandomString)
}

func (r SynapseSqlPoolExtendedAuditingPolicyResource) storageAccountBehindFireWall(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account_network_rules" "test" {
  storage_account_id         = azurerm_storage_account.test.id
  default_action             = "Deny"
  ip_rules                   = ["127.0.0.1"]
  virtual_network_subnet_ids = [azurerm_subnet.test.id]
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_synapse_workspace.test.identity.0.principal_id
}

resource "azurerm_synapse_sql_pool_extended_auditing_policy" "test" {
  sql_pool_id      = azurerm_synapse_sql_pool.test.id
  storage_endpoint = azurerm_storage_account.test.primary_blob_endpoint

  depends_on = [
    azurerm_role_assignment.test,
    azurerm_storage_account_network_rules.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r SynapseSqlPoolExtendedAuditingPolicyResource) monitorTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-LAW-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-EH-%[2]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "acctestEHRule"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  listen              = true
  send                = true
  manage              = true
}
`, r.template(data), data.RandomInteger)
}

func (r SynapseSqlPoolExtendedAuditingPolicyResource) logAnalytics(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "acctest-DS-%[2]d"
  target_resource_id         = azurerm_synapse_sql_pool.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  log {
    category = "SQLSecurityAuditEvents"
    enabled  = true

    retention_policy {
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"

    retention_policy {
      enabled = false
    }
  }

  // log, metric will return all disabled categories
  lifecycle {
    ignore_changes = [log, metric]
  }
}

resource "azurerm_synapse_sql_pool_extended_auditing_policy" "test" {
  sql_pool_id            = azurerm_synapse_sql_pool.test.id
  log_monitoring_enabled = true
}
`, r.monitorTemplate(data), data.RandomInteger)
}

func (r SynapseSqlPoolExtendedAuditingPolicyResource) logAnalyticsAndStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "acctest-DS-%[2]d"
  target_resource_id         = azurerm_synapse_sql_pool.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  log {
    category = "SQLSecurityAuditEvents"
    enabled  = true

    retention_policy {
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"

    retention_policy {
      enabled = false
    }
  }

  // log, metric will return all disabled categories
  lifecycle {
    ignore_changes = [log, metric]
  }
}

resource "azurerm_synapse_sql_pool_extended_auditing_policy" "test" {
  sql_pool_id                = azurerm_synapse_sql_pool.test.id
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  log_monitoring_enabled     = true
}
`, r.monitorTemplate(data), data.RandomInteger)
}

func (r SynapseSqlPoolExtendedAuditingPolicyResource) eventhub(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                           = "acctest-DS-%[2]d"
  target_resource_id             = azurerm_synapse_sql_pool.test.id
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.test.id
  eventhub_name                  = azurerm_eventhub.test.name


  log {
    category = "SQLSecurityAuditEvents"
    enabled  = true

    retention_policy {
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"

    retention_policy {
      enabled = false
    }
  }

  // log, metric will return all disabled categories
  lifecycle {
    ignore_changes = [log, metric]
  }

}

resource "azurerm_synapse_sql_pool_extended_auditing_policy" "test" {
  sql_pool_id            = azurerm_synapse_sql_pool.test.id
  log_monitoring_enabled = true
}
`, r.monitorTemplate(data), data.RandomInteger)
}

func (r SynapseSqlPoolExtendedAuditingPolicyResource) eventhubAndStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                           = "acctest-DS-%[2]d"
  target_resource_id             = azurerm_synapse_sql_pool.test.id
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.test.id
  eventhub_name                  = azurerm_eventhub.test.name


  log {
    category = "SQLSecurityAuditEvents"
    enabled  = true

    retention_policy {
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"

    retention_policy {
      enabled = false
    }
  }

  // log, metric will return all disabled categories
  lifecycle {
    ignore_changes = [log, metric]
  }

}

resource "azurerm_synapse_sql_pool_extended_auditing_policy" "test" {
  sql_pool_id                = azurerm_synapse_sql_pool.test.id
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  log_monitoring_enabled     = true
}
`, r.monitorTemplate(data), data.RandomInteger)
}
