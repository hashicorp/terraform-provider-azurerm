// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlElasticPoolResource struct{}

func TestAccMsSqlElasticPool_basicDTU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDTU(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
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
			Config: r.basicDTU(data, ""),
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

	maintenance_configuration_name := ""
	switch data.Locations.Primary {
	case "westeurope":
		maintenance_configuration_name = "SQL_WestEurope_DB_2"
	case "francecentral":
		maintenance_configuration_name = "SQL_FranceCentral_DB_1"
	default:
		maintenance_configuration_name = "SQL_Default"
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardDTU(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("maintenance_configuration_name").HasValue(maintenance_configuration_name),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep("max_size_gb"),
		{
			Config: r.standardDTU(data, `enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("maintenance_configuration_name").HasValue(maintenance_configuration_name),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
			),
		},
		data.ImportStep("max_size_gb"),
	})
}

func TestAccMsSqlElasticPool_premiumDTUZoneRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	// Limited regional availability for ZRS
	data.Locations.Primary = "westeurope"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.premiumDTUZoneRedundant(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("PremiumPool"),
				check.That(data.ResourceName).Key("zone_redundant").HasValue("true"),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep("max_size_gb"),
		{
			Config: r.premiumDTUZoneRedundant(data, `enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("PremiumPool"),
				check.That(data.ResourceName).Key("zone_redundant").HasValue("true"),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
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
			Config: r.basicVCore(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep("max_size_gb"),
		{
			Config: r.basicVCore(data, `enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
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
			Config: r.basicVCoreMaxSizeBytes(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep("max_size_gb"),
		{
			Config: r.basicVCoreMaxSizeBytes(data, `enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
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
			Config: r.standardDTU(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep("max_size_gb"),
		{
			Config: r.resizeDTU(data, `enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
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
			Config: r.basicVCore(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep("max_size_gb"),
		{
			Config: r.resizeVCore(data, `enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
			),
		},
		data.ImportStep("max_size_gb"),
	})
}

func TestAccMsSqlElasticPool_fsv2FamilyVCore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	// Limited regional availability for Fsv2 family
	data.Locations.Primary = "westeurope"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fsv2VCore(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep("max_size_gb"),
		{
			Config: r.fsv2VCore(data, `enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
			),
		},
		data.ImportStep("max_size_gb"),
	})
}

func TestAccMsSqlElasticPool_dcFamilyVCore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	// Limited regional availability for DC-series
	data.Locations.Primary = "westeurope"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dcVCore(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep("max_size_gb"),
	})
}

func TestAccMsSqlElasticPool_dcFamilyBcTierVCore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	// Limited regional availability for DC-series
	data.Locations.Primary = "westeurope"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dcVCoreBc(data, ""),
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
			Config: r.licenseType(data, "LicenseIncluded", ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("license_type").HasValue("LicenseIncluded"),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep(),
		{
			Config: r.licenseType(data, "BasePrice", `enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("license_type").HasValue("BasePrice"),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlElasticPool_hyperScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hyperScale(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep("max_size_gb"),
		{
			Config: r.hyperScaleUpdate(data, `enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
			),
		},
		data.ImportStep("max_size_gb"),
	})
}

func TestAccMsSqlElasticPool_vCoreToStandardDTU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.noLicenseType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateToStandardDTU(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep(),
		{
			Config: r.licenseType(data, "LicenseIncluded", ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateToStandardDTU(data, `enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlElasticPool_enclaveTypeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDTU(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep("max_size_gb"),
		{
			Config: r.basicDTU(data, `enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
			),
		},
		data.ImportStep("max_size_gb"),
		{
			Config: r.basicDTU(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep("max_size_gb"),
	})
}

func TestAccMsSqlElasticPool_dcFamilyVCoreEnclaveTypeError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MsSqlElasticPoolResource{}

	// Limited regional availability for DC-series
	data.Locations.Primary = "westeurope"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.dcVCore(data, `enclave_type = "VBS"`),
			ExpectError: regexp.MustCompile(`virtualization based security \(VBS\) enclaves are not supported`),
		},
	})
}

func (MsSqlElasticPoolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseSqlElasticPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	existing, err := client.MSSQL.ElasticPoolsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return nil, fmt.Errorf("SQL Elastic Pool %s does not exist", id)
		}

		return nil, fmt.Errorf("reading SQL Elastic Pool %s: %v", id, err)
	}

	return utils.Bool(existing.Model.Id != nil), nil
}

func (r MsSqlElasticPoolResource) basicDTU(data acceptance.TestData, enclaveType string) string {
	return r.templateDTU(data, "BasicPool", "Basic", 50, 4.8828125, 0, 5, false, enclaveType)
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
`, r.templateDTU(data, "BasicPool", "Basic", 50, 4.8828125, 0, 5, false, ""))
}

func (r MsSqlElasticPoolResource) premiumDTUZoneRedundant(data acceptance.TestData, enclaveType string) string {
	return r.templateDTU(data, "PremiumPool", "Premium", 125, 50, 0, 50, true, enclaveType)
}

func (r MsSqlElasticPoolResource) standardDTU(data acceptance.TestData, enclaveType string) string {
	return r.templateDTU(data, "StandardPool", "Standard", 50, 50, 0, 50, false, enclaveType)
}

func (r MsSqlElasticPoolResource) resizeDTU(data acceptance.TestData, enclaveType string) string {
	return r.templateDTU(data, "StandardPool", "Standard", 100, 100, 50, 100, false, enclaveType)
}

func (r MsSqlElasticPoolResource) updateToStandardDTU(data acceptance.TestData, enclaveType string) string {
	return r.templateUpdateToDTU(data, "StandardPool", "Standard", 50, 50, 10, 50, false, enclaveType)
}

func (r MsSqlElasticPoolResource) basicVCore(data acceptance.TestData, enclaveType string) string {
	return r.templateVCore(data, "GP_Gen5", "GeneralPurpose", 4, "Gen5", 0.25, 4, enclaveType)
}

func (r MsSqlElasticPoolResource) fsv2VCore(data acceptance.TestData, enclaveType string) string {
	return r.templateVCore(data, "GP_Fsv2", "GeneralPurpose", 8, "Fsv2", 0, 8, enclaveType)
}

func (r MsSqlElasticPoolResource) dcVCore(data acceptance.TestData, enclaveType string) string {
	return r.templateVCore(data, "GP_DC", "GeneralPurpose", 2, "DC", 2, 2, enclaveType)
}

func (r MsSqlElasticPoolResource) dcVCoreBc(data acceptance.TestData, enclaveType string) string {
	return r.templateVCore(data, "BC_DC", "BusinessCritical", 2, "DC", 2, 2, enclaveType)
}

func (r MsSqlElasticPoolResource) basicVCoreMaxSizeBytes(data acceptance.TestData, enclaveType string) string {
	return r.templateVCoreMaxSizeBytes(data, "GP_Gen5", "GeneralPurpose", 4, "Gen5", 0.25, 4, enclaveType)
}

func (r MsSqlElasticPoolResource) resizeVCore(data acceptance.TestData, enclaveType string) string {
	return r.templateVCore(data, "GP_Gen5", "GeneralPurpose", 8, "Gen5", 0, 8, enclaveType)
}

func (MsSqlElasticPoolResource) templateDTU(data acceptance.TestData, skuName string, skuTier string, skuCapacity int, maxSizeGB float64, databaseSettingsMin int, databaseSettingsMax int, zoneRedundant bool, enclaveType string) string {
	configName := "SQL_Default"
	if skuTier != "Basic" {
		switch data.Locations.Primary {
		case "westeurope":
			configName = "SQL_WestEurope_DB_2"
		case "francecentral":
			configName = "SQL_FranceCentral_DB_1"
		default:
			configName = "SQL_Default"
		}
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "test" {
  name                           = "acctest-pool-dtu-%[1]d"
  resource_group_name            = azurerm_resource_group.test.name
  location                       = azurerm_resource_group.test.location
  server_name                    = azurerm_mssql_server.test.name
  max_size_gb                    = %.7[6]f
  zone_redundant                 = %[9]t
  maintenance_configuration_name = "%[10]s"
  %[11]s

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
`, data.RandomInteger, data.Locations.Primary, skuName, skuTier, skuCapacity, maxSizeGB, databaseSettingsMin, databaseSettingsMax, zoneRedundant, configName, enclaveType)
}

func (MsSqlElasticPoolResource) templateUpdateToDTU(data acceptance.TestData, skuName string, skuTier string, skuCapacity int, maxSizeGB float64, databaseSettingsMin int, databaseSettingsMax int, zoneRedundant bool, enclaveType string) string {
	configName := "SQL_Default"
	if skuTier != "Basic" {
		switch data.Locations.Primary {
		case "westeurope":
			configName = "SQL_WestEurope_DB_2"
		case "francecentral":
			configName = "SQL_FranceCentral_DB_1"
		default:
			configName = "SQL_Default"
		}
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
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
  server_name         = azurerm_mssql_server.test.name
  max_size_gb         = %.7[6]f
  zone_redundant      = %[9]t
  %[11]s

  maintenance_configuration_name = "%[10]s"

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
`, data.RandomInteger, data.Locations.Primary, skuName, skuTier, skuCapacity, maxSizeGB, databaseSettingsMin, databaseSettingsMax, zoneRedundant, configName, enclaveType)
}

func (MsSqlElasticPoolResource) templateVCore(data acceptance.TestData, skuName string, skuTier string, skuCapacity int, skuFamily string, databaseSettingsMin float64, databaseSettingsMax float64, enclaveType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
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
  server_name         = azurerm_mssql_server.test.name
  max_size_gb         = 5
  %[9]s
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
`, data.RandomInteger, data.Locations.Primary, skuName, skuTier, skuCapacity, skuFamily, databaseSettingsMin, databaseSettingsMax, enclaveType)
}

func (MsSqlElasticPoolResource) templateVCoreMaxSizeBytes(data acceptance.TestData, skuName string, skuTier string, skuCapacity int, skuFamily string, databaseSettingsMin float64, databaseSettingsMax float64, enclaveType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
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
  server_name         = azurerm_mssql_server.test.name
  max_size_bytes      = 214748364800
  %[9]s

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
`, data.RandomInteger, data.Locations.Primary, skuName, skuTier, skuCapacity, skuFamily, databaseSettingsMin, databaseSettingsMax, enclaveType)
}

func (MsSqlElasticPoolResource) templateHyperScale(data acceptance.TestData, skuName string, skuTier string, skuCapacity int, skuFamily string, databaseSettingsMin float64, databaseSettingsMax float64, enclaveType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
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
  server_name         = azurerm_mssql_server.test.name
  %[9]s
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
`, data.RandomInteger, data.Locations.Primary, skuName, skuTier, skuCapacity, skuFamily, databaseSettingsMin, databaseSettingsMax, enclaveType)
}

func (MsSqlElasticPoolResource) noLicenseType(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
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
  server_name         = azurerm_mssql_server.test.name
  max_size_gb         = 50
  zone_redundant      = false

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
`, data.RandomInteger, data.Locations.Primary)
}

func (MsSqlElasticPoolResource) licenseType(data acceptance.TestData, licenseType string, enclaveType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
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
  server_name         = azurerm_mssql_server.test.name
  max_size_gb         = 50
  zone_redundant      = false
  license_type        = "%[3]s"
  %[4]s

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
`, data.RandomInteger, data.Locations.Primary, licenseType, enclaveType)
}

func (r MsSqlElasticPoolResource) hyperScale(data acceptance.TestData, enclaveType string) string {
	return r.templateHyperScale(data, "HS_Gen5", "Hyperscale", 4, "Gen5", 0.25, 4, enclaveType)
}

func (r MsSqlElasticPoolResource) hyperScaleUpdate(data acceptance.TestData, enclaveType string) string {
	return r.templateHyperScale(data, "HS_Gen5", "Hyperscale", 4, "Gen5", 0, 4, enclaveType)
}
