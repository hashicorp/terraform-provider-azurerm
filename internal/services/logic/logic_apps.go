// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowtriggers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

// NOTE: this file is not a recommended way of developing Terraform resources; this exists to work around the fact that this API is dynamic (by its nature)
func flattenLogicAppActionRunAfter(input map[string]interface{}) []interface{} {
	if len(input) == 0 {
		return nil
	}
	output := []interface{}{}
	for k, v := range input {
		output = append(output, map[string]interface{}{
			"action_name":   k,
			"action_result": v.([]interface{})[0],
		})
	}

	return output
}

func expandLogicAppActionRunAfter(input []interface{}) map[string]interface{} {
	if len(input) == 0 {
		return nil
	}
	output := map[string]interface{}{}
	for _, v := range input {
		b := v.(map[string]interface{})
		output[b["action_name"].(string)] = []string{b["action_result"].(string)}
	}

	return output
}

func resourceLogicAppActionUpdate(d *pluginsdk.ResourceData, meta interface{}, workflowId workflows.WorkflowId, actionId parse.ActionId, vals map[string]interface{}, resourceName string) error {
	return resourceLogicAppComponentUpdate(d, meta, "Action", "actions", workflowId, actionId.ID(), actionId.Name, vals, resourceName)
}

func resourceLogicAppTriggerUpdate(d *pluginsdk.ResourceData, meta interface{}, workflowId workflows.WorkflowId, triggerId workflowtriggers.TriggerId, vals map[string]interface{}, resourceName string) error {
	return resourceLogicAppComponentUpdate(d, meta, "Trigger", "triggers", workflowId, triggerId.ID(), triggerId.TriggerName, vals, resourceName)
}

func resourceLogicAppComponentUpdate(d *pluginsdk.ResourceData, meta interface{}, kind string, propertyName string, workflowId workflows.WorkflowId, resourceId string, name string, vals map[string]interface{}, resourceName string) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[DEBUG] Preparing arguments for Logic App Workspace %s %s %q", workflowId, kind, name)

	// lock to prevent against Actions or Triggers conflicting
	locks.ByName(workflowId.WorkflowName, logicAppResourceName)
	defer locks.UnlockByName(workflowId.WorkflowName, logicAppResourceName)

	read, err := client.Get(ctx, workflowId)
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			return fmt.Errorf("[ERROR] Logic App Workflow %s was not found", workflowId)
		}

		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %s: %+v", workflowId, err)
	}

	if read.Model == nil || read.Model.Properties == nil {
		return fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `WorkflowProperties` is nil")
	}

	if read.Model.Properties.Definition == nil {
		return fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `WorkflowProperties.Definition` is nil")
	}

	rawDefinition := *read.Model.Properties.Definition
	definitionMap := rawDefinition.(map[string]interface{})
	vs := definitionMap[propertyName].(map[string]interface{})

	if d.IsNewResource() {
		if _, hasExisting := vs[name]; hasExisting {
			return tf.ImportAsExistsError(resourceName, resourceId)
		}
	}

	vs[name] = vals
	definitionMap[propertyName] = vs
	rawDefinition = definitionMap

	if read.Model.Identity != nil && read.Model.Identity.IdentityIds != nil {
		for k := range read.Model.Identity.IdentityIds {
			read.Model.Identity.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				// this has to be an empty object due to the API design
			}
		}
	}

	properties := workflows.Workflow{
		Location: read.Model.Location,
		Properties: &workflows.WorkflowProperties{
			Definition:                    &rawDefinition,
			Parameters:                    read.Model.Properties.Parameters,
			AccessControl:                 read.Model.Properties.AccessControl,
			IntegrationAccount:            read.Model.Properties.IntegrationAccount,
			IntegrationServiceEnvironment: read.Model.Properties.IntegrationServiceEnvironment,
		},
		Identity: read.Model.Identity,
		Tags:     read.Model.Tags,
	}
	if _, err = client.CreateOrUpdate(ctx, workflowId, properties); err != nil {
		return fmt.Errorf("updating Logic App Workflow %s for %s %q: %+v", workflowId, kind, name, err)
	}

	if d.IsNewResource() {
		d.SetId(resourceId)
	}

	return nil
}

func resourceLogicAppActionRemove(d *pluginsdk.ResourceData, meta interface{}, id workflows.WorkflowId, name string) error {
	return resourceLogicAppComponentRemove(d, meta, "Action", "actions", id, name)
}

func resourceLogicAppTriggerRemove(d *pluginsdk.ResourceData, meta interface{}, id workflows.WorkflowId, name string) error {
	return resourceLogicAppComponentRemove(d, meta, "Trigger", "triggers", id, name)
}

