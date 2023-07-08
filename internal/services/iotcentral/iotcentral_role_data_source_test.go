package iotcentral_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type IoTCentralRoleDataSource struct{}

func TestAccIoTCentralRole_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_iotcentral_role", "test")
	r := IoTCentralRoleDataSource{}

	displayName := "App Administrator"
	expectedId := "ca310b8d-2f4a-44e0-a36e-957c202cd8d4"

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(displayName),
				check.That(data.ResourceName).Key("id").HasValue(expectedId),
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
  display_name = "App Administrator"
}
`, data.RandomString, data.Locations.Primary, data.RandomString, data.RandomString)
}
