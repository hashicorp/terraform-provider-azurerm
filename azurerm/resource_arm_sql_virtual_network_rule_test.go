package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSqlVirtualNetworkRule_basic(t *testing.T) {
	resourceName := "azurerm_sql_virtual_network_rule.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMSqlVirtualNetworkRule_basic(ri, testLocation())
	postConfig := testAccAzureRMSqlVirtualNetworkRule_withUpdates(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlVirtualNetworkRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ignore_missing_vnet_service_endpoint", "false"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlVirtualNetworkRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ignore_missing_vnet_service_endpoint", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMSqlVirtualNetworkRule_switchSubnets(t *testing.T) {
	resourceName := "azurerm_sql_virtual_network_rule.test"
	ri := acctest.RandInt()

	preConfig := testAccAzureRMSqlVirtualNetworkRule_subnetSwitchPre(ri, testLocation())
	postConfig := testAccAzureRMSqlVirtualNetworkRule_subnetSwitchPost(ri, testLocation())

	// Create regex strings that will ensure that one subnet name exists, but not the other
	preConfigRegex := regexp.MustCompile(fmt.Sprintf("(subnet1%d)$|(subnet[^2]%d)$", ri, ri))  //subnet 1 but not 2
	postConfigRegex := regexp.MustCompile(fmt.Sprintf("(subnet2%d)$|(subnet[^1]%d)$", ri, ri)) //subnet 2 but not 1

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlVirtualNetworkRuleExists(resourceName),
					resource.TestMatchResourceAttr(resourceName, "subnet_id", preConfigRegex),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlVirtualNetworkRuleExists(resourceName),
					resource.TestMatchResourceAttr(resourceName, "subnet_id", postConfigRegex),
				),
			},
		},
	})
}

func TestAccAzureRMSqlVirtualNetworkRule_disappears(t *testing.T) {
	resourceName := "azurerm_sql_virtual_network_rule.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSqlVirtualNetworkRule_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlVirtualNetworkRuleExists(resourceName),
					testCheckAzureRMSqlVirtualNetworkRuleDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

/*
	--Testing for Success--
	Test if we are able to create a vnet without the SQL endpoint, but SQL rule
	is still applied since the endpoint validation will be set to false.
*/
func TestAccAzureRMSqlVirtualNetworkRule_IgnoreEndpointValid(t *testing.T) {
	resourceName := "azurerm_sql_virtual_network_rule.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSqlVirtualNetworkRule_ignoreEndpointValid(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlVirtualNetworkRuleExists(resourceName),
				),
			},
		},
	})
}

/*
	--Testing for Failure--
	Test if we are able to create a vnet with out the SQL endpoint, but SQL rule
	is still applied since the endpoint validation will be set to false.
*/
func TestAccAzureRMSqlVirtualNetworkRule_IgnoreEndpointInvalid(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMSqlVirtualNetworkRule_ignoreEndpointInvalid(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Code=\"VirtualNetworkRuleBadRequest\""),
			},
		},
	})
}

/*
	--Testing for Success--
	Test if we are able to create multiple subnets and connect multiple subnets to the
	SQL server.
*/
func TestAccAzureRMSqlVirtualNetworkRule_multipleSubnets(t *testing.T) {
	resourceName1 := "azurerm_sql_virtual_network_rule.rule1"
	resourceName2 := "azurerm_sql_virtual_network_rule.rule2"
	resourceName3 := "azurerm_sql_virtual_network_rule.rule3"
	ri := acctest.RandInt()
	config := testAccAzureRMSqlVirtualNetworkRule_multipleSubnets(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlVirtualNetworkRuleExists(resourceName1),
					testCheckAzureRMSqlVirtualNetworkRuleExists(resourceName2),
					testCheckAzureRMSqlVirtualNetworkRuleExists(resourceName3),
				),
			},
		},
	})
}

/*
	Validation Function Tests - Invalid Name Validation
*/
func TestResourceAzureRMSqlVirtualNetworkRule_invalidNameValidation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		// Must only contain alphanumeric characters or hyphens (4 cases)
		{
			Value:    "test!Rule",
			ErrCount: 1,
		},
		{
			Value:    "test_Rule",
			ErrCount: 1,
		},
		{
			Value:    "test:Rule",
			ErrCount: 1,
		},
		{
			Value:    "test'Rule",
			ErrCount: 1,
		},
		// Cannot be more than 128 characters (1 case)
		{
			Value:    acctest.RandString(129),
			ErrCount: 1,
		},
		// Cannot be empty (1 case)
		{
			Value:    "",
			ErrCount: 1,
		},
		// Cannot end in a hyphen (1 case)
		{
			Value:    "testRule-",
			ErrCount: 1,
		},
		// Cannot start with a number or hyphen (2 cases)
		{
			Value:    "7testRule",
			ErrCount: 1,
		},
		{
			Value:    "-testRule",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateSqlVirtualNetworkRuleName(tc.Value, "azurerm_sql_virtual_network_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM SQL Virtual Network Rule Name to trigger a validation error.")
		}
	}
}

