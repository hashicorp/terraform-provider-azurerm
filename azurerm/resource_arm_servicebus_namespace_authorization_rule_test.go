package azurerm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"strconv"
)

func TestAccAzureRMServiceBusNamespaceAuthorizationRule_rights(t *testing.T) {
	cases := []struct {
		Name   string
		Rights []string
	}{
		{
			Name:   "listen",
			Rights: []string{string(servicebus.Listen)},
		},
		{
			Name:   "send",
			Rights: []string{string(servicebus.Send)},
		},
		{
			Name:   "send-listen",
			Rights: []string{string(servicebus.Listen), string(servicebus.Send)},
		},
		{
			Name:   "manage",
			Rights: []string{string(servicebus.Listen), string(servicebus.Send), string(servicebus.Manage)},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			testAccAzureRMServiceBusNamespaceAuthorizationRule(t, tc.Rights)
		})
	}
}

func testAccAzureRMServiceBusNamespaceAuthorizationRule(t *testing.T, rights []string) {
	resourceName := "azurerm_servicebus_namespace_authorization_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespaceAuthorizationRule_base(acctest.RandInt(), testLocation(), rights),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_name"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttr(resourceName, "rights.#", strconv.Itoa(len(rights))),
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

func TestAccAzureRMServiceBusNamespaceAuthorizationRule_rightsUpdate(t *testing.T) {
	resourceName := "azurerm_servicebus_namespace_authorization_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespaceAuthorizationRule_base(acctest.RandInt(), testLocation(), []string{"Listen"}),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rights.#", "1"),
				),
			},
			{
				Config: testAccAzureRMServiceBusNamespaceAuthorizationRule_base(acctest.RandInt(), testLocation(), []string{"Listen", "Send", "Manage"}),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_name"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttr(resourceName, "rights.#", "3"),
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

func testCheckAzureRMServiceBusNamespaceAuthorizationRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).serviceBusNamespacesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_servicebus_topic_authorization_rule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.GetAuthorizationRule(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}

	return nil
}

func testCheckAzureRMServiceBusNamespaceAuthorizationRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for ServiceBus Namespace: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).serviceBusNamespacesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.GetAuthorizationRule(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: ServiceBus Namespace Authorization Rule %q (namespace %s / resource group: %s) does not exist", name, namespaceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on ServiceBus Namespace: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMServiceBusNamespaceAuthorizationRule_base(rInt int, location string, rights []string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_namespace_authorization_rule" "test" {
  name                = "acctest-%[1]d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  rights              = ["%[3]s"]
}
`, rInt, location, strings.Join(rights, `","`))
}
