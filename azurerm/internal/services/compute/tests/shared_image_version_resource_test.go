package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSharedImageVersion_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageVersionDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then reference it in the image creation
				Config: testAccAzureRMSharedImageVersion_setup(data, userName, password, hostName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
				),
			},
			{
				Config: testAccAzureRMSharedImageVersion_imageVersion(data, userName, password, hostName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageVersionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_image_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_region.#", "1"),
				),
			},
			{
				Config: testAccAzureRMSharedImageVersion_imageVersionUpdated(data, userName, password, hostName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageVersionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_image_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_region.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "1234567890.1234567890.1234567890"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSharedImageVersion_storageAccountTypeLrs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageVersionDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then reference it in the image creation
				Config:  testAccAzureRMSharedImageVersion_setup(data, userName, password, hostName),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
				),
			},
			{
				Config: testAccAzureRMSharedImageVersion_imageVersionStorageAccountType(data, userName, password, hostName, "Standard_LRS"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageVersionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_image_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_region.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSharedImageVersion_storageAccountTypeZrs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageVersionDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then reference it in the image creation
				Config:  testAccAzureRMSharedImageVersion_setup(data, userName, password, hostName),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
				),
			},
			{
				Config: testAccAzureRMSharedImageVersion_imageVersionStorageAccountType(data, userName, password, hostName, "Standard_ZRS"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageVersionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_image_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_region.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSharedImageVersion_specializedImageVersionBySnapshot(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")

	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSharedImageVersion_imageVersionSpecializedBySnapshot(data, userName, password, hostName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageVersionExists(data.ResourceName),
					// the share image version will generate a shared access signature (SAS) on the referenced snapshot and keep it active until the replication is complete
					// in the meantime, the service will return success of creation before the replication complete.
					// therefore in this test, we have to revoke the access on the snapshot in order to do the cleaning work
					revokeSASOnSnapshotForCleanup("azurerm_snapshot.test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSharedImageVersion_specializedImageVersionByVM(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")

	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSharedImageVersion_imageVersionSpecializedByVM(data, userName, password, hostName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageVersionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSharedImageVersion_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageVersionDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then reference it in the image creation
				Config:  testAccAzureRMSharedImageVersion_setup(data, userName, password, hostName),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
				),
			},
			{
				Config: testAccAzureRMSharedImageVersion_imageVersion(data, userName, password, hostName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageVersionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_image_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_region.#", "1"),
				),
			},
			{
				Config:      testAccAzureRMSharedImageVersion_requiresImport(data, userName, password, hostName),
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
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.GalleryImageVersionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func revokeSASOnSnapshotForCleanup(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.SnapshotsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Get the snapshot resource
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		snapShotName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		future, err := client.RevokeAccess(ctx, resourceGroup, snapShotName)
		if err != nil {
			return fmt.Errorf("bad: cannot revoke SAS on the snapshot: %+v", err)
		}
		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("bad: waiting the revoke of SAS on the snapshot: %+v", err)
		}

		return nil
	}
}

// nolint: unparam
func testAccAzureRMSharedImageVersion_setup(data acceptance.TestData, username, password, hostname string) string {
	return testAccAzureRMImage_standaloneImage_setup(data, username, password, hostname, "LRS")
}

func testAccAzureRMSharedImageVersion_provision(data acceptance.TestData, username, password, hostname string) string {
	template := testAccAzureRMImage_standaloneImage_provision(data, username, password, hostname, "LRS", "")
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMSharedImageVersion_imageVersion(data acceptance.TestData, username, password, hostname string) string {
	template := testAccAzureRMSharedImageVersion_provision(data, username, password, hostname)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }

  tags = {
    "foo" = "bar"
  }
}
`, template)
}

func testAccAzureRMSharedImageVersion_setupSpecialized(data acceptance.TestData, username, password, hostname string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name               = "acctsub-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  domain_name_label   = "%s"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_virtual_machine" "testsource" {
  name                  = "testsource"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.testsource.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "mdimagetestsource"
    admin_username = "%s"
    admin_password = "%s"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, hostname, data.RandomInteger, username, password)
}

func testAccAzureRMSharedImageVersion_provisionSpecialized(data acceptance.TestData, username, password, hostname string) string {
	template := testAccAzureRMSharedImageVersion_setupSpecialized(data, username, password, hostname)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  specialized         = true

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMSharedImageVersion_imageVersionSpecializedBySnapshot(data acceptance.TestData, username, password, hostname string) string {
	template := testAccAzureRMSharedImageVersion_provisionSpecialized(data, username, password, hostname)
	return fmt.Sprintf(`
%s

resource "azurerm_snapshot" "test" {
  name                = "acctestsnapshot%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_resource_id  = azurerm_virtual_machine.testsource.storage_os_disk.0.managed_disk_id
}

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  os_disk_snapshot_id = azurerm_snapshot.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }

  tags = {
    "foo" = "bar"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMSharedImageVersion_imageVersionSpecializedByVM(data acceptance.TestData, username, password, hostname string) string {
	template := testAccAzureRMSharedImageVersion_provisionSpecialized(data, username, password, hostname)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  managed_image_id = azurerm_virtual_machine.testsource.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }

  tags = {
    "foo" = "bar"
  }
}
`, template)
}

func testAccAzureRMSharedImageVersion_imageVersionStorageAccountType(data acceptance.TestData, username, password, hostname string, storageAccountType string) string {
	template := testAccAzureRMSharedImageVersion_provision(data, username, password, hostname)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
    storage_account_type   = "%s"
  }
}
`, template, storageAccountType)
}

func testAccAzureRMSharedImageVersion_requiresImport(data acceptance.TestData, username, password, hostname string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "import" {
  name                = azurerm_shared_image_version.test.name
  gallery_name        = azurerm_shared_image_version.test.gallery_name
  image_name          = azurerm_shared_image_version.test.image_name
  resource_group_name = azurerm_shared_image_version.test.resource_group_name
  location            = azurerm_shared_image_version.test.location
  managed_image_id    = azurerm_shared_image_version.test.managed_image_id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }
}
`, testAccAzureRMSharedImageVersion_imageVersion(data, username, password, hostname))
}

func testAccAzureRMSharedImageVersion_imageVersionUpdated(data acceptance.TestData, username, password, hostname string) string {
	template := testAccAzureRMSharedImageVersion_provision(data, username, password, hostname)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "1234567890.1234567890.1234567890"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }

  target_region {
    name                   = "%s"
    regional_replica_count = 2
  }
}
`, template, data.Locations.Secondary)
}
