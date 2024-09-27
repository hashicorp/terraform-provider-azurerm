// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MsSqlDatabaseExtendedAuditingPolicyResource struct{}

func TestAccMsSqlDatabaseExtendedAuditingPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_extended_auditing_policy", "test")
	r := MsSqlDatabaseExtendedAuditingPolicyResource{}

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

func TestAccMsSqlDatabaseExtendedAuditingPolicy_primaryWithOnlineSecondary(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_extended_auditing_policy", "test")
	r := MsSqlDatabaseExtendedAuditingPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.primaryWithOnlineSecondary(data, "test1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
		{
			Config: r.primaryWithOnlineSecondary(data, "test2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
	})
}

func TestAccMsSqlDatabaseExtendedAuditingPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_extended_auditing_policy", "test")
	r := MsSqlDatabaseExtendedAuditingPolicyResource{}

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

func TestAccMsSqlDatabaseExtendedAuditingPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_extended_auditing_policy", "test")
	r := MsSqlDatabaseExtendedAuditingPolicyResource{}

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

func TestAccMsSqlDatabaseExtendedAuditingPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_extended_auditing_policy", "test")
	r := MsSqlDatabaseExtendedAuditingPolicyResource{}

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
			Config: r.disabled(data),
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

func TestAccMsSqlDatabaseExtendedAuditingPolicy_storageAccBehindFireWall(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_extended_auditing_policy", "test")
	r := MsSqlDatabaseExtendedAuditingPolicyResource{}

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

func TestAccMsSqlDatabaseExtendedAuditingPolicy_logAnalytics(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_extended_auditing_policy", "test")
	r := MsSqlDatabaseExtendedAuditingPolicyResource{}

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

func TestAccMsSqlDatabaseExtendedAuditingPolicy_eventhub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_extended_auditing_policy", "test")
	r := MsSqlDatabaseExtendedAuditingPolicyResource{}

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

func (MsSqlDatabaseExtendedAuditingPolicyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DatabaseExtendedAuditingPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	dbId := commonids.NewSqlDatabaseID(id.SubscriptionId, id.ResourceGroup, id.ServerName, id.DatabaseName)

	resp, err := client.MSSQL.BlobAuditingPoliciesClient.ExtendedDatabaseBlobAuditingPoliciesGet(ctx, dbId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("%s does not exist", *id)
		}

		return nil, fmt.Errorf("reading %s: %v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (MsSqlDatabaseExtendedAuditingPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest-sqlserver-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "AdminPassword123!"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[1]d"
  server_id = azurerm_mssql_server.test.id
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r MsSqlDatabaseExtendedAuditingPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database_extended_auditing_policy" "test" {
  database_id                = azurerm_mssql_database.test.id
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, r.template(data))
}

func (r MsSqlDatabaseExtendedAuditingPolicyResource) primaryWithOnlineSecondary(data acceptance.TestData, tag string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "second" {
  name     = "acctestRG-mssql2-%[2]d"
  location = "%[4]s"
}

resource "azurerm_mssql_server" "second" {
  name                         = "acctest-sqlserver2-%[2]d"
  resource_group_name          = azurerm_resource_group.second.name
  location                     = azurerm_resource_group.second.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "secondary" {
  name                        = "acctest-dbs-%[2]d"
  server_id                   = azurerm_mssql_server.second.id
  create_mode                 = "OnlineSecondary"
  creation_source_database_id = azurerm_mssql_database.test.id

  tags = {
    tag = "%[5]s"
  }
}

resource "azurerm_mssql_database_extended_auditing_policy" "test" {
  database_id                = azurerm_mssql_database.test.id
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  depends_on = [azurerm_mssql_database.secondary]
}
`, r.template(data), data.RandomInteger, data.RandomString, data.Locations.Secondary, tag)
}

func (r MsSqlDatabaseExtendedAuditingPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database_extended_auditing_policy" "import" {
  database_id                = azurerm_mssql_database.test.id
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, r.template(data))
}

func (r MsSqlDatabaseExtendedAuditingPolicyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database_extended_auditing_policy" "test" {
  database_id                             = azurerm_mssql_database.test.id
  storage_endpoint                        = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.test.primary_access_key
  storage_account_access_key_is_secondary = false
  retention_in_days                       = 6
  log_monitoring_enabled                  = false
  enabled                                 = true
}
`, r.template(data))
}

func (r MsSqlDatabaseExtendedAuditingPolicyResource) disabled(data acceptance.TestData) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest-sqlserver-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "AdminPassword123!"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[1]d"
  server_id = azurerm_mssql_server.test.id
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database_extended_auditing_policy" "test" {
  database_id = azurerm_mssql_database.test.id
  enabled     = false
}

`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r MsSqlDatabaseExtendedAuditingPolicyResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test2" {
  name                     = "unlikely23exst2acc2%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database_extended_auditing_policy" "test" {
  database_id                             = azurerm_mssql_database.test.id
  storage_endpoint                        = azurerm_storage_account.test2.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.test2.primary_access_key
  storage_account_access_key_is_secondary = true
  retention_in_days                       = 3
}
`, r.template(data), data.RandomString)
}

