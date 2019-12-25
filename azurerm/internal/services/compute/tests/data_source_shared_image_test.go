package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMSharedImage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSharedImage_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMSharedImage_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSharedImage_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testAccDataSourceSharedImage_basic(data acceptance.TestData) string {
	template := testAccAzureRMSharedImage_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_shared_image" "test" {
  name                = "${azurerm_shared_image.test.name}"
  gallery_name        = "${azurerm_shared_image.test.gallery_name}"
  resource_group_name = "${azurerm_shared_image.test.resource_group_name}"
}
`, template)
}

func testAccDataSourceSharedImage_complete(data acceptance.TestData) string {
	template := testAccAzureRMSharedImage_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_shared_image" "test" {
  name                = "${azurerm_shared_image.test.name}"
  gallery_name        = "${azurerm_shared_image.test.gallery_name}"
  resource_group_name = "${azurerm_shared_image.test.resource_group_name}"
}
`, template)
}
