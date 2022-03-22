package healthcare_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"testing"
)

type HealthCareWorkspaceIotConnectorResource struct{}

func TestAccHealthCareIotConnector_basic(t *testing.T) {
	data := acceptance.BuildTestData()
}

func (r HealthCareWorkspaceIotConnectorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
`)
}

func (HealthCareWorkspaceIotConnectorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-healthcareapi-%d"
  location = "%s"
}

resource "azurerm_healthcare_workspace" "test" {
  name                = "wk%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
