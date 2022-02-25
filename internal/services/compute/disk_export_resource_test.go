package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DiskExportResource struct{}

func TestAccDiskExport_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_export", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: DiskExportResource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func (t DiskExportResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-revokedisk-%d"
  location = "%s"
}
resource "azurerm_managed_disk" "test" {
  name                 = "acctestsads%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"
}
resource "azurerm_disk_export" "test" {
  managed_disk_id     = azurerm_managed_disk.test.id
  duration_in_seconds = 300
  access              = "Read"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
