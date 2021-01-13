package compute

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

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"
	storageType := "LRS"

	data.DataSourceTest(t, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: ImageResource{}.standaloneImage_setup(data, userName, password, hostName, storageType),
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: ImageResource{}.standaloneImage_provision(data, userName, password, hostName, storageType, ""),
			Check:  resource.ComposeTestCheckFunc(),
		},
		{
			Config: r.basic(data, userName, password, hostName, storageType),
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

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"
	storageType := "LRS"

	data.DataSourceTest(t, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: ImageResource{}.standaloneImage_setup(data, userName, password, hostName, storageType),
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: ImageResource{}.standaloneImage_provision(data, userName, password, hostName, storageType, ""),
			Check:  resource.ComposeTestCheckFunc(),
		},
		{
			Config:      r.tagsFilterError(data, userName, password, hostName, storageType),
			ExpectError: regexp.MustCompile("no images were found that match the specified tags"),
		},
	})
}

func TestAccDataSourceAzureRMImages_tagsFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_images", "test")
	r := ImagesDataSource{}

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"
	storageType := "LRS"

	data.DataSourceTest(t, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: ImageResource{}.standaloneImage_setup(data, userName, password, hostName, storageType),
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: ImageResource{}.standaloneImage_provision(data, userName, password, hostName, storageType, ""),
			Check:  resource.ComposeTestCheckFunc(),
		},
		{
			Config: r.tagsFilter(data, userName, password, hostName, storageType),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("images.#").HasValue("1"),
			),
		},
	})
}

func (ImagesDataSource) basic(data acceptance.TestData, userName, password, hostName, storageType string) string {
	return fmt.Sprintf(`
%s

data "azurerm_images" "test" {
  resource_group_name = azurerm_image.test.resource_group_name
}
`, ImageResource{}.standaloneImage_provision(data, userName, password, hostName, storageType, ""))
}

func (ImagesDataSource) tagsFilterError(data acceptance.TestData, userName, password, hostName, storageType string) string {
	return fmt.Sprintf(`
%s

data "azurerm_images" "test" {
  resource_group_name = azurerm_image.test.resource_group_name
  tags_filter = {
    environment = "error"
  }
}
`, ImageResource{}.standaloneImage_provision(data, userName, password, hostName, storageType, ""))
}

func (ImagesDataSource) tagsFilter(data acceptance.TestData, userName, password, hostName, storageType string) string {
	return fmt.Sprintf(`
%s

data "azurerm_images" "test" {
  resource_group_name = azurerm_image.test.resource_group_name
  tags_filter = {
    environment = "Dev"
  }
}
`, ImageResource{}.standaloneImage_provision(data, userName, password, hostName, storageType, ""))
}
