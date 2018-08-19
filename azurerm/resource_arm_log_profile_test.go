package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMLogProfile_basic(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogProfile_basic(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogProfileExists("azurerm_log_profile.test"),
				),
			},
		},
	})
}

func TestAccAzureRMLogProfile_servicebus(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogProfile_servicebus(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogProfileExists("azurerm_log_profile.test"),
				),
			},
		},
	})
}

func TestAccAzureRMLogProfile_complete(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogProfile_complete(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogProfileExists("azurerm_log_profile.test"),
				),
			},
		},
	})
}

func TestAccAzureRMLogProfile_disappears(t *testing.T) {
	resourceName := "azurerm_log_profile.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(10)
	config := testAccAzureRMLogProfile_basic(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogProfileExists(resourceName),
					testCheckAzureRMLogProfileDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMLogProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).logProfilesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_profile" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resp, err := client.Get(ctx, name)
		if err != nil {
			return nil
		}

		return fmt.Errorf("Log Profile still exists:\n%#v", *resp.ID)
	}

	return nil
}

func testCheckAzureRMLogProfileExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).logProfilesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		name := rs.Primary.Attributes["name"]
		resp, err := client.Get(ctx, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Log Profile %q does not exist", name)
			}

			return fmt.Errorf("Bad: Get on logProfilesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMLogProfileDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).logProfilesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, err := client.Delete(ctx, name)
		if err != nil {
			return fmt.Errorf("Error deleting Log Profile %q: %+v", name, err)
		}

		return nil
	}
}

func testAccAzureRMLogProfile_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctest%d-rg"
	location = "%s"
}

resource "azurerm_storage_account" "test" {
	name                     = "%s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "GRS"
}
	
resource "azurerm_log_profile" "test" {
	name = "basic"

	categories = [
		"Action",
	]
	
	locations = [
		"%s"
	]
	
	storage_account_id  = "${azurerm_storage_account.test.id}"

	retention_policy {
		enabled = true
		days = 7
	}
}
`, rInt, location, rString, location)
}

func testAccAzureRMLogProfile_servicebus(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctest%d-rg"
	location = "%s"
}
	
resource "azurerm_servicebus_namespace" "test" {
  name                = "a%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "standard"
}

resource "azurerm_servicebus_namespace_authorization_rule" "test" {
  name                = "%s"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  listen = true
  send   = true
  manage = true
}

resource "azurerm_log_profile" "test" {
	name = "default"

	categories = [
		"Action"
	]
	
	locations = [
		"%s",
	]
	
	service_bus_rule_id = "${azurerm_servicebus_namespace_authorization_rule.test.id}"
	
	retention_policy {
		enabled = false
		days = 0
	}
}
`, rInt, location, rString, rString, location)
}

func testAccAzureRMLogProfile_complete(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctest%d-rg"
	location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
	
resource "azurerm_eventhub_namespace" "test" {
  name                = "a%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  capacity            = 2
}

resource "azurerm_log_profile" "test" {
	name = "complete"

	categories = [
		"Action",
		"Delete",
		"Write",
	]
	
	locations = [
		"east asia",
		"southeastasia",
		"centralus",
		"eastus",
		"eastus2",
		"westus",
		"northcentralus",
		"southcentralus",
		"northeurope",
		"westeurope",
		"japan west",
		"japaneast",
		"brazilsouth",
		"australiaeast",
		"australiasoutheast",
		"southindia",
		"centralindia",
		"west india",
		"canadacentral",
		"canadaeast",
		"uksouth",
		"ukwest",
		"westcentralus",
		"westus2",
		"koreacentral",
		"koreasouth",
		"francecentral",
		"francesouth",
		"australiacentral",
		"australiacentral2",
		"global",
	]
	
	# RootManageSharedAccessKey is created by default with listen, send, manage permissions
	service_bus_rule_id = "${azurerm_eventhub_namespace.test.id}/authorizationrules/RootManageSharedAccessKey"
	storage_account_id  = "${azurerm_storage_account.test.id}"

	retention_policy {
		enabled = true
		days = 7
	}
}
`, rInt, location, rString, rString)
}
