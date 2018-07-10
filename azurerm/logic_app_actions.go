package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2016-06-01/logic"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLogicAppActionUpdate(d *schema.ResourceData, meta interface{}, logicAppId string, name string, vals map[string]interface{}) error {
	client := meta.(*ArmClient).logicWorkflowsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(logicAppId)
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	logicAppName := id.Path["workflows"]

	log.Printf("[DEBUG] Preparing arguments for Logic App Workspace %q (Resource Group %q) Action %q", logicAppName, resourceGroup, name)

	// lock to prevent against Actions, Parameters or Actions conflicting
	azureRMLockByName(logicAppName, logicAppResourceName)
	defer azureRMUnlockByName(logicAppName, logicAppResourceName)

	read, err := client.Get(ctx, resourceGroup, logicAppName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %q (Resource Group %q): %+v", logicAppName, resourceGroup, err)
	}

	if read.WorkflowProperties == nil {
		return fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `WorkflowProperties` is nil")
	}

	if read.WorkflowProperties.Definition == nil {
		return fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `WorkflowProperties.Definition` is nil")
	}

	definition := read.WorkflowProperties.Definition.(map[string]interface{})
	actions := definition["actions"].(map[string]interface{})
	actions[name] = vals
	definition["actions"] = actions

	properties := logic.Workflow{
		Location: read.Location,
		WorkflowProperties: &logic.WorkflowProperties{
			Definition: definition,
			Parameters: read.WorkflowProperties.Parameters,
		},
		Tags: read.Tags,
	}

	_, err = client.CreateOrUpdate(ctx, resourceGroup, logicAppName, properties)
	if err != nil {
		return fmt.Errorf("Error updating Logic App Workspace %q (Resource Group %q) for Action %q: %+v", logicAppName, resourceGroup, name, err)
	}

	if d.IsNewResource() {
		d.SetId(fmt.Sprintf("%s/actions/%s", *read.ID, name))
	}

	return nil
}

func resourceArmLogicAppActionRemove(d *schema.ResourceData, meta interface{}, resourceGroup, logicAppName, name string) error {
	client := meta.(*ArmClient).logicWorkflowsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[DEBUG] Preparing arguments for Logic App Workspace %q (Resource Group %q) Action %q Deletion", logicAppName, resourceGroup, name)

	// lock to prevent against Actions, Parameters or Actions conflicting
	azureRMLockByName(logicAppName, logicAppResourceName)
	defer azureRMUnlockByName(logicAppName, logicAppResourceName)

	read, err := client.Get(ctx, resourceGroup, logicAppName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %q (Resource Group %q): %+v", logicAppName, resourceGroup, err)
	}

	if read.WorkflowProperties == nil {
		return fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `WorkflowProperties` is nil")
	}

	if read.WorkflowProperties.Definition == nil {
		return fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `WorkflowProperties.Definition` is nil")
	}

	definition := read.WorkflowProperties.Definition.(map[string]interface{})
	actions := definition["actions"].(map[string]interface{})
	delete(actions, name)
	definition["actions"] = actions

	properties := logic.Workflow{
		Location: read.Location,
		WorkflowProperties: &logic.WorkflowProperties{
			Definition: definition,
			Parameters: read.WorkflowProperties.Parameters,
		},
		Tags: read.Tags,
	}

	_, err = client.CreateOrUpdate(ctx, resourceGroup, logicAppName, properties)
	if err != nil {
		return fmt.Errorf("Error removing Action %q from Logic App Workspace %q (Resource Group %q): %+v", name, logicAppName, resourceGroup, err)
	}

	return nil
}

func retrieveLogicAppAction(meta interface{}, resourceGroup, logicAppName, name string) (*map[string]interface{}, *logic.Workflow, error) {
	client := meta.(*ArmClient).logicWorkflowsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[DEBUG] Preparing arguments for Logic App Workspace %q (Resource Group %q) Action %q", logicAppName, resourceGroup, name)

	// lock to prevent against Actions, Parameters or Actions conflicting
	azureRMLockByName(logicAppName, logicAppResourceName)
	defer azureRMUnlockByName(logicAppName, logicAppResourceName)

	read, err := client.Get(ctx, resourceGroup, logicAppName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return nil, nil, nil
		}

		return nil, nil, fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %q (Resource Group %q): %+v", logicAppName, resourceGroup, err)
	}

	if read.WorkflowProperties == nil {
		return nil, nil, fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `WorkflowProperties` is nil")
	}

	if read.WorkflowProperties.Definition == nil {
		return nil, nil, fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `WorkflowProperties.Definition` is nil")
	}

	definition := read.WorkflowProperties.Definition.(map[string]interface{})
	actions := definition["actions"].(map[string]interface{})
	action := actions[name]
	if action == nil {
		return nil, nil, nil
	}

	v := action.(map[string]interface{})
	return &v, &read, nil
}