func (MsSqlDatabaseExtendedAuditingPolicyResource) storageAccountBehindFireWall(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest-sqlserver-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "AdminPassword123!"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[1]d"
  server_id = azurerm_mssql_server.test.id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                              = "acctestsubnet%[1]d"
  resource_group_name               = azurerm_resource_group.test.name
  virtual_network_name              = azurerm_virtual_network.test.name
  address_prefixes                  = ["10.0.2.0/24"]
  service_endpoints                 = ["Microsoft.Storage", "Microsoft.Sql"]
  private_endpoint_network_policies = "Disabled"
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  network_rules {
    default_action             = "Deny"
    ip_rules                   = ["127.0.0.1"]
    virtual_network_subnet_ids = [azurerm_subnet.test.id]
    bypass                     = ["AzureServices"]
  }

  identity {
    type = "SystemAssigned"

  }
}

resource "azurerm_mssql_virtual_network_rule" "sqlvnetrule" {
  name      = "sql-vnet-rule"
  server_id = azurerm_mssql_server.test.id
  subnet_id = azurerm_subnet.test.id

}

resource "azurerm_mssql_firewall_rule" "test" {
  name             = "FirewallRule1"
  server_id        = azurerm_mssql_server.test.id
  start_ip_address = "0.0.0.0"
  end_ip_address   = "0.0.0.0"
}


resource "azurerm_mssql_database_extended_auditing_policy" "test" {
  database_id      = azurerm_mssql_database.test.id
  storage_endpoint = azurerm_storage_account.test.primary_blob_endpoint

  depends_on = [
    azurerm_role_assignment.test,
  ]
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_role_assignment" "test" {
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_mssql_server.test.identity.0.principal_id
}

resource "azurerm_mssql_server_extended_auditing_policy" "test" {
  storage_endpoint       = azurerm_storage_account.test.primary_blob_endpoint
  server_id              = azurerm_mssql_server.test.id
  retention_in_days      = 6
  log_monitoring_enabled = false

  storage_account_subscription_id = "%s"

  depends_on = [
    azurerm_role_assignment.test,
    azurerm_storage_account.test,
  ]
}


`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.Client().SubscriptionID)
}

func (r MsSqlDatabaseExtendedAuditingPolicyResource) monitorTemplate(data acceptance.TestData) string {
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

resource "azurerm_mssql_server_extended_auditing_policy" "test" {
  server_id              = azurerm_mssql_server.test.id
  log_monitoring_enabled = true
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseExtendedAuditingPolicyResource) logAnalytics(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "acctest-DS-%[2]d"
  target_resource_id         = azurerm_mssql_database.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  enabled_log {
    category = "SQLSecurityAuditEvents"

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

  // enabled_log, metric will return all disabled categories
  lifecycle {
    ignore_changes = [enabled_log, metric]
  }
}

resource "azurerm_mssql_database_extended_auditing_policy" "test" {
  database_id            = azurerm_mssql_database.test.id
  log_monitoring_enabled = true
}
`, r.monitorTemplate(data), data.RandomInteger)
}

func (r MsSqlDatabaseExtendedAuditingPolicyResource) logAnalyticsAndStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "acctest-DS-%[2]d"
  target_resource_id         = azurerm_mssql_database.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  enabled_log {
    category = "SQLSecurityAuditEvents"

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

  // enabled_log, metric will return all disabled categories
  lifecycle {
    ignore_changes = [enabled_log, metric]
  }
}

resource "azurerm_mssql_database_extended_auditing_policy" "test" {
  database_id                = azurerm_mssql_database.test.id
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  log_monitoring_enabled     = true
}
`, r.monitorTemplate(data), data.RandomInteger)
}

func (r MsSqlDatabaseExtendedAuditingPolicyResource) eventhub(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                           = "acctest-DS-%[2]d"
  target_resource_id             = azurerm_mssql_database.test.id
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.test.id
  eventhub_name                  = azurerm_eventhub.test.name


  enabled_log {
    category = "SQLSecurityAuditEvents"

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

  // enabled_log, metric will return all disabled categories
  lifecycle {
    ignore_changes = [enabled_log, metric]
  }

}

resource "azurerm_mssql_database_extended_auditing_policy" "test" {
  database_id            = azurerm_mssql_database.test.id
  log_monitoring_enabled = true
}
`, r.monitorTemplate(data), data.RandomInteger)
}

func (r MsSqlDatabaseExtendedAuditingPolicyResource) eventhubAndStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                           = "acctest-DS-%[2]d"
  target_resource_id             = azurerm_mssql_database.test.id
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.test.id
  eventhub_name                  = azurerm_eventhub.test.name


  enabled_log {
    category = "SQLSecurityAuditEvents"

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

  // enabled_log, metric will return all disabled categories
  lifecycle {
    ignore_changes = [enabled_log, metric]
  }

}

resource "azurerm_mssql_database_extended_auditing_policy" "test" {
  database_id                = azurerm_mssql_database.test.id
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  log_monitoring_enabled     = true
}
`, r.monitorTemplate(data), data.RandomInteger)
}
