package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DataSourceManagedDiskExportDataSource struct{}

func TestAccDataSourceDiskExportSas_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_disk_export", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: DataSourceManagedDiskExportDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("managed_disk_id").HasValue("/subscriptions/42cbb0b8a331-abaf-4e69-8d8f-14b86a40/resourceGroups/disksrg/providers/Microsoft.Compute/disks/disk1"),
				check.That(data.ResourceName).Key("access").HasValue("Write"),
				check.That(data.ResourceName).Key("duration_in_seconds").HasValue("30"),
				check.That(data.ResourceName).Key("sas").Exists(),
			),
		},
	})
}

func (d DataSourceManagedDiskExportDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-disk-%d"
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
data "azurerm_managed_disk_export" "test" {
	managed_disk_id  = azurerm_managed_disk.disk.id
	duration_in_seconds = 300
	access = "Read"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
