package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMSqlElasticPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlElasticPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSqlElasticPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlElasticPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPoolExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSqlElasticPool_requiresImport),
		},
	})
}

func TestAccAzureRMSqlElasticPool_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlElasticPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPoolExists(data.ResourceName),
					testCheckAzureRMSqlElasticPoolDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMSqlElasticPool_resizeDtu(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlElasticPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "dtu", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "pool_size", "5000"),
				),
			},
			{
				Config: testAccAzureRMSqlElasticPool_resizedDtu(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "dtu", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "pool_size", "10000"),
				),
			},
		},
	})
}

func testCheckAzureRMSqlElasticPoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Sql.ElasticPoolsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		poolName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, serverName, poolName)
		if err != nil {
			return fmt.Errorf("Bad: Get on sqlElasticPoolsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: SQL Elastic Pool %q on server: %q (resource group: %q) does not exist", poolName, serverName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMSqlElasticPoolDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Sql.ElasticPoolsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sql_elasticpool" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		poolName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, serverName, poolName)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("SQL Elastic Pool still exists:\n%#v", resp.ElasticPoolProperties)
		}
	}

	return nil
}

func testCheckAzureRMSqlElasticPoolDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Sql.ElasticPoolsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		poolName := rs.Primary.Attributes["name"]

		if _, err := client.Delete(ctx, resourceGroup, serverName, poolName); err != nil {
			return fmt.Errorf("Bad: Delete on sqlElasticPoolsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMSqlElasticPool_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%[1]d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_sql_elasticpool" "test" {
  name                = "acctest-pool-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  edition             = "Basic"
  dtu                 = 50
  pool_size           = 5000
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMSqlElasticPool_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sql_elasticpool" "import" {
  name                = azurerm_sql_elasticpool.test.name
  resource_group_name = azurerm_sql_elasticpool.test.resource_group_name
  location            = azurerm_sql_elasticpool.test.location
  server_name         = azurerm_sql_elasticpool.test.server_name
  edition             = azurerm_sql_elasticpool.test.edition
  dtu                 = azurerm_sql_elasticpool.test.dtu
  pool_size           = azurerm_sql_elasticpool.test.pool_size
}
`, testAccAzureRMSqlElasticPool_basic(data))
}

func testAccAzureRMSqlElasticPool_resizedDtu(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%[1]d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_sql_elasticpool" "test" {
  name                = "acctest-pool-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  edition             = "Basic"
  dtu                 = 100
  pool_size           = 10000
}
`, data.RandomInteger, data.Locations.Primary)
}
