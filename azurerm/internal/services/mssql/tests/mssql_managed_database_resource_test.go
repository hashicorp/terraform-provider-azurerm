package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMsSqlManagedDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedDatabase_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedDatabase_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedDatabaseExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMMssqlManagedDatabase_requiresImport),
		},
	})
}

func TestAccAzureRMMsSqlManagedDatabase_updateCollation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedDatabase_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "SQL_Latin1_General_CP1_CI_AS"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlManagedDatabase_UpdateCollation(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "Estonian_100_CS_AS_SC_UTF8"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedDatabase_updateCatalogCollation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedDatabase_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "catalog_collation", "SQL_Latin1_General_CP1_CI_AS"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlManagedDatabase_UpdateCatalogCollation(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "catalog_collation", "DATABASE_DEFAULT"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedDatabase_createPITRMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedDatabase_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),

			{
				PreConfig: func() { time.Sleep(7 * time.Minute) },
				Config:    testAccAzureRMMsSqlManagedDatabase_createPITRMode(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedDatabaseExists("azurerm_mssql_managed_database.pitr"),
				),
			},

			data.ImportStep("create_mode", "source_database_id", "restore_point_in_time"),
		},
	})
}

func testAccAzureRMMsSqlManagedDatabase_basicTemplate(data acceptance.TestData) string {
	return testAccAzureRMMsSqlManagedDatabase_basic(data, "SQL_Latin1_General_CP1_CI_AS", "SQL_Latin1_General_CP1_CI_AS")
}

func testAccAzureRMMsSqlManagedDatabase_UpdateCollation(data acceptance.TestData) string {
	return testAccAzureRMMsSqlManagedDatabase_basic(data, "Estonian_100_CS_AS_SC_UTF8", "SQL_Latin1_General_CP1_CI_AS")
}

func testAccAzureRMMsSqlManagedDatabase_UpdateCatalogCollation(data acceptance.TestData) string {
	return testAccAzureRMMsSqlManagedDatabase_basic(data, "SQL_Latin1_General_CP1_CI_AS", "DATABASE_DEFAULT")
}

func testCheckAzureRMMsSqlManagedDatabaseExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ManagedDatabasesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		managedInstanceId := rs.Primary.Attributes["managed_instance_id"]

		id, err := azure.ParseAzureResourceID(managedInstanceId)
		if err != nil {
			return err
		}

		managedInstanceName := id.Path["managedInstances"]
		resourceGroup := id.ResourceGroup

		resp, err := client.Get(ctx, resourceGroup, managedInstanceName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on ManagedDatabase: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Managed database %q  (Managed Sql Instance %q, resource group: %q) does not exist", name, managedInstanceName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMMsSqlManagedDatabaseDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ManagedDatabasesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mssql_managed_database" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		managedInstanceId := rs.Primary.Attributes["managed_instance_id"]
		id, err := azure.ParseAzureResourceID(managedInstanceId)
		if err != nil {
			return err
		}
		managedInstanceName := id.Path["managedInstances"]
		resourceGroup := id.ResourceGroup

		if resp, err := client.Get(ctx, resourceGroup, managedInstanceName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Get on MsSql Database Client: %+v", err)
			}
		}
		return nil
	}

	return nil
}

func testAccAzureRMMsSqlManagedDatabase_basic(data acceptance.TestData, collation string, catalogCollation string) string {
	template := testAccAzureRMMsSqlManagedDatabase_prepareDependencies(data)
	return fmt.Sprintf(`%s

resource "azurerm_mssql_managed_database" "test" {
	name                         = "acctest-db-%[2]d"
	managed_instance_name          = azurerm_mssql_managed_instance.test.name
	resource_group_name				= azurerm_mssql_managed_instance.test.resource_group_name
	collation					 = "%[3]s"
	catalog_collation           =  "%[4]s"
  }
`, template, data.RandomInteger, collation, catalogCollation)
}

func testAccAzureRMMsSqlManagedDatabase_createPITRMode(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlManagedDatabase_prepareDependencies(data)
	return fmt.Sprintf(`%s

resource "azurerm_mssql_managed_database" "pitr" {
	name                         = "acctest-dbp-%d"
	managed_instance_name          = azurerm_mssql_managed_instance.test.name
	resource_group_name				= azurerm_mssql_managed_instance.test.resource_group_name
	create_mode 				= "PointInTimeRestore"
  	restore_point_in_time       = "%s"
  	source_database_id          =  azurerm_mssql_managed_database.test.id
}
`, template, data.RandomInteger, time.Now().Add(time.Duration(7)*time.Minute).UTC().Format(time.RFC3339))
}

func testAccAzureRMMssqlManagedDatabase_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlManagedDatabase_basic(data, "SQL_Latin1_General_CP1_CI_AS", "SQL_Latin1_General_CP1_CI_AS")
	return fmt.Sprintf(`%s

resource "azurerm_mssql_managed_database" "import" {
	name                         = "acctest-db-%[2]d"
	managed_instance_name          = azurerm_mssql_managed_instance.test.name
	resource_group_name				= azurerm_mssql_managed_instance.test.resource_group_name
	collation					 = "%[3]s"
  }
`, template, data.RandomInteger, "SQL_Latin1_General_CP1_CI_AS", "Default", "SQL_Latin1_General_CP1_CI_AS")
}

func testAccAzureRMMsSqlManagedDatabase_prepareDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_network_security_group" "test" {
  name                = "accTestNetworkSecurityGroup-%[1]d"
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-%[1]d-network"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  virtual_network_name = azurerm_virtual_network.test.name
  resource_group_name  = azurerm_resource_group.test.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "miDelegation"
    service_delegation {
      name = "Microsoft.Sql/managedInstances"
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_route_table" "test" {
  name                = "test-routetable-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  route {
    name                   = "test"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}

resource "azurerm_mssql_managed_instance" "test" {
  name                         = "acctest-mi-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "AcceptanceTestUser"
  administrator_login_password = "LengthyPassword@1234"
  subnet_id                    = azurerm_subnet.test.id
  identity {
    type = "SystemAssigned"
  }
  sku {
    capacity = 8
    family   = "Gen5"
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
  }
  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]
}
	`, data.RandomInteger, data.Locations.Primary)
}
