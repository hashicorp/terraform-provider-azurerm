package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DiskEncryptionSetDataSource struct {
}

func TestAccDataSourceDiskEncryptionSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("auto_key_rotation_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccDataSourceDiskEncryptionSet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("auto_key_rotation_enabled").HasValue("true"),
			),
		},
	})
}

func (DiskEncryptionSetDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_disk_encryption_set" "test" {
  name                = azurerm_disk_encryption_set.test.name
  resource_group_name = azurerm_disk_encryption_set.test.resource_group_name
}
`, DiskEncryptionSetResource{}.basic(data))
}

func (DiskEncryptionSetDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_disk_encryption_set" "test" {
  name                = azurerm_disk_encryption_set.test.name
  resource_group_name = azurerm_disk_encryption_set.test.resource_group_name
}
`, DiskEncryptionSetResource{}.complete(data))
}
