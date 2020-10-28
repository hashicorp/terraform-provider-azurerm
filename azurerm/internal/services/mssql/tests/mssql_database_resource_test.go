package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMsSqlDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMMsSqlDatabase_requiresImport),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "SQL_AltDiction_CP850_CI_AI"),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", "BasePrice"),
					resource.TestCheckResourceAttr(data.ResourceName, "max_size_gb", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_Gen5_2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "Test"),
				),
			},
			data.ImportStep("sample_name"),
			{
				Config: testAccAzureRMMsSqlDatabase_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", "LicenseIncluded"),
					resource.TestCheckResourceAttr(data.ResourceName, "max_size_gb", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "Staging"),
				),
			},
			data.ImportStep("sample_name"),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_elasticPool(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlDatabase_elasticPool(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "elastic_pool_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "ElasticPool"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlDatabase_elasticPoolDisassociation(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_GP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_GP(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_Gen5_2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_GP_Serverless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_GPServerless(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_pause_delay_in_minutes", "70"),
					resource.TestCheckResourceAttr(data.ResourceName, "min_capacity", "0.75"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_S_Gen5_2"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlDatabase_GPServerlessUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_pause_delay_in_minutes", "90"),
					resource.TestCheckResourceAttr(data.ResourceName, "min_capacity", "1.25"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_S_Gen5_2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_BC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_BC(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "read_scale", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "BC_Gen5_2"),
					resource.TestCheckResourceAttr(data.ResourceName, "zone_redundant", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlDatabase_BCUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "read_scale", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "BC_Gen5_2"),
					resource.TestCheckResourceAttr(data.ResourceName, "zone_redundant", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_HS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_HS(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "read_replica_count", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "HS_Gen5_2"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlDatabase_HSUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "read_replica_count", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "HS_Gen5_2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_createCopyMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "copy")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_createCopyMode(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "SQL_AltDiction_CP850_CI_AI"),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", "BasePrice"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_Gen5_2"),
				),
			},
			data.ImportStep("create_mode", "creation_source_database_id"),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_createPITRMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),

			{
				PreConfig: func() { time.Sleep(7 * time.Minute) },
				Config:    testAccAzureRMMsSqlDatabase_createPITRMode(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists("azurerm_mssql_database.pitr"),
				),
			},

			data.ImportStep("create_mode", "creation_source_database_id", "restore_point_in_time"),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_createSecondaryMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "secondary")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_createSecondaryMode(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "SQL_AltDiction_CP850_CI_AI"),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", "BasePrice"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_Gen5_2"),
				),
			},
			data.ImportStep("create_mode", "creation_source_database_id", "sample_name"),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_createRestoreMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_createRestoreMode(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep("create_mode", "creation_source_database_id"),

			{
				PreConfig: func() { time.Sleep(8 * time.Minute) },
				Config:    testAccAzureRMMsSqlDatabase_createRestoreModeDBDeleted(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},

			data.ImportStep(),

			{
				PreConfig: func() { time.Sleep(8 * time.Minute) },
				Config:    testAccAzureRMMsSqlDatabase_createRestoreModeDBRestored(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					testCheckAzureRMMsSqlDatabaseExists("azurerm_mssql_database.restore"),
				),
			},

			data.ImportStep("create_mode", "restore_dropped_database_id"),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_threatDetectionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_threatDetectionPolicy(data, "Enabled"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "threat_detection_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "threat_detection_policy.0.state", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "threat_detection_policy.0.retention_days", "15"),
					resource.TestCheckResourceAttr(data.ResourceName, "threat_detection_policy.0.disabled_alerts.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "threat_detection_policy.0.email_account_admins", "Enabled"),
				),
			},
			data.ImportStep("sample_name", "threat_detection_policy.0.storage_account_access_key"),
			{
				Config: testAccAzureRMMsSqlDatabase_threatDetectionPolicy(data, "Disabled"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "threat_detection_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "threat_detection_policy.0.state", "Disabled"),
				),
			},
			data.ImportStep("sample_name", "threat_detection_policy.0.storage_account_access_key"),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_withBlobAuditingPolices(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlDatabase_withBlobAuditingPolices(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep("extended_auditing_policy.0.storage_account_access_key"),
			{
				Config: testAccAzureRMMsSqlDatabase_withBlobAuditingPolicesUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep("extended_auditing_policy.0.storage_account_access_key"),
			{
				Config: testAccAzureRMMsSqlDatabase_withBlobAuditingPolicesDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_updateSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlDatabase_updateSku(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlDatabase_updateSku2(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_minCapacity0(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_minCapacity0(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_withLongTermRetentionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlDatabase_withLongTermRetentionPolicy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlDatabase_withLongTermRetentionPolicyUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlDatabase_withShortTermRetentionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlDatabase_withShortTermRetentionPolicy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlDatabase_withShortTermRetentionPolicyUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
		},
	})
}

func testCheckAzureRMMsSqlDatabaseExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.DatabasesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.MsSqlDatabaseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.MsSqlServer, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("MsSql Database %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
			}

			return fmt.Errorf("Get on MsSql Database Client: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMsSqlDatabaseDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.DatabasesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mssql_database" {
			continue
		}

		id, err := parse.MsSqlDatabaseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.MsSqlServer, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Get on MsSql Database Client: %+v", err)
			}
		}
		return nil
	}

	return nil
}

func testAccAzureRMMsSqlDatabase_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest-sqlserver-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMMsSqlDatabase_basic(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%d"
  server_id = azurerm_sql_server.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "import" {
  name      = azurerm_mssql_database.test.name
  server_id = azurerm_sql_server.test.id
}
`, template)
}

func testAccAzureRMMsSqlDatabase_complete(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name         = "acctest-db-%[2]d"
  server_id    = azurerm_sql_server.test.id
  collation    = "SQL_AltDiction_CP850_CI_AI"
  license_type = "BasePrice"
  max_size_gb  = 1
  sample_name  = "AdventureWorksLT"
  sku_name     = "GP_Gen5_2"

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_update(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name         = "acctest-db-%[2]d"
  server_id    = azurerm_sql_server.test.id
  collation    = "SQL_AltDiction_CP850_CI_AI"
  license_type = "LicenseIncluded"
  max_size_gb  = 2
  sku_name     = "GP_Gen5_2"

  tags = {
    ENV = "Staging"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_elasticPool(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  max_size_gb         = 5

  sku {
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
    capacity = 4
    family   = "Gen5"
  }

  per_database_settings {
    min_capacity = 0.25
    max_capacity = 4
  }
}

resource "azurerm_mssql_database" "test" {
  name            = "acctest-db-%[2]d"
  server_id       = azurerm_sql_server.test.id
  elastic_pool_id = azurerm_mssql_elasticpool.test.id
  sku_name        = "ElasticPool"
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_elasticPoolDisassociation(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  max_size_gb         = 5

  sku {
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
    capacity = 4
    family   = "Gen5"
  }

  per_database_settings {
    min_capacity = 0.25
    max_capacity = 4
  }
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_sql_server.test.id
  sku_name  = "GP_Gen5_2"
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_GP(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%d"
  server_id = azurerm_sql_server.test.id
  sku_name  = "GP_Gen5_2"
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_GPServerless(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name                        = "acctest-db-%d"
  server_id                   = azurerm_sql_server.test.id
  auto_pause_delay_in_minutes = 70
  min_capacity                = 0.75
  sku_name                    = "GP_S_Gen5_2"
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_GPServerlessUpdate(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name                        = "acctest-db-%d"
  server_id                   = azurerm_sql_server.test.id
  auto_pause_delay_in_minutes = 90
  min_capacity                = 1.25
  sku_name                    = "GP_S_Gen5_2"
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_HS(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name               = "acctest-db-%d"
  server_id          = azurerm_sql_server.test.id
  read_replica_count = 2
  sku_name           = "HS_Gen5_2"
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_HSUpdate(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name               = "acctest-db-%d"
  server_id          = azurerm_sql_server.test.id
  read_replica_count = 4
  sku_name           = "HS_Gen5_2"
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_BC(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name           = "acctest-db-%d"
  server_id      = azurerm_sql_server.test.id
  read_scale     = true
  sku_name       = "BC_Gen5_2"
  zone_redundant = true
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_BCUpdate(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name           = "acctest-db-%d"
  server_id      = azurerm_sql_server.test.id
  read_scale     = false
  sku_name       = "BC_Gen5_2"
  zone_redundant = false
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_createCopyMode(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_complete(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "copy" {
  name                        = "acctest-dbc-%d"
  server_id                   = azurerm_sql_server.test.id
  create_mode                 = "Copy"
  creation_source_database_id = azurerm_mssql_database.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_createPITRMode(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "pitr" {
  name                        = "acctest-dbp-%d"
  server_id                   = azurerm_sql_server.test.id
  create_mode                 = "PointInTimeRestore"
  restore_point_in_time       = "%s"
  creation_source_database_id = azurerm_mssql_database.test.id

}
`, template, data.RandomInteger, time.Now().Add(time.Duration(7)*time.Minute).UTC().Format(time.RFC3339))
}

func testAccAzureRMMsSqlDatabase_createSecondaryMode(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_complete(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "second" {
  name     = "acctestRG-mssql2-%[2]d"
  location = "%[3]s"
}

resource "azurerm_sql_server" "second" {
  name                         = "acctest-sqlserver2-%[2]d"
  resource_group_name          = azurerm_resource_group.second.name
  location                     = azurerm_resource_group.second.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "secondary" {
  name                        = "acctest-dbs-%[2]d"
  server_id                   = azurerm_sql_server.second.id
  create_mode                 = "Secondary"
  creation_source_database_id = azurerm_mssql_database.test.id

}
`, template, data.RandomInteger, data.Locations.Secondary)
}

func testAccAzureRMMsSqlDatabase_createRestoreMode(data acceptance.TestData) string {
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
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}


resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[1]d"
  server_id = azurerm_mssql_server.test.id
}

resource "azurerm_mssql_database" "copy" {
  name                        = "acctest-dbc-%[1]d"
  server_id                   = azurerm_mssql_server.test.id
  create_mode                 = "Copy"
  creation_source_database_id = azurerm_mssql_database.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMMsSqlDatabase_createRestoreModeDBDeleted(data acceptance.TestData) string {
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
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}


resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[1]d"
  server_id = azurerm_mssql_server.test.id
}

`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMMsSqlDatabase_createRestoreModeDBRestored(data acceptance.TestData) string {
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
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}


resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[1]d"
  server_id = azurerm_mssql_server.test.id
}

resource "azurerm_mssql_database" "restore" {
  name                        = "acctest-dbr-%[1]d"
  server_id                   = azurerm_mssql_server.test.id
  create_mode                 = "Restore"
  restore_dropped_database_id = azurerm_mssql_server.test.restorable_dropped_database_ids[0]
}

`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMMsSqlDatabase_threatDetectionPolicy(data acceptance.TestData, state string) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "test%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_mssql_database" "test" {
  name         = "acctest-db-%[2]d"
  server_id    = azurerm_sql_server.test.id
  collation    = "SQL_AltDiction_CP850_CI_AI"
  license_type = "BasePrice"
  max_size_gb  = 1
  sample_name  = "AdventureWorksLT"
  sku_name     = "GP_Gen5_2"

  threat_detection_policy {
    retention_days             = 15
    state                      = "%[3]s"
    disabled_alerts            = ["Sql_Injection"]
    email_account_admins       = "Enabled"
    storage_account_access_key = azurerm_storage_account.test.primary_access_key
    storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
    use_server_default         = "Disabled"
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, state)
}

func testAccAzureRMMsSqlDatabase_withBlobAuditingPolices(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_sql_server.test.id
  extended_auditing_policy {
    storage_endpoint                        = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_access_key              = azurerm_storage_account.test.primary_access_key
    storage_account_access_key_is_secondary = true
    retention_in_days                       = 6
  }
}
`, template, data.RandomIntOfLength(15), data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_withBlobAuditingPolicesUpdated(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_sql_server.test.id
  extended_auditing_policy {
    storage_endpoint                        = azurerm_storage_account.test2.primary_blob_endpoint
    storage_account_access_key              = azurerm_storage_account.test2.primary_access_key
    storage_account_access_key_is_secondary = false
    retention_in_days                       = 3
  }
}
`, template, data.RandomIntOfLength(15), data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_withBlobAuditingPolicesDisabled(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name                     = "acctest-db-%[3]d"
  server_id                = azurerm_sql_server.test.id
  extended_auditing_policy = []
}
`, template, data.RandomIntOfLength(15), data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_updateSku(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%d"
  server_id = azurerm_sql_server.test.id
  sku_name  = "HS_Gen5_2"
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_updateSku2(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%d"
  server_id = azurerm_sql_server.test.id
  sku_name  = "HS_Gen5_4"
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_minCapacity0(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%d"
  server_id = azurerm_sql_server.test.id

  min_capacity = 0
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_withLongTermRetentionPolicy(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_sql_server.test.id
  long_term_retention_policy {
    weekly_retention  = "P1W"
    monthly_retention = "P1M"
    yearly_retention  = "P1Y"
    week_of_year      = 1
  }
}
`, template, data.RandomIntOfLength(15), data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_withLongTermRetentionPolicyUpdated(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_sql_server.test.id
  long_term_retention_policy {
    weekly_retention = "P1W"
    yearly_retention = "P1Y"
    week_of_year     = 2
  }
}
`, template, data.RandomIntOfLength(15), data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_withShortTermRetentionPolicy(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_sql_server.test.id
  short_term_retention_policy {
    retention_days = 8
  }
}
`, template, data.RandomIntOfLength(15), data.RandomInteger)
}

func testAccAzureRMMsSqlDatabase_withShortTermRetentionPolicyUpdated(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_sql_server.test.id
  short_term_retention_policy {
    retention_days = 10
  }
}
`, template, data.RandomIntOfLength(15), data.RandomInteger)
}
