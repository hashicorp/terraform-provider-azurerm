package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMsSqlManagedInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstance_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstance_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMMsSqlManagedInstance_requiresImport),
		},
	})
}

func TestAccAzureRMMsSqlManagedInstance_updateAdminCredentials(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstance_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "miadmin"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlManagedInstance_updateAdminCredentials(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "miadmin2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedInstance_updateSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstance_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "8"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.family", "Gen5"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "GP_Gen5"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "GeneralPurpose"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlManagedInstance_updateSku(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "24"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.family", "Gen5"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "BC_Gen5"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "BusinessCritical"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedInstance_updateLicence(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstance_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", "BasePrice"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlManagedInstance_updateLicence(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", "LicenseIncluded"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedInstance_updateCollation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstance_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "SQL_Latin1_General_CP1_CI_AS"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlManagedInstance_updateCollation(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "Latin1_General_100_CS_AS_SC"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedInstance_updateProxyOverride(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstance_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "proxy_override", "Redirect"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlManagedInstance_updateProxyOverride(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "proxy_override", "Proxy"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedInstance_updateStorageSize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstance_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_size_gb", "64"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlManagedInstance_updateStorageSize(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_size_gb", "128"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedInstance_updateVcores(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstance_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vcores", "8"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlManagedInstance_updateVcores(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vcores", "4"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedInstance_updateTimeZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstance_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "timezone_id", "UTC"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlManagedInstance_updateTimeZone(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "timezone_id", "India Standard Time"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedInstance_updateTlsVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstance_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "minimal_tls_version", "1.1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlManagedInstance_updateTlsVersion(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "minimal_tls_version", "1.2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedInstance_DnsZonePartner(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")
	dns := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "dnspartner")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstance_DnsZonePartner(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					testCheckAzureRMMsSqlManagedInstanceExists(dns.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMMsSqlManagedInstance_basicTemplate(data acceptance.TestData) string {
	return testAccAzureRMMsSqlManagedInstance_basic(data, "miadmin", "LengthyPassword@1234", 8, "GP_Gen5", "GeneralPurpose", "BasePrice", "SQL_Latin1_General_CP1_CI_AS", "Redirect", 64, 8, "UTC", "1.1")
}

func testAccAzureRMMsSqlManagedInstance_updateAdminCredentials(data acceptance.TestData) string {
	return testAccAzureRMMsSqlManagedInstance_basic(data, "miadmin2", "LengthyPassword@4321", 16, "GP_Gen5", "GeneralPurpose", "BasePrice", "SQL_Latin1_General_CP1_CI_AS", "Redirect", 64, 8, "UTC", "1.1")
}

func testAccAzureRMMsSqlManagedInstance_updateSku(data acceptance.TestData) string {
	return testAccAzureRMMsSqlManagedInstance_basic(data, "miadmin", "LengthyPassword@1234", 24, "BC_Gen5", "BusinessCritical", "BasePrice", "SQL_Latin1_General_CP1_CI_AS", "Redirect", 64, 8, "UTC", "1.1")
}

func testAccAzureRMMsSqlManagedInstance_updateLicence(data acceptance.TestData) string {
	return testAccAzureRMMsSqlManagedInstance_basic(data, "miadmin", "LengthyPassword@1234", 32, "GP_Gen5", "GeneralPurpose", "LicenseIncluded", "SQL_Latin1_General_CP1_CI_AS", "Redirect", 64, 8, "UTC", "1.1")
}

func testAccAzureRMMsSqlManagedInstance_updateCollation(data acceptance.TestData) string {
	return testAccAzureRMMsSqlManagedInstance_basic(data, "miadmin", "LengthyPassword@1234", 8, "GP_Gen5", "GeneralPurpose", "BasePrice", "Latin1_General_100_CS_AS_SC", "Redirect", 64, 8, "UTC", "1.1")
}

func testAccAzureRMMsSqlManagedInstance_updateProxyOverride(data acceptance.TestData) string {
	return testAccAzureRMMsSqlManagedInstance_basic(data, "miadmin", "LengthyPassword@1234", 16, "GP_Gen5", "GeneralPurpose", "BasePrice", "SQL_Latin1_General_CP1_CI_AS", "Proxy", 64, 8, "UTC", "1.1")
}

func testAccAzureRMMsSqlManagedInstance_updateStorageSize(data acceptance.TestData) string {
	return testAccAzureRMMsSqlManagedInstance_basic(data, "miadmin", "LengthyPassword@1234", 32, "GP_Gen5", "GeneralPurpose", "BasePrice", "SQL_Latin1_General_CP1_CI_AS", "Redirect", 128, 8, "UTC", "1.1")
}

func testAccAzureRMMsSqlManagedInstance_updateVcores(data acceptance.TestData) string {
	return testAccAzureRMMsSqlManagedInstance_basic(data, "miadmin", "LengthyPassword@1234", 24, "GP_Gen5", "GeneralPurpose", "BasePrice", "SQL_Latin1_General_CP1_CI_AS", "Redirect", 64, 4, "UTC", "1.1")
}

func testAccAzureRMMsSqlManagedInstance_updateTimeZone(data acceptance.TestData) string {
	return testAccAzureRMMsSqlManagedInstance_basic(data, "miadmin", "LengthyPassword@1234", 24, "GP_Gen5", "GeneralPurpose", "BasePrice", "SQL_Latin1_General_CP1_CI_AS", "Redirect", 64, 8, "India Standard Time", "1.1")
}

func testAccAzureRMMsSqlManagedInstance_updateTlsVersion(data acceptance.TestData) string {
	return testAccAzureRMMsSqlManagedInstance_basic(data, "miadmin", "LengthyPassword@1234", 8, "GP_Gen5", "GeneralPurpose", "BasePrice", "SQL_Latin1_General_CP1_CI_AS", "Redirect", 64, 8, "UTC", "1.2")
}

func testAccAzureRMMsSqlManagedInstance_basic(data acceptance.TestData, adminLogin string, adminPassword string, skuCapacity int, skuName string, skuTier string, licenseType string, collation string, proxyOverride string, storageSize int, vcores int, timeZoneId string, minTlsVersion string) string {
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
  administrator_login          = "%[3]s"
  administrator_login_password = "%[4]s"
  subnet_id                    = azurerm_subnet.test.id
  identity {
    type = "SystemAssigned"
  }
  sku {
    capacity = %[5]d
    family   = "Gen5"
    name     = "%[6]s"
    tier     = "%[7]s"
  }
  license_type          = "%[8]s"
  collation             = "%[9]s"
  proxy_override        = "%[10]s"
  storage_size_gb       = %[11]d
  vcores                = %[12]d
  data_endpoint_enabled = true
  timezone_id           = "%[13]s"
  minimal_tls_version   = "%[14]s"
  depends_on = [
	azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, adminLogin, adminPassword, skuCapacity, skuName, skuTier, licenseType, collation, proxyOverride, storageSize, vcores, timeZoneId, minTlsVersion)
}

func testAccAzureRMMsSqlManagedInstance_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlManagedInstance_basicTemplate(data)
	return fmt.Sprintf(`%s

resource "azurerm_mssql_managed_instance" "import" {
  name                         = "acctest-mi-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "miadministrator"
  administrator_login_password = "Password@1234"
  subnet_id                    = azurerm_subnet.test.id
  identity {
    type = "SystemAssigned"
  }
  sku {
    capacity = 16
    family   = "Gen5"
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
  }
  license_type          = "LicenseIncluded"
  collation             = "SQL_Latin1_General_CP1_CI_AS"
  proxy_override        = "Redirect"
  storage_size_gb       = 256
  vcores                = 16
  data_endpoint_enabled = true
  timezone_id           = "UTC"
  minimal_tls_version   = "1.1"
}
`, template, data.RandomInteger)
}

func testAccAzureRMMsSqlManagedInstance_DnsZonePartner(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlManagedInstance_basicTemplate(data)
	return fmt.Sprintf(`%s

resource "azurerm_mssql_managed_instance" "dnspartner" {
  name                         = "acctest-mi-dns-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "miadministrator2"
  administrator_login_password = "LengthyPassword@4321"
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
  license_type          = "BasePrice"
  collation             = "SQL_Latin1_General_CP1_CI_AS"
  proxy_override        = "Redirect"
  storage_size_gb       = 64
  vcores                = 8
  data_endpoint_enabled = true
  timezone_id           = "UTC"
  dns_zone_partner      = azurerm_mssql_managed_instance.test.id
  minimal_tls_version   = "1.1"
}
`, template, data.RandomInteger)
}

func testCheckAzureRMMsSqlManagedInstanceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ManagedInstancesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		miName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, miName)
		if err != nil {
			return fmt.Errorf("Bad: Get on ManagedInstancesClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Managed Sql Instance %q (resource group: %q) does not exist", miName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMMsSqlManagedInstanceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ManagedInstancesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mssql_managed_instance" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		miName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, miName)
		if err != nil {
			return nil
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("MsSql Elastic Pool still exists %s: %+v", miName, err)
		}
	}

	return nil
}
