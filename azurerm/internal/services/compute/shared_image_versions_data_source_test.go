package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type SharedImageVersionsDataSource struct {
}

func TestAccDataSourceSharedImageVersions_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_versions", "test")
	r := SharedImageVersionsDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  SharedImageVersionResource{}.setup(data),
			Destroy: false,
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("images.0.managed_image_id").Exists(),
				check.That(data.ResourceName).Key("images.0.target_region.#").HasValue("1"),
				check.That(data.ResourceName).Key("images.0.target_region.0.storage_account_type").HasValue("Standard_LRS"),
			),
		},
	})
}

func TestAccDataSourceSharedImageVersions_tagsFilterError(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_versions", "test")
	r := SharedImageVersionsDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  SharedImageVersionResource{}.setup(data),
			Destroy: false,
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config:      r.tagsFilterError(data),
			ExpectError: regexp.MustCompile("unable to find any images"),
		},
	})
}

func TestAccDataSourceSharedImageVersions_tagsFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_versions", "test")
	r := SharedImageVersionsDataSource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  SharedImageVersionResource{}.setup(data),
			Destroy: false,
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.tagsFilter(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("images.#").HasValue("1"),
			),
		},
	})
}

func (SharedImageVersionsDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_versions" "test" {
  gallery_name        = azurerm_shared_image_version.test.gallery_name
  image_name          = azurerm_shared_image_version.test.image_name
  resource_group_name = azurerm_shared_image_version.test.resource_group_name
}
`, SharedImageVersionResource{}.imageVersion(data))
}

func (SharedImageVersionsDataSource) tagsFilterError(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_versions" "test" {
  gallery_name        = azurerm_shared_image_version.test.gallery_name
  image_name          = azurerm_shared_image_version.test.image_name
  resource_group_name = azurerm_shared_image_version.test.resource_group_name

  tags_filter = {
    "foo" = "error"
  }
}
`, SharedImageVersionResource{}.imageVersion(data))
}

func (SharedImageVersionsDataSource) tagsFilter(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_versions" "test" {
  gallery_name        = azurerm_shared_image_version.test.gallery_name
  image_name          = azurerm_shared_image_version.test.image_name
  resource_group_name = azurerm_shared_image_version.test.resource_group_name

  tags_filter = {
    "foo" = "bar"
  }
}
`, SharedImageVersionResource{}.imageVersion(data))
}
