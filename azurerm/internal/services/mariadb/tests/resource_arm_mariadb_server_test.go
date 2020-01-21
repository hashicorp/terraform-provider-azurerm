package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMariaDbServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "10.2"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_enforcement", "Enabled"),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMariaDbServer_basicOldSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbServer_basicOldSku(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "10.2"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_enforcement", "Enabled"),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMariaDbServer_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMariaDbServer_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_mariadb_server"),
			},
		},
	})
}

func TestAccAzureRMMariaDbServer_basicMaxStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbServer_basicMaxStorage(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "10.2"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_enforcement", "Enabled"),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMariaDbServer_generalPurpose(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbServer_generalPurpose(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMariaDbServer_memoryOptimized(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbServer_memoryOptimizedGeoRedundant(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMariaDbServer_updatePassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMMariaDbServer_basicUpdatedPassword(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMMariaDbServer_updated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "B_Gen5_2"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "10.2"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_profile.0.storage_mb", "51200"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
				),
			},
			{
				Config: testAccAzureRMMariaDbServer_basicUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "B_Gen5_1"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "10.2"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_profile.0.storage_mb", "640000"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
				),
			},
		},
	})
}

func TestAccAzureRMMariaDbServer_updateSKU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbServer_generalPurpose(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "GP_Gen5_32"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "32"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "GeneralPurpose"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.family", "Gen5"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_profile.0.storage_mb", "640000"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
				),
			},
			{
				Config: testAccAzureRMMariaDbServer_memoryOptimized(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "MO_Gen5_16"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "16"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "MemoryOptimized"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.family", "Gen5"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_profile.0.storage_mb", "4096000"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
				),
			},
		},
	})
}

func TestAccAzureRMMariaDbServer_storageAutogrow(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_profile.0.auto_grow", "Enabled"),
				),
			},
			{
				Config: testAccAzureRMMariaDbServer_storageAutogrowUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_profile.0.auto_grow", "Disabled"),
				),
			},
		},
	})
}

func testCheckAzureRMMariaDbServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MariaDB.ServersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for MariaDB Server: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MariaDB Server %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on mariadbServersClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMariaDbServerDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MariaDB.ServersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mariadb_server" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("MariaDB Server still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMMariaDbServer_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "B_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMariaDbServer_basicOldSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "B_Gen5_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMariaDbServer_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMariaDbServer_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mariadb_server" "import" {
  name                = "${azurerm_mariadb_server.test.name}"
  location            = "${azurerm_mariadb_server.test.location}"
  resource_group_name = "${azurerm_mariadb_server.test.resource_group_name}"

  sku_name = "B_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, template)
}

func testAccAzureRMMariaDbServer_basicUpdatedPassword(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "B_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "R3dH0TCh1l1P3pp3rs!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMariaDbServer_basicUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "B_Gen5_1"

  storage_profile {
    storage_mb            = 640000
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMariaDbServer_basicMaxStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "B_Gen5_2"

  storage_profile {
    storage_mb            = 947200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMariaDbServer_generalPurpose(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "GP_Gen5_32"

  storage_profile {
    storage_mb            = 640000
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMariaDbServer_memoryOptimized(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "MO_Gen5_16"

  storage_profile {
    storage_mb            = 4096000
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMariaDbServer_memoryOptimizedGeoRedundant(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "MO_Gen5_16"

  storage_profile {
    storage_mb            = 4096000
    backup_retention_days = 7
    geo_redundant_backup  = "Enabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMariaDbServer_storageAutogrowUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "B_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
    auto_grow             = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
