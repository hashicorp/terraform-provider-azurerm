package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPostgreSQLDatabase_basic(t *testing.T) {
	resourceName := "azurerm_postgresql_database.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLDatabase_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "charset", "UTF8"),
					resource.TestCheckResourceAttr(resourceName, "collation", "English_United States.1252"),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLDatabase_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_postgresql_database.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLDatabase_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "charset", "UTF8"),
					resource.TestCheckResourceAttr(resourceName, "collation", "English_United States.1252"),
				),
			},
			{
				Config:      testAccAzureRMPostgreSQLDatabase_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_postgresql_database"),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLDatabase_collationWithHyphen(t *testing.T) {
	resourceName := "azurerm_postgresql_database.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLDatabase_collationWithHyphen(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "charset", "UTF8"),
					resource.TestCheckResourceAttr(resourceName, "collation", "En-US"),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLDatabase_charsetLowercase(t *testing.T) {
	resourceName := "azurerm_postgresql_database.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLDatabase_charsetLowercase(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "charset", "UTF8"),
					resource.TestCheckResourceAttr(resourceName, "collation", "English_United States.1252"),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLDatabase_charsetMixedcase(t *testing.T) {
	resourceName := "azurerm_postgresql_database.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLDatabase_charsetMixedcase(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "charset", "UTF8"),
					resource.TestCheckResourceAttr(resourceName, "collation", "English_United States.1252"),
				),
			},
		},
	})
}

func testCheckAzureRMPostgreSQLDatabaseExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for PostgreSQL Database: %s", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.DatabasesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: PostgreSQL Database %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on postgresqlDatabasesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLDatabaseDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.DatabasesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_postgresql_database" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, serverName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("PostgreSQL Database still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMPostgreSQLDatabase_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpsqlsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "GP_Gen5_2"
    capacity = 2
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement              = "Enabled"
}

resource "azurerm_postgresql_database" "test" {
  name                = "acctestdb_%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_postgresql_server.test.name}"
  charset             = "UTF8"
  collation           = "English_United States.1252"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPostgreSQLDatabase_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_database" "import" {
  name                = "${azurerm_postgresql_database.test.name}"
  resource_group_name = "${azurerm_postgresql_database.test.resource_group_name}"
  server_name         = "${azurerm_postgresql_database.test.server_name}"
  charset             = "${azurerm_postgresql_database.test.charset}"
  collation           = "${azurerm_postgresql_database.test.collation}"
}
`, testAccAzureRMPostgreSQLDatabase_basic(rInt, location))
}

func testAccAzureRMPostgreSQLDatabase_collationWithHyphen(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpsqlsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "GP_Gen5_2"
    capacity = 2
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement              = "Enabled"
}

resource "azurerm_postgresql_database" "test" {
  name                = "acctestdb_%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_postgresql_server.test.name}"
  charset             = "UTF8"
  collation           = "En-US"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPostgreSQLDatabase_charsetLowercase(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpsqlsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "GP_Gen5_2"
    capacity = 2
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement              = "Enabled"
}

resource "azurerm_postgresql_database" "test" {
  name                = "acctestdb_%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_postgresql_server.test.name}"
  charset             = "utf8"
  collation           = "English_United States.1252"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPostgreSQLDatabase_charsetMixedcase(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpsqlsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "GP_Gen5_2"
    capacity = 2
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement              = "Enabled"
}

resource "azurerm_postgresql_database" "test" {
  name                = "acctestdb_%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_postgresql_server.test.name}"
  charset             = "Utf8"
  collation           = "English_United States.1252"
}
`, rInt, location, rInt, rInt)
}
