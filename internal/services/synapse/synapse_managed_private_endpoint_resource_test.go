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

type SynapseManagedPrivateEndpointResource struct{}

func TestAccSynapseManagedPrivateEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

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

func TestAccSynapseManagedPrivateEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

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

func TestAccSynapseManagedPrivateEndpoint_autoApproveCognitive(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoApproveCognitive(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseManagedPrivateEndpoint_autoApproveCosmos(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoApproveCosmos(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseManagedPrivateEndpoint_autoApproveKeyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoApproveKeyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseManagedPrivateEndpoint_autoApproveMariaDB(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoApproveMariaDB(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseManagedPrivateEndpoint_autoApproveMySQL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoApproveMySQL(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseManagedPrivateEndpoint_autoApprovePostgreSQL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoApprovePostgreSQL(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseManagedPrivateEndpoint_autoApprovePurview(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoApprovePurview(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseManagedPrivateEndpoint_autoApproveSQL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoApproveSQL(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseManagedPrivateEndpoint_autoApproveSearch(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoApproveSearch(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseManagedPrivateEndpoint_autoApproveStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoApproveStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseManagedPrivateEndpoint_autoApproveSynapse(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoApproveSynapse(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SynapseManagedPrivateEndpointResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ManagedPrivateEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	environment := client.Account.Environment
	managedPrivateEndpointsClient, err := client.Synapse.ManagedPrivateEndpointsClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return nil, err
	}
	resp, err := managedPrivateEndpointsClient.Get(ctx, id.ManagedVirtualNetworkName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if state.Attributes["is_manual_connection"] == "false" && *resp.Properties.ConnectionState.Status != "Approved" {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (r SynapseManagedPrivateEndpointResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test_endpoint" {
  name                     = "acctestacce%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_synapse_managed_private_endpoint" "test" {
  name                 = "acctestEndpoint%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  target_resource_id   = azurerm_storage_account.test_endpoint.id
  subresource_name     = "blob"

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger, data.RandomString)
}

func (r SynapseManagedPrivateEndpointResource) autoApproveCognitive(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account" "test" {
  name                  = "acctestcogacc-%[2]d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  kind                  = "Face"
  sku_name              = "F0"
  custom_subdomain_name = "acctestcogacc-%[2]d"

  public_network_access_enabled = false

  network_acls {
    default_action = "Deny"
  }
}

resource "azurerm_synapse_managed_private_endpoint" "test" {
  name                 = "acctestEndpoint%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  target_resource_id   = azurerm_cognitive_account.test.id
  subresource_name     = "account"
  is_manual_connection = false

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r SynapseManagedPrivateEndpointResource) autoApproveCosmos(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Eventual"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_synapse_managed_private_endpoint" "test" {
  name                 = "acctestEndpoint%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  target_resource_id   = azurerm_cosmosdb_account.test.id
  subresource_name     = "sql"
  is_manual_connection = false

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r SynapseManagedPrivateEndpointResource) autoApproveKeyVault(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%[2]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_synapse_managed_private_endpoint" "test" {
  name                 = "acctestEndpoint%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  target_resource_id   = azurerm_key_vault.test.id
  subresource_name     = "vault"
  is_manual_connection = false

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r SynapseManagedPrivateEndpointResource) autoApproveMariaDB(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "GP_Gen5_2"
  version             = "10.3"

  public_network_access_enabled = false

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  ssl_enforcement_enabled      = true
  storage_mb                   = 51200
}

resource "azurerm_synapse_managed_private_endpoint" "test" {
  name                 = "acctestEndpoint%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  target_resource_id   = azurerm_mariadb_server.test.id
  subresource_name     = "mariadbServer"
  is_manual_connection = false

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r SynapseManagedPrivateEndpointResource) autoApproveMySQL(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mysql_server" "test" {
  name                             = "acctestmysqlsvr-%[2]d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  sku_name                         = "GP_Gen5_2"
  administrator_login              = "acctestun"
  administrator_login_password     = "H@Sh1CoR3!"
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_1"
  storage_mb                       = 51200
  version                          = "5.7"
}

resource "azurerm_synapse_managed_private_endpoint" "test" {
  name                 = "acctestEndpoint%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  target_resource_id   = azurerm_mysql_server.test.id
  subresource_name     = "mysqlServer"
  is_manual_connection = false

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r SynapseManagedPrivateEndpointResource) autoApprovePostgreSQL(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "GP_Gen5_2"
  version             = "9.5"

  public_network_access_enabled = false

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  ssl_enforcement_enabled      = true
  storage_mb                   = 51200
}

resource "azurerm_synapse_managed_private_endpoint" "test" {
  name                 = "acctestEndpoint%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  target_resource_id   = azurerm_postgresql_server.test.id
  subresource_name     = "postgresqlServer"
  is_manual_connection = false

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r SynapseManagedPrivateEndpointResource) autoApprovePurview(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_purview_account" "test" {
  name                = "acctestsw%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_synapse_managed_private_endpoint" "test" {
  name                 = "acctestEndpoint%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  target_resource_id   = azurerm_purview_account.test.id
  subresource_name     = "portal"
  is_manual_connection = false

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r SynapseManagedPrivateEndpointResource) autoApproveSearch(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_search_service" "test" {
  name                = "acctestsearchservice%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"

  public_network_access_enabled = false

  tags = {
    environment = "staging"
  }
}

resource "azurerm_synapse_managed_private_endpoint" "test" {
  name                 = "acctestEndpoint%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  target_resource_id   = azurerm_search_service.test.id
  subresource_name     = "searchService"
  is_manual_connection = false

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r SynapseManagedPrivateEndpointResource) autoApproveSQL(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_synapse_managed_private_endpoint" "test" {
  name                 = "acctestEndpoint%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  target_resource_id   = azurerm_sql_server.test.id
  subresource_name     = "SQLServer"
  is_manual_connection = false

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r SynapseManagedPrivateEndpointResource) autoApproveStorage(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test_endpoint" {
  name                     = "acctestacce%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_synapse_managed_private_endpoint" "test" {
  name                 = "acctestEndpoint%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  target_resource_id   = azurerm_storage_account.test_endpoint.id
  subresource_name     = "blob"
  is_manual_connection = false

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger, data.RandomString)
}

func (r SynapseManagedPrivateEndpointResource) autoApproveSynapse(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_data_lake_gen2_filesystem" "test_endpoint" {
  name               = "accteste-%[2]d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_synapse_workspace" "test_endpoint" {
  name                                 = "acctestswe%[2]d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test_endpoint.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
}

resource "azurerm_synapse_managed_private_endpoint" "test" {
  name                 = "acctestEndpoint%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  target_resource_id   = azurerm_synapse_workspace.test_endpoint.id
  subresource_name     = "sqlOnDemand"
  is_manual_connection = false

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger, data.RandomString)
}

func (r SynapseManagedPrivateEndpointResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_managed_private_endpoint" "import" {
  name                 = azurerm_synapse_managed_private_endpoint.test.name
  synapse_workspace_id = azurerm_synapse_managed_private_endpoint.test.synapse_workspace_id
  target_resource_id   = azurerm_synapse_managed_private_endpoint.test.target_resource_id
  subresource_name     = azurerm_synapse_managed_private_endpoint.test.subresource_name
}
`, config)
}

func (r SynapseManagedPrivateEndpointResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-synapse-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%[1]d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%[1]d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  managed_virtual_network_enabled      = true
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_synapse_firewall_rule" "test" {
  name                 = "AllowAll"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
