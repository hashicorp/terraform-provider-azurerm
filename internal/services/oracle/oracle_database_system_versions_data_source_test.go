package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DatabaseVersionsDataSource struct{}

func TestDbVersionsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_database_system_versions", "test")
	r := DatabaseVersionsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.0.version").Exists(),
				check.That(data.ResourceName).Key("versions.0.pluggable_database_supported").Exists(),
			),
		},
	})
}

func TestDbVersionsDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_database_system_versions", "test")
	r := DatabaseVersionsDataSource{}

	const testShape = "VM.Standard.x86"

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data, testShape),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.0.version").Exists(),
				check.That(data.ResourceName).Key("versions.0.pluggable_database_supported").Exists(),
			),
		},
	})
}

func TestDbVersionsDataSource_ShapeFamilyFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_database_system_versions", "test")
	r := DatabaseVersionsDataSource{}

	const testFamily = "VIRTUALMACHINE"

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.shapeFamilyFilter(data, testFamily),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.0.version").Exists(),
				check.That(data.ResourceName).Key("versions.0.pluggable_database_supported").Exists(),
			),
		},
	})
}

func (d DatabaseVersionsDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_oracle_database_system_versions" "test" {
  location = "%s"
}
`, data.Locations.Primary)
}

func (d DatabaseVersionsDataSource) complete(data acceptance.TestData, shape string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_oracle_database_system_versions" "test" {
  location                          = "%[1]s"
  database_system_shape             = "%[2]s"
  upgrade_supported                 = true
  database_software_image_supported = true
}
`, data.Locations.Primary, shape)
}

func (d DatabaseVersionsDataSource) shapeFamilyFilter(data acceptance.TestData, family string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_oracle_database_system_versions" "test" {
  location                          = "%[1]s"
  shape_family                      = "%[2]s"
  upgrade_supported                 = true
  database_software_image_supported = false
}
`, data.Locations.Primary, family)
}
