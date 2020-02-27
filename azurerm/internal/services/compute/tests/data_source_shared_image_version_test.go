package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMSharedImageVersion_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_version", "test")
	username := "testadmin"
	password := "Password1234!"
	hostname := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageVersionDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then reference it in the image creation
				Config:  testAccAzureRMSharedImageVersion_setup(data, username, password, hostname),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", username, password, hostname, "22", data.Locations.Primary),
				),
			},
			{
				Config: testAccDataSourceSharedImageVersion_basic(data, username, password, hostname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_image_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_region.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_region.0.storage_account_type", "Standard_LRS"),
				),
			},
		},
	})
}

func testAccDataSourceSharedImageVersion_basic(data acceptance.TestData, username, password, hostname string) string {
	template := testAccAzureRMSharedImageVersion_imageVersion(data, username, password, hostname)
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_version" "test" {
  name                = "${azurerm_shared_image_version.test.name}"
  gallery_name        = "${azurerm_shared_image_version.test.gallery_name}"
  image_name          = "${azurerm_shared_image_version.test.image_name}"
  resource_group_name = "${azurerm_shared_image_version.test.resource_group_name}"
}
`, template)
}
