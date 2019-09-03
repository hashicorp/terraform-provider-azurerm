package azurerm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func testCheckAzureRMLogicAppActionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		logicAppId := rs.Primary.Attributes["logic_app_id"]
		id, err := azure.ParseAzureResourceID(logicAppId)
		if err != nil {
			return err
		}

		actionName := rs.Primary.Attributes["name"]
		workflowName := id.Path["workflows"]
		resourceGroup := id.ResourceGroup

		client := testAccProvider.Meta().(*ArmClient).logic.WorkflowsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, workflowName)
		if err != nil {
			return fmt.Errorf("Bad: Get on logicWorkflowsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Logic App Workflow %q (resource group %q) does not exist", workflowName, resourceGroup)
		}

		definition := resp.WorkflowProperties.Definition.(map[string]interface{})
		actions := definition["actions"].(map[string]interface{})

		exists := false
		for k := range actions {
			if strings.EqualFold(k, actionName) {
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("Action %q was not found on Logic App %q (Resource Group %q)", actionName, workflowName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMLogicAppTriggerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		logicAppId := rs.Primary.Attributes["logic_app_id"]
		id, err := azure.ParseAzureResourceID(logicAppId)
		if err != nil {
			return err
		}

		triggerName := rs.Primary.Attributes["name"]
		workflowName := id.Path["workflows"]
		resourceGroup := id.ResourceGroup

		client := testAccProvider.Meta().(*ArmClient).logic.WorkflowsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, workflowName)
		if err != nil {
			return fmt.Errorf("Bad: Get on logicWorkflowsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Logic App Workflow %q (resource group %q) does not exist", workflowName, resourceGroup)
		}

		definition := resp.WorkflowProperties.Definition.(map[string]interface{})
		triggers := definition["triggers"].(map[string]interface{})

		exists := false
		for k := range triggers {
			if strings.EqualFold(k, triggerName) {
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("Trigger %q was not found on Logic App %q (Resource Group %q)", triggerName, workflowName, resourceGroup)
		}

		return nil
	}
}
