package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PostgreSQLServerResource struct {
}

func TestAccPostgreSQLServer_basicNinePointFive(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "9.5"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_basicNinePointFiveDeprecated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDeprecated(data, "9.5"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.gp(data, "9.5"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_basicNinePointSix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "9.6"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_basicTenPointZero(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "10.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_gpTenPointZero(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gp(data, "10.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_moTenPointZero(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mo(data, "10.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_basicEleven(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_basicWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithIdentity(data, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_autogrowOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autogrow(data, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.gp(data, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "10.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccPostgreSQLServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_updatedDeprecated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDeprecated(data, "9.6"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.completeDeprecated(data, "9.6"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.basicDeprecated(data, "9.6"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_updated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gp(data, "9.6"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.complete2(data, "9.6"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_completeDeprecatedUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeDeprecated(data, "9.6"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_updateSKU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sku(data, "10.0", "B_Gen5_2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.sku(data, "10.0", "GP_Gen5_2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.sku(data, "10.0", "MO_Gen5_16"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_createReplica(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gp(data, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.createReplica(data, "GP_Gen5_2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_updateReplicaToDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "replica")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.createReplica(data, "GP_Gen5_2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateReplicaToDefault(data, "GP_Gen5_2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("creation_source_server_id"),
	})
}

func TestAccPostgreSQLServer_scaleReplicas(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.createReplicas(data, "GP_Gen5_2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
				check.That("azurerm_postgresql_server.replica1").ExistsInAzure(r),
				check.That("azurerm_postgresql_server.replica1").Key("sku_name").HasValue("GP_Gen5_2"),
				check.That("azurerm_postgresql_server.replica2").ExistsInAzure(r),
				check.That("azurerm_postgresql_server.replica2").Key("sku_name").HasValue("GP_Gen5_2"),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.createReplicas(data, "GP_Gen5_4"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_4"),
				check.That("azurerm_postgresql_server.replica1").ExistsInAzure(r),
				check.That("azurerm_postgresql_server.replica1").Key("sku_name").HasValue("GP_Gen5_4"),
				check.That("azurerm_postgresql_server.replica2").ExistsInAzure(r),
				check.That("azurerm_postgresql_server.replica2").Key("sku_name").HasValue("GP_Gen5_4"),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.createReplicas(data, "GP_Gen5_2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
				check.That("azurerm_postgresql_server.replica1").ExistsInAzure(r),
				check.That("azurerm_postgresql_server.replica1").Key("sku_name").HasValue("GP_Gen5_2"),
				check.That("azurerm_postgresql_server.replica2").ExistsInAzure(r),
				check.That("azurerm_postgresql_server.replica2").Key("sku_name").HasValue("GP_Gen5_2"),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_createPointInTimeRestore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	restoreTime := time.Now().Add(30 * time.Minute)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gp(data, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			PreConfig: func() { time.Sleep(30 * time.Minute) },
			Config:    r.createPointInTimeRestore(data, "11", restoreTime.Format(time.RFC3339)),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccPostgreSQLServer_threatDetectionEmptyAttrs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.emptyAttrs(data, "9.5"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestMinTlsVersionOnServerUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	r := PostgreSQLServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.beforeUpdate(data, "9.6", "TLS1_2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.afterUpdate(data, "9.6", "TLS1_0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (t PostgreSQLServerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.ServersClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Postgresql Server (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (PostgreSQLServerResource) template(data acceptance.TestData, sku, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "%s"
  version    = "%s"
  storage_mb = 51200

  ssl_enforcement_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sku, version)
}

func (r PostgreSQLServerResource) basic(data acceptance.TestData, version string) string {
	return r.template(data, "B_Gen5_1", version)
}

func (PostgreSQLServerResource) basicWithIdentity(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "B_Gen5_1"
  version    = "%s"
  storage_mb = 51200

  ssl_enforcement_enabled = true

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func (r PostgreSQLServerResource) mo(data acceptance.TestData, version string) string {
	return r.template(data, "MO_Gen5_2", version)
}

func (r PostgreSQLServerResource) gp(data acceptance.TestData, version string) string {
	return r.template(data, "GP_Gen5_2", version)
}

func (PostgreSQLServerResource) autogrow(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name          = "GP_Gen5_2"
  version           = "%s"
  auto_grow_enabled = true

  ssl_enforcement_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func (PostgreSQLServerResource) basicDeprecated(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name = "GP_Gen5_2"
  version  = "%s"

  storage_profile {
    storage_mb = 51200
  }

  ssl_enforcement_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func (r PostgreSQLServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_server" "import" {
  name                = azurerm_postgresql_server.test.name
  location            = azurerm_postgresql_server.test.location
  resource_group_name = azurerm_postgresql_server.test.resource_group_name

  administrator_login          = azurerm_postgresql_server.test.administrator_login
  administrator_login_password = azurerm_postgresql_server.test.administrator_login_password

  sku_name   = azurerm_postgresql_server.test.sku_name
  version    = azurerm_postgresql_server.test.version
  storage_mb = azurerm_postgresql_server.test.storage_mb

  ssl_enforcement_enabled = azurerm_postgresql_server.test.ssl_enforcement_enabled
}
`, r.basic(data, "10.0"))
}

func (PostgreSQLServerResource) completeDeprecated(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  version  = "%s"
  sku_name = "GP_Gen5_2"

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"

  infrastructure_encryption_enabled = true
  public_network_access_enabled     = false
  ssl_minimal_tls_version_enforced  = "TLS1_2"

  ssl_enforcement_enabled = true

  storage_profile {
    storage_mb            = 640000
    backup_retention_days = 7
    geo_redundant_backup  = "Enabled"
    auto_grow             = "Enabled"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func (PostgreSQLServerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!updated"

  sku_name   = "GP_Gen5_4"
  version    = "9.6"
  storage_mb = 640000

  backup_retention_days        = 7
  geo_redundant_backup_enabled = true
  auto_grow_enabled            = true

  infrastructure_encryption_enabled = true
  public_network_access_enabled     = false
  ssl_enforcement_enabled           = true
  ssl_minimal_tls_version_enforced  = "TLS1_2"

  threat_detection_policy {
    enabled              = true
    disabled_alerts      = ["Sql_Injection", "Data_Exfiltration"]
    email_account_admins = true
    email_addresses      = ["kt@example.com", "admin@example.com"]

    retention_days = 7
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (PostgreSQLServerResource) complete2(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!updated"

  sku_name   = "GP_Gen5_4"
  version    = "%[3]s"
  storage_mb = 640000

  backup_retention_days        = 14
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = false

  infrastructure_encryption_enabled = false
  public_network_access_enabled     = true
  ssl_enforcement_enabled           = false

  threat_detection_policy {
    enabled              = true
    disabled_alerts      = ["Sql_Injection"]
    email_account_admins = true
    email_addresses      = ["kt@example.com"]

    retention_days = 7
  }
}
`, data.RandomInteger, data.Locations.Primary, version)
}

func (PostgreSQLServerResource) sku(data acceptance.TestData, version, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "%s"
  storage_mb = 51200
  version    = "%s"

  ssl_enforcement_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sku, version)
}

func (r PostgreSQLServerResource) createReplica(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "replica" {
  name     = "acctestRG-psql-%[2]d-replica"
  location = "%[3]s"
}

resource "azurerm_postgresql_server" "replica" {
  name                = "acctest-psql-server-%[2]d-replica"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.replica.name

  sku_name = "%[4]s"
  version  = "11"

  create_mode               = "Replica"
  creation_source_server_id = azurerm_postgresql_server.test.id

  public_network_access_enabled = false
  ssl_enforcement_enabled       = true
}
`, r.template(data, sku, "11"), data.RandomInteger, data.Locations.Secondary, sku)
}

func (r PostgreSQLServerResource) updateReplicaToDefault(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "replica" {
  name     = "acctestRG-psql-%[2]d-replica"
  location = "%[3]s"
}

resource "azurerm_postgresql_server" "replica" {
  name                = "acctest-psql-server-%[2]d-replica"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.replica.name

  sku_name    = "%[4]s"
  version     = "11"
  create_mode = "Default"

  public_network_access_enabled = false
  ssl_enforcement_enabled       = true
}
`, r.template(data, sku, "11"), data.RandomInteger, data.Locations.Secondary, sku)
}

func (r PostgreSQLServerResource) createReplicas(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "replica1" {
  name     = "acctestRG-psql-%[2]d-replica1"
  location = "%[3]s"
}

resource "azurerm_postgresql_server" "replica1" {
  name                = "acctest-psql-server-%[2]d-replica1"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.replica1.name

  sku_name = "%[4]s"
  version  = "11"

  create_mode               = "Replica"
  creation_source_server_id = azurerm_postgresql_server.test.id

  ssl_enforcement_enabled = true
}

resource "azurerm_postgresql_server" "replica2" {
  name                = "acctest-psql-server-%[2]d-replica2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "%[4]s"
  version  = "11"

  create_mode               = "Replica"
  creation_source_server_id = azurerm_postgresql_server.test.id

  ssl_enforcement_enabled = true
}
`, r.template(data, sku, "11"), data.RandomInteger, data.Locations.Secondary, sku)
}

func (r PostgreSQLServerResource) createPointInTimeRestore(data acceptance.TestData, version, restoreTime string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_postgresql_server" "restore" {
  name                = "acctest-psql-server-%[2]d-restore"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name   = "GP_Gen5_2"
  version    = "%[4]s"
  storage_mb = 51200

  create_mode               = "PointInTimeRestore"
  creation_source_server_id = azurerm_postgresql_server.test.id
  restore_point_in_time     = "%[3]s"

  ssl_enforcement_enabled = true
}
`, r.basic(data, version), data.RandomInteger, restoreTime, version)
}

func (PostgreSQLServerResource) emptyAttrs(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!updated"

  sku_name   = "GP_Gen5_4"
  version    = "%[3]s"
  storage_mb = 640000

  ssl_enforcement_enabled          = false
  ssl_minimal_tls_version_enforced = "TLSEnforcementDisabled"

  threat_detection_policy {
    enabled              = true
    email_account_admins = true

    retention_days = 7
  }
}
`, data.RandomInteger, data.Locations.Primary, version)
}

func (PostgreSQLServerResource) beforeUpdate(data acceptance.TestData, version string, tlsVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!updated"

  sku_name   = "GP_Gen5_4"
  version    = "%[3]s"
  storage_mb = 640000

  backup_retention_days = 7
  auto_grow_enabled     = true

  public_network_access_enabled    = false
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "%[4]s"
}
`, data.RandomInteger, data.Locations.Primary, version, tlsVersion)
}

func (PostgreSQLServerResource) afterUpdate(data acceptance.TestData, version string, tlsVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!updated"

  sku_name   = "GP_Gen5_4"
  version    = "%[3]s"
  storage_mb = 640000

  backup_retention_days = 7
  auto_grow_enabled     = true

  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "%[4]s"

}
`, data.RandomInteger, data.Locations.Primary, version, tlsVersion)
}
