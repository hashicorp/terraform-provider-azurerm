package azurerm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAnalysisServicesServer_basic(t *testing.T) {
	resourceName := "azurerm_analysis_services_server.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAnalysisServicesServer_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(resourceName),
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

func TestAccAzureRMAnalysisServicesServer_withTags(t *testing.T) {
	resourceName := "azurerm_analysis_services_server.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMAnalysisServicesServer_withTags(ri, testLocation())
	postConfig := testAccAzureRMAnalysisServicesServer_withTagsUpdate(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.label", "test"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.label", "test1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "prod"),
				),
			},
		},
	})
}

func TestAccAzureRMAnalysisServicesServer_querypoolConnectionMode(t *testing.T) {
	resourceName := "azurerm_analysis_services_server.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMAnalysisServicesServer_querypoolConnectionMode(ri, testLocation(), "All")
	postConfig := testAccAzureRMAnalysisServicesServer_querypoolConnectionMode(ri, testLocation(), "ReadOnly")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "querypool_connection_mode", "All"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "querypool_connection_mode", "ReadOnly"),
				),
			},
		},
	})
}

func TestAccAzureRMAnalysisServicesServer_firewallSettings(t *testing.T) {
	resourceName := "azurerm_analysis_services_server.test"
	ri := tf.AccRandTimeInt()

	config1 := testAccAzureRMAnalysisServicesServer_firewallSettings(ri, testLocation(), true, make([]map[string]string, 0))

	firewallRules2 := make([]map[string]string, 1)

	firewallRules2[0] = make(map[string]string)
	firewallRules2[0]["name"] = "test1"
	firewallRules2[0]["range_start"] = "92.123.234.11"
	firewallRules2[0]["range_end"] = "92.123.234.12"

	config2 := testAccAzureRMAnalysisServicesServer_firewallSettings(ri, testLocation(), false, firewallRules2)

	firewallRules3 := make([]map[string]string, 2)

	firewallRules3[0] = make(map[string]string)
	firewallRules3[0]["name"] = "test1"
	firewallRules3[0]["range_start"] = "92.123.234.11"
	firewallRules3[0]["range_end"] = "92.123.234.13"

	firewallRules3[1] = make(map[string]string)
	firewallRules3[1]["name"] = "test2"
	firewallRules3[1]["range_start"] = "226.202.187.57"
	firewallRules3[1]["range_end"] = "226.208.192.47"

	config3 := testAccAzureRMAnalysisServicesServer_firewallSettings(ri, testLocation(), true, firewallRules3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_power_bi_service", "true"),
					resource.TestCheckResourceAttr(resourceName, "ipv4_firewall_rule.#", "0"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_power_bi_service", "false"),
					resource.TestCheckResourceAttr(resourceName, "ipv4_firewall_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ipv4_firewall_rule.0.name", "test1"),
					resource.TestCheckResourceAttr(resourceName, "ipv4_firewall_rule.0.range_start", "92.123.234.11"),
					resource.TestCheckResourceAttr(resourceName, "ipv4_firewall_rule.0.range_end", "92.123.234.12"),
				),
			},
			{
				Config: config3,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_power_bi_service", "true"),
					resource.TestCheckResourceAttr(resourceName, "ipv4_firewall_rule.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ipv4_firewall_rule.0.name", "test1"),
					resource.TestCheckResourceAttr(resourceName, "ipv4_firewall_rule.0.range_start", "92.123.234.11"),
					resource.TestCheckResourceAttr(resourceName, "ipv4_firewall_rule.0.range_end", "92.123.234.13"),
					resource.TestCheckResourceAttr(resourceName, "ipv4_firewall_rule.1.name", "test2"),
					resource.TestCheckResourceAttr(resourceName, "ipv4_firewall_rule.1.range_start", "226.202.187.57"),
					resource.TestCheckResourceAttr(resourceName, "ipv4_firewall_rule.1.range_end", "226.208.192.47"),
				),
			},
		},
	})
}

