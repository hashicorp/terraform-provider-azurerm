package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
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
		_, errors := storage.ValidateArmStorageQueueName(tc.Value, "azurerm_storage_queue")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the ARM Storage Queue Name to trigger a validation error")
		}
	}
}

func TestAccAzureRMStorageQueue_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageQueue_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageQueueExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageQueue_basicAzureADAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageQueue_basicAzureADAuth(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageQueueExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageQueue_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageQueue_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageQueueExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStorageQueue_requiresImport),
		},
	})
}

func TestAccAzureRMStorageQueue_metaData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageQueue_metaData(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageQueueExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageQueue_metaDataUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageQueueExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMStorageQueueExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Queue %q: %s", accountName, name, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
		}

		queuesClient, err := storageClient.QueuesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Queues Client: %s", err)
		}

		metaData, err := queuesClient.GetMetaData(ctx, accountName, name)
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
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
	storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_queue" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Queue %q: %s", accountName, name, err)
		}
		// expected if this has been deleted
		if account == nil {
			return nil
		}

		queuesClient, err := storageClient.QueuesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Queues Client: %s", err)
		}

		props, err := queuesClient.GetMetaData(ctx, accountName, name)
		if err != nil {
			return nil
		}

		return fmt.Errorf("Queue still exists: %+v", props)
	}

	return nil
}

func testAccAzureRMStorageQueue_basic(data acceptance.TestData) string {
	template := testAccAzureRMStorageQueue_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%d"
  storage_account_name = "${azurerm_storage_account.test.name}"
}
`, template, data.RandomInteger)
}

func testAccAzureRMStorageQueue_basicAzureADAuth(data acceptance.TestData) string {
	template := testAccAzureRMStorageQueue_basic(data)
	return fmt.Sprintf(`
provider "azurerm" {
  storage_use_azuread = true
}

%s
`, template)
}

func testAccAzureRMStorageQueue_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStorageQueue_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "import" {
  name                 = "${azurerm_storage_queue.test.name}"
  storage_account_name = "${azurerm_storage_queue.test.storage_account_name}"
}
`, template)
}

func testAccAzureRMStorageQueue_metaData(data acceptance.TestData) string {
	template := testAccAzureRMStorageQueue_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue-%d"
  storage_account_name = "${azurerm_storage_account.test.name}"

  metadata = {
    hello = "world"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMStorageQueue_metaDataUpdated(data acceptance.TestData) string {
	template := testAccAzureRMStorageQueue_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMStorageQueue_template(data acceptance.TestData) string {
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

`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
