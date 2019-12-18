package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSharedImageVersion_basic(t *testing.T) {
	resourceName := "azurerm_shared_image_version.test"

	ri := tf.AccRandTimeInt()
	resourceGroup := fmt.Sprintf("acctestRG-%d", ri)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", ri)
	sshPort := "22"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageVersionDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then reference it in the image creation
				Config:  testAccAzureRMSharedImageVersion_setup(ri, acceptance.Location(), userName, password, hostName),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, acceptance.Location()),
				),
			},
			{
				Config: testAccAzureRMSharedImageVersion_imageVersion(ri, acceptance.Location(), userName, password, hostName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageVersionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "managed_image_id"),
					resource.TestCheckResourceAttr(resourceName, "target_region.#", "1"),
				),
			},
			{
				Config: testAccAzureRMSharedImageVersion_imageVersionUpdated(ri, acceptance.Location(), acceptance.AltLocation(), userName, password, hostName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageVersionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "managed_image_id"),
					resource.TestCheckResourceAttr(resourceName, "target_region.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "name", "1234567890.1234567890.1234567890"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSharedImageVersion_storageAccountTypeLrs(t *testing.T) {
	resourceName := "azurerm_shared_image_version.test"

	ri := tf.AccRandTimeInt()
	resourceGroup := fmt.Sprintf("acctestRG-%d", ri)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", ri)
	sshPort := "22"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageVersionDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then reference it in the image creation
				Config:  testAccAzureRMSharedImageVersion_setup(ri, acceptance.Location(), userName, password, hostName),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, acceptance.Location()),
				),
			},
			{
				Config: testAccAzureRMSharedImageVersion_imageVersionStorageAccountType(ri, acceptance.Location(), userName, password, hostName, "Standard_LRS"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageVersionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "managed_image_id"),
					resource.TestCheckResourceAttr(resourceName, "target_region.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSharedImageVersion_storageAccountTypeZrs(t *testing.T) {
	resourceName := "azurerm_shared_image_version.test"

	ri := tf.AccRandTimeInt()
	resourceGroup := fmt.Sprintf("acctestRG-%d", ri)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", ri)
	sshPort := "22"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageVersionDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then reference it in the image creation
				Config:  testAccAzureRMSharedImageVersion_setup(ri, acceptance.Location(), userName, password, hostName),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, acceptance.Location()),
				),
			},
			{
				Config: testAccAzureRMSharedImageVersion_imageVersionStorageAccountType(ri, acceptance.Location(), userName, password, hostName, "Standard_ZRS"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageVersionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "managed_image_id"),
					resource.TestCheckResourceAttr(resourceName, "target_region.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSharedImageVersion_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	resourceName := "azurerm_shared_image_version.test"

	ri := tf.AccRandTimeInt()
	resourceGroup := fmt.Sprintf("acctestRG-%d", ri)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", ri)
	sshPort := "22"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageVersionDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then reference it in the image creation
				Config:  testAccAzureRMSharedImageVersion_setup(ri, acceptance.Location(), userName, password, hostName),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, acceptance.Location()),
				),
			},
			{
				Config: testAccAzureRMSharedImageVersion_imageVersion(ri, acceptance.Location(), userName, password, hostName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageVersionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "managed_image_id"),
					resource.TestCheckResourceAttr(resourceName, "target_region.#", "1"),
				),
			},
			{
				Config:      testAccAzureRMSharedImageVersion_requiresImport(ri, acceptance.Location(), userName, password, hostName),
				ExpectError: acceptance.RequiresImportError("azurerm_shared_image_version"),
			},
		},
	})
}

func testCheckAzureRMSharedImageVersionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.GalleryImageVersionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_shared_image_version" {
			continue
		}

		imageVersion := rs.Primary.Attributes["name"]
		imageName := rs.Primary.Attributes["image_name"]
		galleryName := rs.Primary.Attributes["gallery_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, galleryName, imageName, imageVersion, "")
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("Shared Image Version still exists:\n%+v", resp)
	}

	return nil
}

func testCheckAzureRMSharedImageVersionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		imageVersion := rs.Primary.Attributes["name"]
		imageName := rs.Primary.Attributes["image_name"]
		galleryName := rs.Primary.Attributes["gallery_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Shared Image Version: %s", imageName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.GalleryImageVersionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, galleryName, imageName, imageVersion, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on galleryImageVersionsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Shared Image Version %q (Image %q / Gallery %q / Resource Group: %q) does not exist", imageVersion, imageName, galleryName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMSharedImageVersion_setup(rInt int, location, username, password, hostname string) string {
	return testAccAzureRMImage_standaloneImage_setup(rInt, username, password, hostname, location, "LRS")
}

func testAccAzureRMSharedImageVersion_provision(rInt int, location, username, password, hostname string) string {
	template := testAccAzureRMImage_standaloneImage_provision(rInt, username, password, hostname, location, "LRS", "")
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%d"
  gallery_name        = "${azurerm_shared_image_gallery.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  os_type             = "Linux"

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
  }
}
`, template, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMSharedImageVersion_imageVersion(rInt int, location, username, password, hostname string) string {
	template := testAccAzureRMSharedImageVersion_provision(rInt, location, username, password, hostname)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = "${azurerm_shared_image_gallery.test.name}"
  image_name          = "${azurerm_shared_image.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  managed_image_id    = "${azurerm_image.test.id}"

  target_region {
    name                   = "${azurerm_resource_group.test.location}"
    regional_replica_count = 1
  }
}
`, template)
}

func testAccAzureRMSharedImageVersion_imageVersionStorageAccountType(rInt int, location, username, password, hostname string, storageAccountType string) string {
	template := testAccAzureRMSharedImageVersion_provision(rInt, location, username, password, hostname)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = "${azurerm_shared_image_gallery.test.name}"
  image_name          = "${azurerm_shared_image.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  managed_image_id    = "${azurerm_image.test.id}"

  target_region {
    name                   = "${azurerm_resource_group.test.location}"
	regional_replica_count = 1
	storage_account_type   = "%s"
  }
}
`, template, storageAccountType)
}

func testAccAzureRMSharedImageVersion_requiresImport(rInt int, location, username, password, hostname string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "import" {
  name                = "${azurerm_shared_image_version.test.name}"
  gallery_name        = "${azurerm_shared_image_version.test.gallery_name}"
  image_name          = "${azurerm_shared_image_version.test.image_name}"
  resource_group_name = "${azurerm_shared_image_version.test.resource_group_name}"
  location            = "${azurerm_shared_image_version.test.location}"
  managed_image_id    = "${azurerm_shared_image_version.test.managed_image_id}"

  target_region {
    name                   = "${azurerm_resource_group.test.location}"
    regional_replica_count = 1
  }
}
`, testAccAzureRMSharedImageVersion_imageVersion(rInt, location, username, password, hostname))
}

func testAccAzureRMSharedImageVersion_imageVersionUpdated(rInt int, location, altLocation, username, password, hostname string) string {
	template := testAccAzureRMSharedImageVersion_provision(rInt, location, username, password, hostname)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "1234567890.1234567890.1234567890"
  gallery_name        = "${azurerm_shared_image_gallery.test.name}"
  image_name          = "${azurerm_shared_image.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  managed_image_id    = "${azurerm_image.test.id}"

  target_region {
    name                   = "${azurerm_resource_group.test.location}"
    regional_replica_count = 1
  }

  target_region {
    name                   = "%s"
    regional_replica_count = 2
  }
}
`, template, altLocation)
}
