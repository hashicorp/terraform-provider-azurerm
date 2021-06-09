package compute_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type SharedImageVersionDataSource struct {
}

func TestAccDataSourceSharedImageVersion_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_version", "test")
	r := SharedImageVersionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: SharedImageVersionResource{}.setup(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: SharedImageVersionResource{}.imageVersion(data),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("1"),
				check.That(data.ResourceName).Key("target_region.0.storage_account_type").HasValue("Standard_LRS"),
			),
		},
	})
}

func TestAccDataSourceSharedImageVersion_latest(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_version", "test")
	r := SharedImageVersionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: SharedImageVersionResource{}.setup(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.additionalVersion(data),
		},
		{
			Config: r.customName(data, "latest"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("1"),
				check.That(data.ResourceName).Key("target_region.0.storage_account_type").HasValue("Standard_LRS"),
			),
		},
	})
}

func TestAccDataSourceSharedImageVersion_recent(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_version", "test")
	r := SharedImageVersionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  SharedImageVersionResource{}.setup(data),
			Destroy: false,
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.additionalVersion(data),
		},
		{
			Config: r.customName(data, "recent"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("1"),
				check.That(data.ResourceName).Key("target_region.0.storage_account_type").HasValue("Standard_LRS"),
			),
		},
	})
}

func (SharedImageVersionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_version" "test" {
  name                = azurerm_shared_image_version.test.name
  gallery_name        = azurerm_shared_image_version.test.gallery_name
  image_name          = azurerm_shared_image_version.test.image_name
  resource_group_name = azurerm_shared_image_version.test.resource_group_name
}
`, SharedImageVersionResource{}.imageVersion(data))
}

func (SharedImageVersionDataSource) additionalVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test2" {
  name                = "0.0.2"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }
}
`, SharedImageVersionResource{}.imageVersion(data))
}

func (r SharedImageVersionDataSource) customName(data acceptance.TestData, name string) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_version" "test" {
  name                = "%s"
  gallery_name        = azurerm_shared_image_version.test2.gallery_name
  image_name          = azurerm_shared_image_version.test2.image_name
  resource_group_name = azurerm_shared_image_version.test2.resource_group_name
}
`, r.additionalVersion(data), name)
}