func resourceLogicAppComponentRemove(d *pluginsdk.ResourceData, meta interface{}, kind, propertyName string, id workflows.WorkflowId, name string) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[DEBUG] Preparing arguments for Logic App Workspace %q (Resource Group %q) %s %q Deletion", id.WorkflowName, id.ResourceGroupName, kind, name)

	// lock to prevent against Actions, Parameters or Actions conflicting
	locks.ByName(id.WorkflowName, logicAppResourceName)
	defer locks.UnlockByName(id.WorkflowName, logicAppResourceName)

	read, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Error making Read request on %s: %+v", id.ID(), err)
	}

	if read.Model == nil {
		return fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `Model` is nil")
	}

	if read.Model.Properties == nil {
		return fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `Properties` is nil")
	}

	if read.Model.Properties.Definition == nil {
		return fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `WorkflowProperties.Definition` is nil")
	}

	var definition interface{}
	definitionRaw := *read.Model.Properties.Definition
	definitionMap := definitionRaw.(map[string]interface{})
	vs := definitionMap[propertyName].(map[string]interface{})
	delete(vs, name)
	definitionMap[propertyName] = vs
	definition = definitionMap

	properties := workflows.Workflow{
		Location: read.Model.Location,
		Properties: &workflows.WorkflowProperties{
			Definition:                    &definition,
			Parameters:                    read.Model.Properties.Parameters,
			AccessControl:                 read.Model.Properties.AccessControl,
			IntegrationAccount:            read.Model.Properties.IntegrationAccount,
			IntegrationServiceEnvironment: read.Model.Properties.IntegrationServiceEnvironment,
		},
		Tags: read.Model.Tags,
	}

	if _, err = client.CreateOrUpdate(ctx, id, properties); err != nil {
		return fmt.Errorf("removing %s %q from %s: %+v", kind, name, id.ID(), err)
	}

	return nil
}

func retrieveLogicAppAction(d *pluginsdk.ResourceData, meta interface{}, id workflows.WorkflowId, name string) (*map[string]interface{}, *workflows.Workflow, error) {
	return retrieveLogicAppComponent(d, meta, "Action", "actions", id, name)
}

func retrieveLogicAppTrigger(d *pluginsdk.ResourceData, meta interface{}, id workflowtriggers.TriggerId) (*map[string]interface{}, *workflows.Workflow, *string, error) {
	workflowId := workflows.NewWorkflowID(id.SubscriptionId, id.ResourceGroupName, id.WorkflowName)

	t, app, err := retrieveLogicAppComponent(d, meta, "Trigger", "triggers", workflowId, id.TriggerName)

	if err != nil || t == nil {
		return nil, nil, nil, err
	}

	trigger := *t
	tType := trigger["type"]
	if tType == nil {
		return nil, nil, nil, fmt.Errorf("[ERROR] `type` was nil for %s", id.ID())
	}

	log.Printf("[DEBUG] trigger type is %s", tType.(string))

	if IsCallbackType(tType.(string)) {
		url, err := retreiveLogicAppTriggerCallbackUrl(d, meta, id)
		if err != nil {
			return nil, nil, nil, err
		}

		if url == nil {
			return nil, nil, nil, fmt.Errorf("[ERROR] `callback_url` was nil for %s", id.ID())
		}

		return t, app, url, err
	}

	return t, app, nil, err
}

func IsCallbackType(tType string) bool {
	cTypes := []string{"ApiConnectionWebhook", "HTTPWebhook", "Request"}

	valid := validation.StringInSlice(cTypes, false)
	_, errors := valid(tType, "callback_url")

	return len(errors) == 0
}

func retreiveLogicAppTriggerCallbackUrl(d *pluginsdk.ResourceData, meta interface{}, id workflowtriggers.TriggerId) (*string, error) {
	client := meta.(*clients.Client).Logic
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[DEBUG] Preparing arguments for Logic App Workspace %q (Resource Group %q) %s %q", id.WorkflowName, id.ResourceGroupName, "trigger", id.TriggerName)

	// lock to prevent against Actions, Parameters or Actions conflicting
	locks.ByName(id.WorkflowName, logicAppResourceName)
	defer locks.UnlockByName(id.WorkflowName, logicAppResourceName)

	result, err := client.TriggersClient.ListCallbackURL(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error getting trigger callback URL (%w)", err)
	}

	if result.Model == nil {
		return nil, fmt.Errorf("[ERROR] model was nil for %s", id.ID())
	}

	return result.Model.Value, nil
}

func retrieveLogicAppComponent(d *pluginsdk.ResourceData, meta interface{}, kind, propertyName string, id workflows.WorkflowId, name string) (*map[string]interface{}, *workflows.Workflow, error) {
	client := meta.(*clients.Client).Logic.WorkflowClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[DEBUG] Preparing arguments for %s: %s %q", id.ID(), kind, name)

	// lock to prevent against Actions, Parameters or Actions conflicting
	locks.ByName(id.WorkflowName, logicAppResourceName)
	defer locks.UnlockByName(id.WorkflowName, logicAppResourceName)

	read, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			return nil, nil, nil
		}

		return nil, nil, fmt.Errorf("[ERROR] Error making Read request %s: %+v", id.ID(), err)
	}

	if read.Model == nil {
		return nil, nil, fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `Model` is nil")
	}

	if read.Model.Properties == nil {
		return nil, nil, fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `Properties` is nil")
	}

	if read.Model.Properties.Definition == nil {
		return nil, nil, fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `Properties.Definition` is nil")
	}

	definitionRaw := *read.Model.Properties.Definition
	definitionMap := definitionRaw.(map[string]interface{})
	vs := definitionMap[propertyName].(map[string]interface{})
	v := vs[name]
	if v == nil {
		return nil, nil, nil
	}

	result := v.(map[string]interface{})
	return &result, read.Model, nil
}
