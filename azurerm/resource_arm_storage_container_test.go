package azurerm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
	if !features.ShouldResourcesBeImported() {
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

func TestAccAzureRMStorageContainer_metaData(t *testing.T) {
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
				Config: testAccAzureRMStorageContainer_metaData(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMStorageContainer_metaDataUpdated(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMStorageContainer_metaDataEmpty(ri, rs, location),
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

func TestAccAzureRMStorageContainer_disappears(t *testing.T) {
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
					testAccARMStorageContainerDisappears(resourceName),
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

func TestAccAzureRMStorageContainer_web(t *testing.T) {
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
				Config: testAccAzureRMStorageContainer_web(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "$web"),
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

		storageClient := testAccProvider.Meta().(*ArmClient).Storage
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		containerName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Container %q (Account %s): %s", containerName, accountName, err)
		}
		if resourceGroup == nil {
			return fmt.Errorf("Unable to locate Resource Group for Storage Container %q (Account %s) - assuming removed & removing from state", containerName, accountName)
		}

		client, err := storageClient.ContainersClient(ctx, *resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Error building Containers Client: %s", err)
		}

		resp, err := client.GetProperties(ctx, accountName, containerName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Container %q (Account %q / Resource Group %q) does not exist", containerName, accountName, *resourceGroup)
			}

			return fmt.Errorf("Bad: Get on ContainersClient: %+v", err)
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

		storageClient := testAccProvider.Meta().(*ArmClient).Storage
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		containerName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Container %q (Account %s): %s", containerName, accountName, err)
		}
		if resourceGroup == nil {
			return fmt.Errorf("Unable to locate Resource Group for Storage Container %q (Account %s) - assuming removed & removing from state", containerName, accountName)
		}

		client, err := storageClient.ContainersClient(ctx, *resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Error building Containers Client: %s", err)
		}

		if _, err := client.Delete(ctx, accountName, containerName); err != nil {
			return fmt.Errorf("Error deleting Container %q (Account %q): %s", containerName, accountName, err)
		}

		return nil
	}
}

func testCheckAzureRMStorageContainerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_container" {
			continue
		}

		storageClient := testAccProvider.Meta().(*ArmClient).Storage
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		containerName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Container %q (Account %s): %s", containerName, accountName, err)
		}

		if resourceGroup == nil {
			return nil
		}

		client, err := storageClient.ContainersClient(ctx, *resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Error building Containers Client: %s", err)
		}

		props, err := client.GetProperties(ctx, accountName, containerName)
		if err != nil {
			return nil
		}

		return fmt.Errorf("Container still exists: %+v", props)
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

func testAccAzureRMStorageContainer_metaData(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageContainer_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"

  metadata = {
    hello = "world"
  }
}
`, template)
}

func testAccAzureRMStorageContainer_metaDataUpdated(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageContainer_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"

  metadata = {
    hello = "world"
    panda = "pops"
  }
}
`, template)
}

func testAccAzureRMStorageContainer_metaDataEmpty(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageContainer_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"

  metadata = {
  }
}
`, template)
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

func testAccAzureRMStorageContainer_web(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageContainer_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "$web"
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
		"$web",
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
		"$notweb",
		strings.Repeat("w", 65),
	}
	for _, v := range invalidNames {
		_, errors := validateArmStorageContainerName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Storage Container Name", v)
		}
	}
}
