package iotcentral_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

const (
	appAdminRoleDisplayName   = "App Administrator"
	appAdminRoleId            = "ca310b8d-2f4a-44e0-a36e-957c202cd8d4"
	appBuilderRoleDisplayName = "App Builder"
	appBuilderRoleId          = "344138e9-8de4-4497-8c54-5237e96d6aaf"
	orgAdminRoleDisplayName   = "Org Administrator"
	orgAdminRoleId            = "c495eb57-eb18-489e-9802-62c474e5645c"
)

type IoTCentralRoleDataSource struct{}

func TestAccIoTCentralRole_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_iotcentral_role", "test")
	r := IoTCentralRoleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(appAdminRoleDisplayName),
				check.That(data.ResourceName).Key("id").HasValue(appAdminRoleId),
			),
		},
	})
}

func (d IoTCentralRoleDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "iotctroletest-%s"
  location = "%s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "iotctroletest-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sub_domain          = "iotctroletest-%s"

  sku = "ST0"
}

data "azurerm_iotcentral_role" "test" {
  sub_domain   = azurerm_iotcentral_application.test.sub_domain
  display_name = "%s"
}
`, data.RandomString, data.Locations.Primary, data.RandomString, data.RandomString, appAdminRoleDisplayName)
}
