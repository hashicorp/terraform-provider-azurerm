package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func init() {
	resource.AddTestSweepers("azurerm_monitor_log_profile", &resource.Sweeper{
		Name: "azurerm_monitor_log_profile",
		F:    testSweepMonitorLogProfiles,
	})
}

func testSweepMonitorLogProfiles(region string) error {
	armClient, err := buildConfigForSweepers()
	if err != nil {
		return fmt.Errorf("Error building config for sweepers: %+v", err)
	}

	client := (*armClient).monitorLogProfilesClient
	ctx := (*armClient).StopContext

	results, err := client.List(ctx)
	if err != nil {
		return fmt.Errorf("Error Listing on Log Profiles: %+v", err)
	}

	for _, logProfile := range *results.Value {
		name := *logProfile.Name
		// Use region as resource location as there's no location associated with a log profile
		if !shouldSweepAcceptanceTestResource(name, region, region) {
			continue
		}

		_, err := client.Delete(ctx, name)
		if err != nil {
			return fmt.Errorf("Error deleting Log Profile %q: %+v", name, err)
		}
	}

	return nil
}

// NOTE: this is a combined test rather than separate split out tests due to
// Azure only being happy about provisioning one per subscription at once
// (which our test suite can't easily workaround)

// this occasionally fails due to the rapid provisioning and deprovisioning,
// running the exact same test afterwards always results in a pass.

func TestAccAzureRMMonitorLogProfile(t *testing.T) {
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"basic":      testAccAzureRMMonitorLogProfile_basic,
			"servicebus": testAccAzureRMMonitorLogProfile_servicebus,
			"complete":   testAccAzureRMMonitorLogProfile_complete,
			"disappears": testAccAzureRMMonitorLogProfile_disappears,
		},
		"datasource": {
			"eventhub":       testAccDataSourceAzureRMMonitorLogProfile_eventhub,
			"storageaccount": testAccDataSourceAzureRMMonitorLogProfile_storageaccount,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccAzureRMMonitorLogProfile_basic(t *testing.T) {
	resourceName := "azurerm_monitor_log_profile.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorLogProfile_basicConfig(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogProfileExists(resourceName),
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

func testAccAzureRMMonitorLogProfile_servicebus(t *testing.T) {
	resourceName := "azurerm_monitor_log_profile.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorLogProfile_servicebusConfig(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogProfileExists(resourceName),
				),
			},
		},
	})
}

func testAccAzureRMMonitorLogProfile_complete(t *testing.T) {
	resourceName := "azurerm_monitor_log_profile.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorLogProfile_completeConfig(ri, rs, testLocation(), testAltLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogProfileExists(resourceName),
				),
			},
		},
	})
}

func testAccAzureRMMonitorLogProfile_disappears(t *testing.T) {
	resourceName := "azurerm_monitor_log_profile.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(10)
	config := testAccAzureRMMonitorLogProfile_basicConfig(ri, rs, testLocation())

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
	client := testAccProvider.Meta().(*ArmClient).monitorLogProfilesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_log_profile" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resp, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Log Profile still exists:\n%#v", *resp.ID)
			}
		}
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

		client := testAccProvider.Meta().(*ArmClient).monitorLogProfilesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		name := rs.Primary.Attributes["name"]
		resp, err := client.Get(ctx, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Log Profile %q does not exist", name)
			}

			return fmt.Errorf("Bad: Get on monitorLogProfilesClient: %+v", err)
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

		client := testAccProvider.Meta().(*ArmClient).monitorLogProfilesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, err := client.Delete(ctx, name)
		if err != nil {
			return fmt.Errorf("Error deleting Log Profile %q: %+v", name, err)
		}

		return nil
	}
}

func testAccAzureRMMonitorLogProfile_basicConfig(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name     = "acctestrg-%d"
    location = "%s"
}

resource "azurerm_storage_account" "test" {
    name                     = "acctestsa%s"
    resource_group_name      = "${azurerm_resource_group.test.name}"
    location                 = "${azurerm_resource_group.test.location}"
    account_tier             = "Standard"
    account_replication_type = "GRS"
}
    
resource "azurerm_monitor_log_profile" "test" {
    name = "acctestlp-%d"

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
`, rInt, location, rString, rInt, location)
}

func testAccAzureRMMonitorLogProfile_servicebusConfig(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name     = "acctestrg-%d"
    location = "%s"
}
    
resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestsbns-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "standard"
}

resource "azurerm_servicebus_namespace_authorization_rule" "test" {
  name                = "acctestsbrule-%s"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  listen = true
  send   = true
  manage = true
}

resource "azurerm_monitor_log_profile" "test" {
    name = "acctestlp-%d"

    categories = [
        "Action"
    ]
    
    locations = [
        "%s",
    ]
    
    servicebus_rule_id = "${azurerm_servicebus_namespace_authorization_rule.test.id}"
    
    retention_policy {
        enabled = false
    } 
}
`, rInt, location, rString, rString, rInt, location)
}

func testAccAzureRMMonitorLogProfile_completeConfig(rInt int, rString string, location string, altLocation string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name     = "acctestrg-%d"
    location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
    
resource "azurerm_eventhub_namespace" "test" {
  name                = "acctestehns-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  capacity            = 2
}

resource "azurerm_monitor_log_profile" "test" {
    name = "acctestlp-%d"

    categories = [
        "Action",
        "Delete",
        "Write",
    ]
    
    locations = [
        "%s",
        "%s",
    ]
    
    # RootManageSharedAccessKey is created by default with listen, send, manage permissions
    servicebus_rule_id = "${azurerm_eventhub_namespace.test.id}/authorizationrules/RootManageSharedAccessKey"
    storage_account_id  = "${azurerm_storage_account.test.id}"

    retention_policy {
        enabled = true
        days = 7
    }
}
`, rInt, location, rString, rString, rInt, location, altLocation)
}
