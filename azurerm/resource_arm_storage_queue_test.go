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

func TestResourceAzureRMStorageQueueName_Validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "testing_123",
			ErrCount: 1,
		},
		{
			Value:    "testing123-",
			ErrCount: 1,
		},
		{
			Value:    "-testing123",
			ErrCount: 1,
		},
		{
			Value:    "TestingSG",
			ErrCount: 1,
		},
		{
			Value:    acctest.RandString(256),
			ErrCount: 1,
		},
		{
			Value:    acctest.RandString(1),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateArmStorageQueueName(tc.Value, "azurerm_storage_queue")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the ARM Storage Queue Name to trigger a validation error")
		}
	}
}

func TestAccAzureRMStorageQueue_basic(t *testing.T) {
	resourceName := "azurerm_storage_queue.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageQueue_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageQueueExists(resourceName),
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

func TestAccAzureRMStorageQueue_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_storage_queue.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageQueue_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageQueueExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStorageQueue_requiresImport(ri, rs, location),
				ExpectError: testRequiresImportError("azurerm_storage_queue"),
			},
		},
	})
}

func TestAccAzureRMStorageQueue_metaData(t *testing.T) {
	resourceName := "azurerm_storage_queue.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageQueue_metaData(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageQueueExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMStorageQueue_metaDataUpdated(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageQueueExists(resourceName),
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

func testCheckAzureRMStorageQueueExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		storageClient := testAccProvider.Meta().(*ArmClient).Storage

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Queue %q (Account %s): %s", name, accountName, err)
		}
		if resourceGroup == nil {
			return fmt.Errorf("Unable to locate Resource Group for Storage Queue %q (Account %s) - assuming removed", name, accountName)
		}

		queueClient, err := storageClient.QueuesClient(ctx, *resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Error building Queues Client: %s", err)
		}

		metaData, err := queueClient.GetMetaData(ctx, accountName, name)
		if err != nil {
			if utils.ResponseWasNotFound(metaData.Response) {
				return fmt.Errorf("Bad: Storage Queue %q (storage account: %q) does not exist", name, accountName)
			}

			return fmt.Errorf("Bad: error retrieving Storage Queue %q (storage account: %q): %s", name, accountName, err)
		}

		return nil
	}
}

func testCheckAzureRMStorageQueueDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_queue" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		storageClient := testAccProvider.Meta().(*ArmClient).Storage

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Queue %q (Account %s): %s", name, accountName, err)
		}

		// expected if this has been deleted
		if resourceGroup == nil {
			return nil
		}

		queueClient, err := storageClient.QueuesClient(ctx, *resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Error building Queues Client: %s", err)
		}

		props, err := queueClient.GetMetaData(ctx, accountName, name)
		if err != nil {
			return nil
		}

		return fmt.Errorf("Queue still exists: %+v", props)
	}

	return nil
}

func testAccAzureRMStorageQueue_basic(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageQueue_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%d"
  storage_account_name = "${azurerm_storage_account.test.name}"
}
`, template, rInt)
}

func testAccAzureRMStorageQueue_requiresImport(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageQueue_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "import" {
  name                 = "${azurerm_storage_queue.test.name}"
  storage_account_name = "${azurerm_storage_queue.test.storage_account_name}"
}
`, template)
}

func testAccAzureRMStorageQueue_metaData(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageQueue_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%d"
  storage_account_name = "${azurerm_storage_account.test.name}"

  metadata = {
    hello = "world"
  }
}
`, template, rInt)
}

func testAccAzureRMStorageQueue_metaDataUpdated(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageQueue_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%d"
  storage_account_name = "${azurerm_storage_account.test.name}"

  metadata = {
    hello = "world"
    rick  = "M0rty"
  }
}
`, template, rInt)
}

func testAccAzureRMStorageQueue_template(rInt int, rString string, location string) string {
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