// Currently unsure how to test it as the email addresses need to exist in the AD
//func TestAccAzureRMAnalysisServicesServer_adminUsers(t *testing.T) {
//	resourceName := "azurerm_analysis_services_server.test"
//	ri := tf.AccRandTimeInt()
//	preAdminUsers := []string{"admin@domain.tld"}
//	postAdminUsers := []string{"admin@domain.tld", "admin2@domain.tld"}
//	preConfig := testAccAzureRMAnalysisServicesServer_adminUsers(ri, testLocation(), preAdminUsers)
//	postConfig := testAccAzureRMAnalysisServicesServer_adminUsers(ri, testLocation(), postAdminUsers)
//
//	resource.ParallelTest(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: preConfig,
//				Check: resource.ComposeTestCheckFunc(
//					testCheckAzureRMAnalysisServicesServerExists(resourceName),
//					resource.TestCheckResourceAttr(resourceName, "admin_users.#", "1"),
//					resource.TestCheckResourceAttr(resourceName, "admin_users.0", preAdminUsers[0]),
//				),
//			},
//			{
//				Config: postConfig,
//				Check: resource.ComposeTestCheckFunc(
//					testCheckAzureRMAnalysisServicesServerExists(resourceName),
//					resource.TestCheckResourceAttr(resourceName, "admin_users.#", "2"),
//					resource.TestCheckResourceAttr(resourceName, "admin_users.0", postAdminUsers[0]),
//					resource.TestCheckResourceAttr(resourceName, "admin_users.1", postAdminUsers[1]),
//				),
//			},
//		},
//	})
//}

func testAccAzureRMAnalysisServicesServer_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku 				  = "B1"
}
`, rInt, location, rInt)
}

func testAccAzureRMAnalysisServicesServer_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku 				  = "B1"

  tags = {
	label = "test"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAnalysisServicesServer_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku 				  = "B1"

  tags = {
	label = "test1"
	env   = "prod"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAnalysisServicesServer_querypoolConnectionMode(rInt int, location, connectionMode string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                		= "acctestass%d"
  location            		= "${azurerm_resource_group.test.location}"
  resource_group_name 		= "${azurerm_resource_group.test.name}"
  sku 				  		= "B1"
  querypool_connection_mode = "%s"
}
`, rInt, location, rInt, connectionMode)
}

func testAccAzureRMAnalysisServicesServer_firewallSettings(rInt int, location string, enablePowerBIService bool, ipRules []map[string]string) string {
	ipRulesStr := make([]string, len(ipRules))
	for i, ipRule := range ipRules {
		ipRulesStr[i] = fmt.Sprintf(`ipv4_firewall_rule {
  name        = "%s"
  range_start = "%s"
  range_end   = "%s"
}
`, ipRule["name"], ipRule["range_start"], ipRule["range_end"])
	}

	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                		= "acctestass%d"
  location            		= "${azurerm_resource_group.test.location}"
  resource_group_name 		= "${azurerm_resource_group.test.name}"
  sku 				  		= "B1"
  enable_power_bi_service   = %t

  %s
}
`, rInt, location, rInt, enablePowerBIService, strings.Join(ipRulesStr, "\n"))
}

//func testAccAzureRMAnalysisServicesServer_adminUsers(rInt int, location string, adminUsers []string) string {
//	return fmt.Sprintf(`
//resource "azurerm_resource_group" "test" {
//  name     = "acctestRG-%d"
//  location = "%s"
//}
//
//resource "azurerm_analysis_services_server" "test" {
//  name                		= "acctestass%d"
//  location            		= "${azurerm_resource_group.test.location}"
//  resource_group_name 		= "${azurerm_resource_group.test.name}"
//  sku 				  		= "B1"
//  admin_users 				= ["%s"]
//}
//`, rInt, location, rInt, strings.Join(adminUsers, "\", \""))
//}

func testCheckAzureRMAnalysisServicesServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).analysisServicesServerClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_analysis_services_server" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.GetDetails(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMAnalysisServicesServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		analysisServicesServerName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Analysis Services Server: %s", analysisServicesServerName)
		}

		client := testAccProvider.Meta().(*ArmClient).analysisServicesServerClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.GetDetails(ctx, resourceGroup, analysisServicesServerName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Analysis Services Server %q (resource group: %q) does not exist", analysisServicesServerName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on analysisServicesServerClient: %+v", err)
		}

		return nil
	}
}
