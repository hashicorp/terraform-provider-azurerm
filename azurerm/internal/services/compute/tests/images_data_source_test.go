package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMImages_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_images", "test")

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"
	storageType := "LRS"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then reference it in the image creation
				Config: testAccAzureRMImage_standaloneImage_setup(data, userName, password, hostName, storageType),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
				),
			},
			{
				Config: testAccAzureRMImage_standaloneImage_provision(data, userName, password, hostName, storageType, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageExists("azurerm_image.test", true),
				),
			},
			{
				Config: testAccDataSourceImages_basic(data, userName, password, hostName, storageType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "images.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "images.0.os_disk.0.os_type", "Linux"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMImages_tagsFilterError(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_images", "test")

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"
	storageType := "LRS"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then reference it in the image creation
				Config: testAccAzureRMImage_standaloneImage_setup(data, userName, password, hostName, storageType),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
				),
			},
			{
				Config: testAccAzureRMImage_standaloneImage_provision(data, userName, password, hostName, storageType, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageExists("azurerm_image.test", true),
				),
			},
			{
				Config:      testAccDataSourceImages_tagsFilterError(data, userName, password, hostName, storageType),
				ExpectError: regexp.MustCompile("no images were found that match the specified tags"),
			},
		},
	})
}

func TestAccDataSourceAzureRMImages_tagsFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_images", "test")

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"
	storageType := "LRS"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then reference it in the image creation
				Config: testAccAzureRMImage_standaloneImage_setup(data, userName, password, hostName, storageType),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
				),
			},
			{
				Config: testAccAzureRMImage_standaloneImage_provision(data, userName, password, hostName, storageType, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageExists("azurerm_image.test", true),
				),
			},
			{
				Config: testAccDataSourceImages_tagsFilter(data, userName, password, hostName, storageType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "images.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceImages_basic(data acceptance.TestData, userName, password, hostName, storageType string) string {
	template := testAccAzureRMImage_standaloneImage_provision(data, userName, password, hostName, storageType, "")
	return fmt.Sprintf(`
%s

data "azurerm_images" "test" {
  resource_group_name = azurerm_image.test.resource_group_name
}
`, template)
}

func testAccDataSourceImages_tagsFilterError(data acceptance.TestData, userName, password, hostName, storageType string) string {
	template := testAccAzureRMImage_standaloneImage_provision(data, userName, password, hostName, storageType, "")
	return fmt.Sprintf(`
%s

data "azurerm_images" "test" {
  resource_group_name = azurerm_image.test.resource_group_name
  tags_filter = {
    environment = "error"
  }
}
`, template)
}

func testAccDataSourceImages_tagsFilter(data acceptance.TestData, userName, password, hostName, storageType string) string {
	template := testAccAzureRMImage_standaloneImage_provision(data, userName, password, hostName, storageType, "")
	return fmt.Sprintf(`
%s

data "azurerm_images" "test" {
  resource_group_name = azurerm_image.test.resource_group_name
  tags_filter = {
    environment = "Dev"
  }
}
`, template)
}
