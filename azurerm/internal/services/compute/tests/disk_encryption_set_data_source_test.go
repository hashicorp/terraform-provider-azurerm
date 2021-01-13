package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DiskEncryptionSetDataSource struct {
}

func TestAccDataSourceDiskEncryptionSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetDataSource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
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
