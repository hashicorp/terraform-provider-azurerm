package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2016-06-01/logic"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// NOTE: this file is not a recommended way of developing Terraform resources; this exists to work around the fact that this API is dynamic (by it's nature)

func resourceLogicAppActionUpdate(d *schema.ResourceData, meta interface{}, logicAppId string, name string, vals map[string]interface{}, resourceName string) error {
	return resourceLogicAppComponentUpdate(d, meta, "Action", "actions", logicAppId, name, vals, resourceName)
}

func resourceLogicAppTriggerUpdate(d *schema.ResourceData, meta interface{}, logicAppId string, name string, vals map[string]interface{}, resourceName string) error {
	return resourceLogicAppComponentUpdate(d, meta, "Trigger", "triggers", logicAppId, name, vals, resourceName)
}

func resourceLogicAppComponentUpdate(d *schema.ResourceData, meta interface{}, kind string, propertyName string, logicAppId string, name string, vals map[string]interface{}, resourceName string) error {
	client := meta.(*ArmClient).logic.WorkflowsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(logicAppId)
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	logicAppName := id.Path["workflows"]

	log.Printf("[DEBUG] Preparing arguments for Logic App Workspace %q (Resource Group %q) %s %q", logicAppName, resourceGroup, kind, name)

	// lock to prevent against Actions or Triggers conflicting
	locks.ByName(logicAppName, logicAppResourceName)
	defer locks.UnlockByName(logicAppName, logicAppResourceName)

	read, err := client.Get(ctx, resourceGroup, logicAppName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("[ERROR] Logic App Workflow %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %q (Resource Group %q): %+v", logicAppName, resourceGroup, err)
	}

	if read.WorkflowProperties == nil {
		return fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `WorkflowProperties` is nil")
	}

	if read.WorkflowProperties.Definition == nil {
		return fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `WorkflowProperties.Definition` is nil")
	}

	resourceId := fmt.Sprintf("%s/%s/%s", *read.ID, propertyName, name)

	definition := read.WorkflowProperties.Definition.(map[string]interface{})
	vs := definition[propertyName].(map[string]interface{})

	if requireResourcesToBeImported && d.IsNewResource() {
		if _, hasExisting := vs[name]; hasExisting {
			return tf.ImportAsExistsError(resourceName, resourceId)
		}
	}

	vs[name] = vals
	definition[propertyName] = vs

	properties := logic.Workflow{
		Location: read.Location,
		WorkflowProperties: &logic.WorkflowProperties{
			Definition: definition,
			Parameters: read.WorkflowProperties.Parameters,
		},
		Tags: read.Tags,
	}

	if _, err = client.CreateOrUpdate(ctx, resourceGroup, logicAppName, properties); err != nil {
		return fmt.Errorf("Error updating Logic App Workspace %q (Resource Group %q) for %s %q: %+v", logicAppName, resourceGroup, kind, name, err)
	}

	if d.IsNewResource() {
		d.SetId(resourceId)
	}

	return nil
}

func resourceLogicAppActionRemove(d *schema.ResourceData, meta interface{}, resourceGroup, logicAppName, name string) error {
	return resourceLogicAppComponentRemove(d, meta, "Action", "actions", resourceGroup, logicAppName, name)
}

func resourceLogicAppTriggerRemove(d *schema.ResourceData, meta interface{}, resourceGroup, logicAppName, name string) error {
	return resourceLogicAppComponentRemove(d, meta, "Trigger", "triggers", resourceGroup, logicAppName, name)
}

func resourceLogicAppComponentRemove(d *schema.ResourceData, meta interface{}, kind, propertyName, resourceGroup, logicAppName, name string) error {
	client := meta.(*ArmClient).logic.WorkflowsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[DEBUG] Preparing arguments for Logic App Workspace %q (Resource Group %q) %s %q Deletion", logicAppName, resourceGroup, kind, name)

	// lock to prevent against Actions, Parameters or Actions conflicting
	locks.ByName(logicAppName, logicAppResourceName)
	defer locks.UnlockByName(logicAppName, logicAppResourceName)

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
	vs := definition[propertyName].(map[string]interface{})
	delete(vs, name)
	definition[propertyName] = vs

	properties := logic.Workflow{
		Location: read.Location,
		WorkflowProperties: &logic.WorkflowProperties{
			Definition: definition,
			Parameters: read.WorkflowProperties.Parameters,
		},
		Tags: read.Tags,
	}

	if _, err = client.CreateOrUpdate(ctx, resourceGroup, logicAppName, properties); err != nil {
		return fmt.Errorf("Error removing %s %q from Logic App Workspace %q (Resource Group %q): %+v", kind, name, logicAppName, resourceGroup, err)
	}

	return nil
}

func retrieveLogicAppAction(meta interface{}, resourceGroup, logicAppName, name string) (*map[string]interface{}, *logic.Workflow, error) {
	return retrieveLogicAppComponent(meta, resourceGroup, "Action", "actions", logicAppName, name)
}

func retrieveLogicAppTrigger(meta interface{}, resourceGroup, logicAppName, name string) (*map[string]interface{}, *logic.Workflow, error) {
	return retrieveLogicAppComponent(meta, resourceGroup, "Trigger", "triggers", logicAppName, name)
}

func retrieveLogicAppComponent(meta interface{}, resourceGroup, kind, propertyName, logicAppName, name string) (*map[string]interface{}, *logic.Workflow, error) {
	client := meta.(*ArmClient).logic.WorkflowsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[DEBUG] Preparing arguments for Logic App Workspace %q (Resource Group %q) %s %q", logicAppName, resourceGroup, kind, name)

	// lock to prevent against Actions, Parameters or Actions conflicting
	locks.ByName(logicAppName, logicAppResourceName)
	defer locks.UnlockByName(logicAppName, logicAppResourceName)

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
	vs := definition[propertyName].(map[string]interface{})
	v := vs[name]
	if v == nil {
		return nil, nil, nil
	}

	result := v.(map[string]interface{})
	return &result, &read, nil
}
