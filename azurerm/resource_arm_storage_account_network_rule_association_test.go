package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMStorageAccountNetworkRulesAssociation_basic(t *testing.T) {
	resourceName := "azurerm_storage_account_network_rules_assocation.testsanra"
	ri := tf.AccRandTimeINt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccountNetworkRulesAssociation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testAccAzureRMStorageAccountNetworkRulesAssociationExists(resourceName),
				),			
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			}
		},
	})
}

func TestAccAzureRMStorageAccountNetworkRulesAssociation_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_storage_account_network_rules_assocation.testsanra"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccountNetworkRulesAssociation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountNetworkRulesAssociationExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStorageAccountNetworkRulesAssociation_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_storage_account_network_rules_assocation"),
			},
		},
	})
}

func TestAccAzureRMStorageAccountNetworkRulesAssociation_deleted(t *testing.T) {
	resourceName := "azurerm_storage_account_network_rules_assocation.testsanra"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccountNetworkRulesAssociation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountNetworkRulesAssociationExists(resourceName),
					testCheckAzureRMStorageAccountNetworkRulesAssociationDisappears(resourceName),
					testCheckAzureRMStorageAccountHasNoNetworkRules(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMStorageAccountNetworkRulesAssociationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		storageAccountId := rs.Primary.Attributes["storage_account_id"]
		parsedId, err := parseAzureResourceID(storageAccountId)
		if err != nil {
			return err
		}

		resourceGroupName := parsedId.ResourceGroup
		storageAccountName := parsedId.Path["storage_account"]

		client := testAccProvider.Meta().(*ArmClient).storageAccountClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroupName, storageAccountName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Storage Account %q (Resource Group: %q) does not exist", storageAccountName, resourceGroupName)
			}

			return fmt.Errorf("Bad: Get on storageAccountClient: %+v", err)
		}

		props := resp.StorageAccountPropertiesFormat
		if props == nil {
			return fmt.Errorf("Properties was nil for Storage Account %q (Resource Group: %q)", storageAccountName, resourceGroupName)
		}

		if props.NetworkRules == nil || props.NetworkRules.ID == nil {
			return fmt.Errorf("No Network Rules association exists for Storage Account %q (Resource Group: %q)", storageAccountName, resourceGroupName)
		}

		return nil
	}
}

func testCheckAzureRMStorageAccountNetworkRulesAssociationDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		storageAccountId := rs.Primary.Attributes["storage_account_id"]
		parsedId, err := parseAzureResourceID(storageAccountId)
		if err != nil {
			return err
		}

		resourceGroup := parsedId.ResourceGroup
		storageAccountName := parsedId.Path["storage_account"]

		client := testAccProvider.Meta().(*ArmClient).storageAccountClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		read, err := client.Get(ctx, resourceGroup, storageAccountName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(read.Response) {
				return fmt.Errorf("Error retrieving Storage Account %q (Resource Group %q): %+v", storageAccountName, resourceGroup, err)
			}
		}

		read.StorageAccountPropertiesFormat.NetworkRules = nil

		future, err := client.CreateOrUpdate(ctx, resourceGroup, storageAccountName, read)
		if err != nil {
			return fmt.Errorf("Error updating Storage Account %q (Resource Group %q): %+v", storageAccountName, resourceGroup, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for completion of Storage Account %q (Resource Group %q): %+v", storageAccountName, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMStorageAccountHasNoNetworkRules(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		storageAccountId := rs.Primary.Attributes["storage_account_id"]
		parsedId, err := parseAzureResourceID(storageAccountId)
		if err != nil {
			return err
		}

		resourceGroupName := parsedId.ResourceGroup
		storageAccountName := parsedId.Path["storage_account"]

		client := testAccProvider.Meta().(*ArmClient).storageAccountClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroupName, storageAccountName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Storage Account %q (Resource Group: %q) does not exist", storageAccountName, resourceGroupName)
			}

			return fmt.Errorf("Bad: Get on storageAccountClient: %+v", err)
		}

		props := resp.StorageAccountPropertiesFormat
		if props == nil {
			return fmt.Errorf("Properties was nil for Storage Account %q (Resource Group: %q)", storageAccountName, resourceGroupName)
		}

		if props.NetworkRules != nil && ((props.NetworkRules.ID == nil) || (props.NetworkRules.ID != nil && *props.NetworkRules.ID == "")) {
			return fmt.Errorf("No Network Rules should exist for Storage Account %q (Resource Group: %q) but got %q", storageAccountName, resourceGroupName, *props.NetworkRules.ID)
		}

		return nil
	}
}

func testAccAzureRMStorageAccountNetworkRulesAssociation_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.testrg.location}"
  resource_group_name = "${azurerm_resource_group.testrg.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.testrg.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "testsa" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = "${azurerm_resource_group.testrg.name}"
  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

resource "azurerm_network_rules" testnr" {
	name                     = "acctestnr%d"
  resource_group_name      = "${azurerm_resource_group.testrg.name}"
	location                 = "${azurerm_resource_group.testrg.location}"
	
	network_rules {
		default_action             = "Deny"
    ip_rules                   = ["127.0.0.1"]
    virtual_network_subnet_ids = ["${azurerm_subnet.test.id}"]
	}
}

resource "azurerm_storage_account_network_rules_association" "testsanra" {
	storage_account_id        = "${azurerm_storage_account_testsa.id}"
	network_rules_id          = "${azurerm_network_rules.testnr.id}"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMStorageAccountNetworkRulesAssociation_requiresImport(rInt int, location string) string {
	template := testAccAzureRMStorageAccountNetworkRulesAssociation_basic(rInt, location)
	return fmt.Sprintf(`
%s

`, template)
}
