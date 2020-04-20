package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDiskEncryptionSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_disk_encryption_set", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDiskEncryptionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDiskEncryptionSet_basic(data),
			},
			{
				Config: testAccDataSourceDiskEncryptionSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
				),
			},
		},
	})
}

func testAccDataSourceDiskEncryptionSet_basic(data acceptance.TestData) string {
	config := testAccAzureRMDiskEncryptionSet_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_disk_encryption_set" "test" {
  name                = azurerm_disk_encryption_set.test.name
  resource_group_name = azurerm_disk_encryption_set.test.resource_group_name
}
`, config)
}
