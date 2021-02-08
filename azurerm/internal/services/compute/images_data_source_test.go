package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ImagesDataSource struct {
}

func TestAccDataSourceAzureRMImages_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_images", "test")
	r := ImagesDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: ImageResource{}.setupUnmanagedDisks(data, "LRS"),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: ImageResource{}.standaloneImageProvision(data, "LRS", ""),
			Check:  resource.ComposeTestCheckFunc(),
		},
		{
			Config: r.basic(data, "LRS"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("images.#").HasValue("1"),
				check.That(data.ResourceName).Key("images.0.os_disk.0.os_type").HasValue("Linux"),
			),
		},
	})
}

func TestAccDataSourceAzureRMImages_tagsFilterError(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_images", "test")
	r := ImagesDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: ImageResource{}.setupUnmanagedDisks(data, "LRS"),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: ImageResource{}.standaloneImageProvision(data, "LRS", ""),
			Check:  resource.ComposeTestCheckFunc(),
		},
		{
			Config:      r.tagsFilterError(data, "LRS"),
			ExpectError: regexp.MustCompile("no images were found that match the specified tags"),
		},
	})
}

func TestAccDataSourceAzureRMImages_tagsFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_images", "test")
	r := ImagesDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: ImageResource{}.setupUnmanagedDisks(data, "LRS"),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: ImageResource{}.standaloneImageProvision(data, "LRS", ""),
			Check:  resource.ComposeTestCheckFunc(),
		},
		{
			Config: r.tagsFilter(data, "LRS"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("images.#").HasValue("1"),
			),
		},
	})
}

func (ImagesDataSource) basic(data acceptance.TestData, storageType string) string {
	return fmt.Sprintf(`
%s

data "azurerm_images" "test" {
  resource_group_name = azurerm_image.test.resource_group_name
}
`, ImageResource{}.standaloneImageProvision(data, storageType, ""))
}

func (ImagesDataSource) tagsFilterError(data acceptance.TestData, storageType string) string {
	return fmt.Sprintf(`
%s

data "azurerm_images" "test" {
  resource_group_name = azurerm_image.test.resource_group_name
  tags_filter = {
    environment = "error"
  }
}
`, ImageResource{}.standaloneImageProvision(data, storageType, ""))
}

func (ImagesDataSource) tagsFilter(data acceptance.TestData, storageType string) string {
	return fmt.Sprintf(`
%s

data "azurerm_images" "test" {
  resource_group_name = azurerm_image.test.resource_group_name
  tags_filter = {
    environment = "Dev"
  }
}
`, ImageResource{}.standaloneImageProvision(data, storageType, ""))
}
