package automation_test

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testCheckAzureRMAutomationConnectionDestroy(s *terraform.State, varType string) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.ConnectionClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	resourceName := "azurerm_automation_connection"
	if varType != "" {
		resourceName = fmt.Sprintf("azurerm_automation_connection_%s", strings.ToLower(varType))
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourceName {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		accName := rs.Primary.Attributes["automation_account_name"]

		if resp, err := conn.Get(ctx, resourceGroup, accName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Automation.ConnectionClient: %+v", err)
			}
		}
	}

	return nil
}

func testCheckAzureRMAutomationConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.ConnectionClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		accName := rs.Primary.Attributes["automation_account_name"]

		if resp, err := conn.Get(ctx, resourceGroup, accName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Connection %q (Resource Group %q / automation account %q) does not exist", name, resourceGroup, accName)
			}
			return fmt.Errorf("bad: Get on automationConnectionClient: %+v", err)
		}

		return nil
	}
}
