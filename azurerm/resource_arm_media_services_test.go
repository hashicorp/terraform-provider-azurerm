package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2018-07-01/media"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func init() {
	resource.AddTestSweepers("azurerm_media_services", &resource.Sweeper{
		Name: "azurerm_media_services",
		F:    testSweepMediaServicesAccount,
	})
}

func testSweepMediaServicesAccount(region string) error {
	armClient, err := buildConfigForSweepers()
	if err != nil {
		return err
	}

	client := (*armClient).mediaServicesClient
	ctx := (*armClient).StopContext

	log.Printf("Retrieving the Media Services Accounts..")
	results, err := client.List(ctx, "testrg")
	if err != nil {
		return fmt.Errorf("Error Listing on Media Services Accounts: %+v", err)
	}

	for _, account := range results.Values() {
		if !shouldSweepAcceptanceTestResource(*account.Name, *account.Location, region) {
			continue
		}

		resourceId, err := parseAzureResourceID(*account.ID)
		if err != nil {
			return err
		}

		resourceGroup := resourceId.ResourceGroup
		name := resourceId.Path["mediaservices"]

		log.Printf("Deleting Media Services Account '%s' in Resource Group '%s'", name, resourceGroup)
		_, err = client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestAccAzureRMMediaServicesAccount_singlestorage(t *testing.T) {

	ri := tf.AccRandTimeInt()
	resourceName := fmt.Sprintf("%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMediaServicesAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: generateTemplate_singlestorage(resourceName, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMMediaServicesAccount("azurerm_media_services.ams", testLocation(), "1"),
				),
			},
		},
	})
}

func generateTemplate_singlestorage(name string, location string) string {
	return fmt.Sprintf(`
	resource "azurerm_resource_group" "testrg" {
			name     = "%[1]s"
			location = "%[2]s"
	  }
	  
	  resource "azurerm_storage_account" "testsa" {
			name                     = "%[1]s"
			resource_group_name      = "${azurerm_resource_group.testrg.name}"
			location                 = "${azurerm_resource_group.testrg.location}"
			account_tier             = "Standard"
			account_replication_type = "GRS"
		
			tags {
			environment = "staging"
			}
	  }
	  
	  resource "azurerm_media_services" "ams" {
	  
			  name = "%[1]s"
			  location = "${azurerm_resource_group.testrg.location}"
			  resource_group_name = "${azurerm_resource_group.testrg.name}"
			  
			  storage_account {
					  id = "${azurerm_storage_account.testsa.id}"
					  is_primary = true
			  }

			  tags {
				environment = "staging"
			  }
	  }
`, name, location)
}

func TestAccAzureRMMediaServicesAccount_multiplestorage(t *testing.T) {

	ri := tf.AccRandTimeInt()
	resourceName := fmt.Sprintf("%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMediaServicesAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: generateTemplate_multiplestorage(resourceName, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMMediaServicesAccount("azurerm_media_services.ams", testLocation(), "2"),
				),
			},
		},
	})
}

func generateTemplate_multiplestorage(name string, location string) string {
	return fmt.Sprintf(`
	resource "azurerm_resource_group" "testrg" {
			name     = "%[1]s"
			location = "%[2]s"
	  }
	  
	  resource "azurerm_storage_account" "testsa1" {
			name                     = "1%[1]s"
			resource_group_name      = "${azurerm_resource_group.testrg.name}"
			location                 = "${azurerm_resource_group.testrg.location}"
			account_tier             = "Standard"
			account_replication_type = "GRS"
		
			tags {
			environment = "staging"
			}
	  }

	  resource "azurerm_storage_account" "testsa2" {
		name                     = "2%[1]s"
		resource_group_name      = "${azurerm_resource_group.testrg.name}"
		location                 = "${azurerm_resource_group.testrg.location}"
		account_tier             = "Standard"
		account_replication_type = "GRS"
	
		tags {
		environment = "staging"
		}
  }
	  
	  resource "azurerm_media_services" "ams" {
	  
			  name = "%[1]s"
			  location = "${azurerm_resource_group.testrg.location}"
			  resource_group_name = "${azurerm_resource_group.testrg.name}"
			  
			  storage_account {
					  id = "${azurerm_storage_account.testsa1.id}"
					  is_primary = true
			  }

			  storage_account {
					id = "${azurerm_storage_account.testsa2.id}"
					is_primary = false
			  }

			  tags {
					environment = "staging"
			  }
	  }
`, name, location)
}

func TestAzureRMMediaServicesAccount_validateStorageConfiguration(t *testing.T) {

	accounts := []media.StorageAccount{
		media.StorageAccount{
			ID:   utils.String("id1"),
			Type: media.Primary,
		},
		media.StorageAccount{
			ID:   utils.String("id2"),
			Type: media.Primary,
		},
	}

	actualErr := validateStorageConfiguration(accounts)
	expectedErr := fmt.Errorf("Error processing storage account 'id2'. Another storage account is already assigned as is_primary = 'true'")

	if actualErr == nil && expectedErr.Error() != actualErr.Error() {
		t.Fatal("validateStorageConfiguration must throw an error when more than one storage account is Primary.")
	}

	accounts = []media.StorageAccount{
		media.StorageAccount{
			ID:   utils.String("id1"),
			Type: media.Primary,
		},
		media.StorageAccount{
			ID:   utils.String("id2"),
			Type: media.Secondary,
		},
	}

	actualErr = validateStorageConfiguration(accounts)

	if actualErr != nil {
		t.Fatal("validateStorageConfiguration must allow multiple storage accounts when only one is Primary.")
	}
}

func testCheckAzureRMMediaServicesAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Media service not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Media Services Account: '%s'", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).mediaServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on mediaServicesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Media Services Account '%s' (resource group: '%s') does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMMediaServicesAccountDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).mediaServicesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_media_services" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Media Services Account still exists:\n%#v", resp)
		}
	}

	return nil
}

func checkAccAzureRMMediaServicesAccount(resourceName string, location string, storageCount string) resource.TestCheckFunc {
	// It would be ideal to also validate which storage account was Primary, but that isn't straight forward.
	// The key for the storage account's Primary flag is non-deterministic (ex: storage_account.1137153885.is_primary) and
	// isn't based on the storage account name.
	return resource.ComposeTestCheckFunc(
		testCheckAzureRMMediaServicesAccountExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttr(resourceName, "location", azureRMNormalizeLocation(location)),
		resource.TestCheckResourceAttr(resourceName, "storage_account.#", storageCount),
	)
}
