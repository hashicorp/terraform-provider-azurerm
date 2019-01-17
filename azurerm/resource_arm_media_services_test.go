package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
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

func TestAccAzureRMMediaServicesAccount_basiccreation(t *testing.T) {

	t.Log("Running my test!")
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_media_services_test"

	template := testAccAzureRMMediaServicesAccount_basic(ri, testLocation())

	t.Log(template)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMediaServicesAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMediaServicesAccount_basic(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMMediaServicesAccount_basic(resourceName, testLocation()),
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

func testAccAzureRMMediaServicesAccount_basic(rInt int, location string) string {
	return fmt.Sprintf(`
	resource "azurerm_resource_group" "testrg" {
		name     = "mstest-%[1]d"
		location = "%s"
	  }
	  
	  resource "azurerm_storage_account" "testsa" {
		name                     = "mstest%[1]d"
		resource_group_name      = "${azurerm_resource_group.testrg.name}"
		location                 = "${azurerm_resource_group.testrg.location}"
		account_tier             = "Standard"
		account_replication_type = "GRS"
	  
		tags {
		  environment = "staging"
		}
	  }
	  
	  resource "azurerm_media_services" "ams" {
	  
			  name = "mstest%[1]d"
			  location = "${azurerm_resource_group.testrg.location}"
			  resource_group_name = "${azurerm_resource_group.testrg.name}"
			  
			  storage_account {
					  id = "${azurerm_storage_account.testsa.id}"
					  is_primary = true
			  }
	  }
`, rInt, location)
}

func testCheckAzureRMMediaServicesAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
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

func checkAccAzureRMMediaServicesAccount_basic(resourceName string, location string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testCheckAzureRMMediaServicesAccountExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttr(resourceName, "location", azureRMNormalizeLocation(location)),
	)
}
