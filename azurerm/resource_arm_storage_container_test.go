package azurerm

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMStorageContainer_basic(t *testing.T) {
	resourceName := "azurerm_storage_container.test"

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(resourceName),
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

func TestAccAzureRMStorageContainer_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_storage_container.test"

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStorageContainer_requiresImport(ri, rs, location),
				ExpectError: testRequiresImportError("azurerm_storage_container"),
			},
		},
	})
}

func TestAccAzureRMStorageContainer_update(t *testing.T) {
	resourceName := "azurerm_storage_container.test"

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_update(ri, rs, location, "private"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container_access_type", "private"),
				),
			},
			{
				Config: testAccAzureRMStorageContainer_update(ri, rs, location, "container"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container_access_type", "container"),
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

func TestAccAzureRMStorageContainer_disappears(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists("azurerm_storage_container.test"),
					testAccARMStorageContainerDisappears("azurerm_storage_container.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMStorageContainer_root(t *testing.T) {
	resourceName := "azurerm_storage_container.test"

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_root(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "$root"),
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

func testCheckAzureRMStorageContainerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		storageAccountName := rs.Primary.Attributes["storage_account_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for storage container: %s", name)
		}

		armClient := testAccProvider.Meta().(*ArmClient)
		ctx := armClient.StopContext
		blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, resourceGroup, storageAccountName)
		if err != nil {
			return err
		}
		if !accountExists {
			return fmt.Errorf("Bad: Storage Account %q does not exist", storageAccountName)
		}

		containers, err := blobClient.ListContainers(storage.ListContainersParameters{
			Prefix:  name,
			Timeout: 90,
		})
		if err != nil {
			return fmt.Errorf("Error listing Storage Container %q containers (storage account: %q) : %+v", name, storageAccountName, err)
		}

		if len(containers.Containers) == 0 {
			return fmt.Errorf("Bad: Storage Container %q (storage account: %q) does not exist", name, storageAccountName)
		}

		var found bool
		for _, container := range containers.Containers {
			if container.Name == name {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("Bad: Storage Container %q (storage account: %q) does not exist", name, storageAccountName)
		}

		return nil
	}
}

func testAccARMStorageContainerDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		armClient := testAccProvider.Meta().(*ArmClient)
		ctx := armClient.StopContext

		storageAccountName := rs.Primary.Attributes["storage_account_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		name := rs.Primary.Attributes["name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for storage container: %s", name)
		}

		blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, resourceGroup, storageAccountName)
		if err != nil {
			return err
		}
		if !accountExists {
			log.Printf("[INFO] Storage Account %q doesn't exist so the container won't exist", storageAccountName)
			return nil
		}

		reference := blobClient.GetContainerReference(name)
		options := &storage.DeleteContainerOptions{}
		_, err = reference.DeleteIfExists(options)
		return err
	}
}

func testCheckAzureRMStorageContainerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_container" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		storageAccountName := rs.Primary.Attributes["storage_account_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for storage container: %s", name)
		}

		armClient := testAccProvider.Meta().(*ArmClient)
		ctx := armClient.StopContext
		blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, resourceGroup, storageAccountName)
		if err != nil {
			//If we can't get keys then the blob can't exist
			return nil
		}
		if !accountExists {
			return nil
		}

		containers, err := blobClient.ListContainers(storage.ListContainersParameters{
			Prefix:  name,
			Timeout: 90,
		})

		if err != nil {
			return nil
		}

		var found bool
		for _, container := range containers.Containers {
			if container.Name == name {
				found = true
			}
		}

		if found {
			return fmt.Errorf("Bad: Storage Container %q (storage account: %q) still exist", name, storageAccountName)
		}
	}

	return nil
}

func testAccAzureRMStorageContainer_basic(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageContainer_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}
`, template)
}

func testAccAzureRMStorageContainer_requiresImport(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageContainer_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "import" {
  name                  = "${azurerm_storage_container.test.name}"
  resource_group_name   = "${azurerm_storage_container.test.resource_group_name}"
  storage_account_name  = "${azurerm_storage_container.test.storage_account_name}"
  container_access_type = "${azurerm_storage_container.test.container_access_type}"
}
`, template)
}

func testAccAzureRMStorageContainer_update(rInt int, rString string, location string, accessType string) string {
	template := testAccAzureRMStorageContainer_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "%s"
}
`, template, accessType)
}

func testAccAzureRMStorageContainer_root(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageContainer_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "$root"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}
`, template)
}

func testAccAzureRMStorageContainer_template(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rString)
}

func TestValidateArmStorageContainerName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
		"$root",
	}
	for _, v := range validNames {
		_, errors := validateArmStorageContainerName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Storage Container Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"InvalidName1",
		"-invalidname1",
		"invalid_name",
		"invalid!",
		"ww",
		"$notroot",
		strings.Repeat("w", 65),
	}
	for _, v := range invalidNames {
		_, errors := validateArmStorageContainerName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Storage Container Name", v)
		}
	}
}