/*
	Function to assert if a rule exists or not.
*/
func testCheckAzureRMSqlVirtualNetworkRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlVirtualNetworkRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, ruleName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: SQL Firewall Rule %q (server %q / resource group %q) was not found", ruleName, serverName, resourceGroup)
			}

			return err
		}

		return nil
	}
}

/*
	Function to delete a rule.
*/
func testCheckAzureRMSqlVirtualNetworkRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sql_virtual_network_rule" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlVirtualNetworkRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, ruleName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Bad: SQL Firewall Rule %q (server %q / resource group %q) still exists: %+v", ruleName, serverName, resourceGroup, resp)
	}

	return nil
}

/*
	Function to assert if that a rule gets deleted.
*/
func testCheckAzureRMSqlVirtualNetworkRuleDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlVirtualNetworkRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		future, err := client.Delete(ctx, resourceGroup, serverName, ruleName)
		if err != nil {
			return fmt.Errorf("Bad: Error deleting SQL Virtual Network Rule: %+v", err)
		}

		err = future.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return nil // we are expecting the rule to not be there
		}

		//if the rule is still there, something
		return fmt.Errorf("Bad: Delete on sqlVirtualNetworkRulesClient: %+v", err)
	}
}

/*
	(This test configuration is intended to succeed.)
	Basic Provisioning Configuration
*/
func testAccAzureRMSqlVirtualNetworkRule_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG_%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/29"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
resource "azurerm_subnet" "test" {
  name = "acctestsubnet%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix = "10.7.29.0/29"
  service_endpoints = ["Microsoft.Sql"]
}
resource "azurerm_sql_server" "test" {
  name = "acctestsqlserver%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location = "${azurerm_resource_group.test.location}"
  version = "12.0"
  administrator_login = "missadmin"
  administrator_login_password = "${md5(%d)}!"
}
resource "azurerm_sql_virtual_network_rule" "test" {
  name = "acctestsqlvnetrule%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name = "${azurerm_sql_server.test.name}"
  subnet_id = "${azurerm_subnet.test.id}"
  ignore_missing_vnet_service_endpoint = false
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

/*
	(This test configuration is intended to succeed.)
	Basic Provisioning Update Configuration (all other properties would recreate the rule)
	ignore_missing_vnet_service_endpoint (false ==> true)
*/
func testAccAzureRMSqlVirtualNetworkRule_withUpdates(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG_%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/29"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
resource "azurerm_subnet" "test" {
  name = "acctestsubnet%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix = "10.7.29.0/29"
  service_endpoints = ["Microsoft.Sql"]
}
resource "azurerm_sql_server" "test" {
  name = "acctestsqlserver%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location = "${azurerm_resource_group.test.location}"
  version = "12.0"
  administrator_login = "missadmin"
  administrator_login_password = "${md5(%d)}!"
}
resource "azurerm_sql_virtual_network_rule" "test" {
  name = "acctestsqlvnetrule%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name = "${azurerm_sql_server.test.name}"
  subnet_id = "${azurerm_subnet.test.id}"
  ignore_missing_vnet_service_endpoint = true
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMSqlVirtualNetworkRule_subnetSwitchPre(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG_%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
resource "azurerm_subnet" "test1" {
  name = "subnet1%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix = "10.7.29.0/25"
  service_endpoints = ["Microsoft.Sql"]
}
resource "azurerm_subnet" "test2" {
  name = "subnet2%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix = "10.7.29.128/25"
  service_endpoints = ["Microsoft.Sql"]
}
resource "azurerm_sql_server" "test" {
  name = "acctestsqlserver%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location = "${azurerm_resource_group.test.location}"
  version = "12.0"
  administrator_login = "missadmin"
  administrator_login_password = "${md5(%d)}!"
}
resource "azurerm_sql_virtual_network_rule" "test" {
  name = "acctestsqlvnetrule%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name = "${azurerm_sql_server.test.name}"
  subnet_id = "${azurerm_subnet.test1.id}"
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMSqlVirtualNetworkRule_subnetSwitchPost(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG_%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
resource "azurerm_subnet" "test1" {
  name = "subnet1%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix = "10.7.29.0/25"
  service_endpoints = ["Microsoft.Sql"]
}
resource "azurerm_subnet" "test2" {
  name = "subnet2%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix = "10.7.29.128/25"
  service_endpoints = ["Microsoft.Sql"]
}
resource "azurerm_sql_server" "test" {
  name = "acctestsqlserver%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location = "${azurerm_resource_group.test.location}"
  version = "12.0"
  administrator_login = "missadmin"
  administrator_login_password = "${md5(%d)}!"
}
resource "azurerm_sql_virtual_network_rule" "test" {
  name = "acctestsqlvnetrule%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name = "${azurerm_sql_server.test.name}"
  subnet_id = "${azurerm_subnet.test2.id}"
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

/*
	(This test configuration is intended to succeed.)
	Succeeds because subnet's service_endpoints does not include 'Microsoft.Sql' and the SQL
    virtual network rule is set to *not* validate that the service_endpoint includes that value.
    The endpoint is purposefully set to Microsoft.Storage.
*/
func testAccAzureRMSqlVirtualNetworkRule_ignoreEndpointValid(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG_%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/29"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
resource "azurerm_subnet" "test" {
  name = "acctestsubnet%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix = "10.7.29.0/29"
  service_endpoints = ["Microsoft.Storage"]
}
resource "azurerm_sql_server" "test" {
  name = "acctestsqlserver%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location = "${azurerm_resource_group.test.location}"
  version = "12.0"
  administrator_login = "missadmin"
  administrator_login_password = "${md5(%d)}!"
}
resource "azurerm_sql_virtual_network_rule" "test" {
  name = "acctestsqlvnetrule%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name = "${azurerm_sql_server.test.name}"
  subnet_id = "${azurerm_subnet.test.id}"
  ignore_missing_vnet_service_endpoint = true
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

/*
	(This test configuration is intended to fail.)
	Fails because subnet's service_endpoints does not include 'Microsoft.Sql' and the SQL
    virtual network rule is set to validate that the service_endpoint includes that value.
    The endpoint is purposefully set to Microsoft.Storage.
*/
func testAccAzureRMSqlVirtualNetworkRule_ignoreEndpointInvalid(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG_%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/29"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
resource "azurerm_subnet" "test" {
  name = "acctestsubnet%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix = "10.7.29.0/29"
  service_endpoints = ["Microsoft.Storage"]
}
resource "azurerm_sql_server" "test" {
  name = "acctestsqlserver%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location = "${azurerm_resource_group.test.location}"
  version = "12.0"
  administrator_login = "missadmin"
  administrator_login_password = "${md5(%d)}!"
}
resource "azurerm_sql_virtual_network_rule" "test" {
  name = "acctestsqlvnetrule%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name = "${azurerm_sql_server.test.name}"
  subnet_id = "${azurerm_subnet.test.id}"
  ignore_missing_vnet_service_endpoint = false
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

/*
	(This test configuration is intended to succeed.)
	This configuration sets up 3 subnets in 2 different virtual networks, and adds
	SQL virtual network rules for all 3 subnets to the SQL server.
*/
func testAccAzureRMSqlVirtualNetworkRule_multipleSubnets(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG_%d"
  location = "%s"
}
resource "azurerm_virtual_network" "vnet1" {
  name                = "acctestvnet1%d"
  address_space       = ["10.7.29.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
resource "azurerm_virtual_network" "vnet2" {
  name                = "acctestvnet2%d"
  address_space       = ["10.1.29.0/29"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
resource "azurerm_subnet" "vnet1_subnet1" {
  name = "acctestsubnet1%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet1.name}"
  address_prefix = "10.7.29.0/29"
  service_endpoints = ["Microsoft.Sql"]
}
resource "azurerm_subnet" "vnet1_subnet2" {
  name = "acctestsubnet2%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet1.name}"
  address_prefix = "10.7.29.128/29"
  service_endpoints = ["Microsoft.Sql"]
}
resource "azurerm_subnet" "vnet2_subnet1" {
  name = "acctestsubnet3%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet2.name}"
  address_prefix = "10.1.29.0/29"
  service_endpoints = ["Microsoft.Sql"]
}
resource "azurerm_sql_server" "test" {
  name = "acctestsqlserver1%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location = "${azurerm_resource_group.test.location}"
  version = "12.0"
  administrator_login = "missadmin"
  administrator_login_password = "${md5(%d)}!"
}
resource "azurerm_sql_virtual_network_rule" "rule1" {
  name = "acctestsqlvnetrule1%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name = "${azurerm_sql_server.test.name}"
  subnet_id = "${azurerm_subnet.vnet1_subnet1.id}"
  ignore_missing_vnet_service_endpoint = false
}
resource "azurerm_sql_virtual_network_rule" "rule2" {
  name = "acctestsqlvnetrule2%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name = "${azurerm_sql_server.test.name}"
  subnet_id = "${azurerm_subnet.vnet1_subnet2.id}"
  ignore_missing_vnet_service_endpoint = false
}
resource "azurerm_sql_virtual_network_rule" "rule3" {
  name = "acctestsqlvnetrule3%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name = "${azurerm_sql_server.test.name}"
  subnet_id = "${azurerm_subnet.vnet2_subnet1.id}"
  ignore_missing_vnet_service_endpoint = false
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}
