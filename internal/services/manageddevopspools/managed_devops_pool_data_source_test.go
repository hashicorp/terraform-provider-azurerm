package manageddevopspools_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagedDevOpsPoolDataSource struct{}

func TestAccManagedDevOpsPoolDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_devops_pool", "test")
	r := ManagedDevOpsPoolDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
			),
		},
	})
}

func (ManagedDevOpsPoolDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_managed_devops_pool" "test" {
  name     = "acctest-pool-%d"
  location = "%s"
}

data "azurerm_managed_devops_pool" "test" {
  name = azurerm_managed_devops_pool.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}
