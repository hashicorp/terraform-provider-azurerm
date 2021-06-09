package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MsSqlElasticPoolResource struct{}

func TestAccMsSqlElasticPool_basicDTU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDTU(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("max_size_gb"),
	})
}

func TestAccMsSqlElasticPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDTU(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMsSqlElasticPool_standardDTU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardDTU(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("max_size_gb"),
	})
}

func TestAccMsSqlElasticPool_premiumDTUZoneRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.premiumDTUZoneRedundant(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "PremiumPool"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "zone_redundant", "true"),
			),
		},
		data.ImportStep("max_size_gb"),
	})
}

func TestAccMsSqlElasticPool_basicVCore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicVCore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("max_size_gb"),
	})
}

func TestAccMsSqlElasticPool_basicVCoreMaxSizeBytes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicVCoreMaxSizeBytes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("max_size_gb"),
	})
}

func TestAccMsSqlElasticPool_resizeDTU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardDTU(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("max_size_gb"),
		{
			Config: r.resizeDTU(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("max_size_gb"),
	})
}

func TestAccMsSqlElasticPool_resizeVCore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicVCore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("max_size_gb"),
		{
			Config: r.resizeVCore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("max_size_gb"),
	})
}

func TestAccMsSqlElasticPool_licenseType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.licenseType(data, "LicenseIncluded"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "license_type", "LicenseIncluded"),
			),
		},
		data.ImportStep(),
		{
			Config: r.licenseType(data, "BasePrice"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "license_type", "BasePrice"),
			),
		},
		data.ImportStep(),
	})
}

func (MsSqlElasticPoolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ElasticPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.ElasticPoolsClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil, fmt.Errorf("SQL Elastic Pool %q (Server %q, Resource Group %q) does not exist", id.Name, id.ServerName, id.ResourceGroup)
		}
		return nil, fmt.Errorf("reading SQL Elastic Pool %q (Server %q, Resource Group %q): %v", id.Name, id.ServerName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r MsSqlElasticPoolResource) basicDTU(data acceptance.TestData) string {
	return r.templateDTU(data, "BasicPool", "Basic", 50, 4.8828125, 0, 5, false)
}

func (r MsSqlElasticPoolResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_elasticpool" "import" {
  name                = azurerm_mssql_elasticpool.test.name
  resource_group_name = azurerm_mssql_elasticpool.test.resource_group_name
  location            = azurerm_mssql_elasticpool.test.location
  server_name         = azurerm_mssql_elasticpool.test.server_name
  max_size_gb         = 4.8828125

  sku {
    name     = "BasicPool"
    tier     = "Basic"
    capacity = 50
  }

  per_database_settings {
    min_capacity = 0
    max_capacity = 5
  }
}
`, r.templateDTU(data, "BasicPool", "Basic", 50, 4.8828125, 0, 5, false))
}

func (r MsSqlElasticPoolResource) premiumDTUZoneRedundant(data acceptance.TestData) string {
	return r.templateDTU(data, "PremiumPool", "Premium", 125, 50, 0, 50, true)
}

func (r MsSqlElasticPoolResource) standardDTU(data acceptance.TestData) string {
	return r.templateDTU(data, "StandardPool", "Standard", 50, 50, 0, 50, false)
}

func (r MsSqlElasticPoolResource) resizeDTU(data acceptance.TestData) string {
	return r.templateDTU(data, "StandardPool", "Standard", 100, 100, 50, 100, false)
}

func (r MsSqlElasticPoolResource) basicVCore(data acceptance.TestData) string {
	return r.templateVCore(data, "GP_Gen5", "GeneralPurpose", 4, "Gen5", 0.25, 4)
}

func (r MsSqlElasticPoolResource) basicVCoreMaxSizeBytes(data acceptance.TestData) string {
	return r.templateVCoreMaxSizeBytes(data, "GP_Gen5", "GeneralPurpose", 4, "Gen5", 0.25, 4)
}

func (r MsSqlElasticPoolResource) resizeVCore(data acceptance.TestData) string {
	return r.templateVCore(data, "GP_Gen5", "GeneralPurpose", 8, "Gen5", 0, 8)
}

func (MsSqlElasticPoolResource) templateDTU(data acceptance.TestData, skuName string, skuTier string, skuCapacity int, maxSizeGB float64, databaseSettingsMin int, databaseSettingsMax int, zoneRedundant bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-dtu-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  max_size_gb         = %.7[6]f
  zone_redundant      = %[9]t

  sku {
    name     = "%[3]s"
    tier     = "%[4]s"
    capacity = %[5]d
  }

  per_database_settings {
    min_capacity = %[7]d
    max_capacity = %[8]d
  }
}
`, data.RandomInteger, data.Locations.Primary, skuName, skuTier, skuCapacity, maxSizeGB, databaseSettingsMin, databaseSettingsMax, zoneRedundant)
}

func (MsSqlElasticPoolResource) templateVCore(data acceptance.TestData, skuName string, skuTier string, skuCapacity int, skuFamily string, databaseSettingsMin float64, databaseSettingsMax float64) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-vcore-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  max_size_gb         = 5

  sku {
    name     = "%[3]s"
    tier     = "%[4]s"
    capacity = %[5]d
    family   = "%[6]s"
  }

  per_database_settings {
    min_capacity = %.2[7]f
    max_capacity = %.2[8]f
  }
}
`, data.RandomInteger, data.Locations.Primary, skuName, skuTier, skuCapacity, skuFamily, databaseSettingsMin, databaseSettingsMax)
}

func (MsSqlElasticPoolResource) templateVCoreMaxSizeBytes(data acceptance.TestData, skuName string, skuTier string, skuCapacity int, skuFamily string, databaseSettingsMin float64, databaseSettingsMax float64) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-vcore-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  max_size_bytes      = 214748364800

  sku {
    name     = "%[3]s"
    tier     = "%[4]s"
    capacity = %[5]d
    family   = "%[6]s"
  }

  per_database_settings {
    min_capacity = %.2[7]f
    max_capacity = %.2[8]f
  }
}
`, data.RandomInteger, data.Locations.Primary, skuName, skuTier, skuCapacity, skuFamily, databaseSettingsMin, databaseSettingsMax)
}

func (MsSqlElasticPoolResource) licenseType(data acceptance.TestData, licenseType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-dtu-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  max_size_gb         = 50
  zone_redundant      = false
  license_type        = "%[3]s"

  sku {
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
    capacity = 4
    family   = "Gen5"
  }

  per_database_settings {
    min_capacity = 0
    max_capacity = 4
  }

}
`, data.RandomInteger, data.Locations.Primary, licenseType)
}
