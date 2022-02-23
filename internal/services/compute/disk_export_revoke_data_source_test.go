package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DataSourceManagedDiskExportRevokeDataSource struct{}

func TestAccDataSourceDiskExportRevokeSas_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_disk_export", "revoke")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: DataSourceManagedDiskExportRevokeDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func (d DataSourceManagedDiskExportRevokeDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-revokedisk-%d"
  location = "%s"
}
resource "azurerm_managed_disk" "disk" {
  name                 = "acctestsads%s"
  location             = azurerm_resource_group.rg.location
  resource_group_name  = azurerm_resource_group.rg.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"
}
data "azurerm_managed_disk_export" "export" {
  managed_disk_id     = azurerm_managed_disk.disk.id
  duration_in_seconds = 300
  access              = "Read"
}
data "azurerm_managed_disk_export_revoke" "revoke" {
  depends_on      = [data.azurerm_managed_disk_export.export]
  managed_disk_id = azurerm_managed_disk.disk.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
