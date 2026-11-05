// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_workflow

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogicAppWorkflowUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workflows.ParseWorkflowID(d.Id())
	if err != nil {
		return err
	}

	// lock to prevent against Actions, Parameters or Triggers conflicting
	locks.ByName(id.WorkflowName, logicAppResourceName)
	defer locks.UnlockByName(id.WorkflowName, logicAppResourceName)

	read, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %s: %+v", id, err)
	}

	if read.Model == nil || read.Model.Properties == nil {
		return fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `WorkflowProperties` is nil")
	}

	location := location.Normalize(d.Get("location").(string))
	workflowParameters, err := expandLogicAppWorkflowWorkflowParameters(d.Get("workflow_parameters").(map[string]any))
	if err != nil {
		return fmt.Errorf("expanding `workflow_parameters`: %+v", err)
	}
	parameters, err := expandLogicAppWorkflowParameters(d.Get("parameters").(map[string]any), workflowParameters)
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]any)

	var definition any
	if read.Model.Properties.Definition != nil {
		definitionRaw := *read.Model.Properties.Definition
		definitionMap := definitionRaw.(map[string]any)
		definitionMap["parameters"] = workflowParameters
		definition = definitionMap
	}

	isEnabled := workflows.WorkflowStateEnabled
	if v := d.Get("enabled").(bool); !v {
		isEnabled = workflows.WorkflowStateDisabled
	}

	identity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	properties := workflows.Workflow{
		Identity: identity,
		Location: pointer.To(location),
		Properties: &workflows.WorkflowProperties{
			Definition: &definition,
			Parameters: parameters,
			State:      &isEnabled,
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("access_control"); ok {
		properties.Properties.AccessControl = expandLogicAppWorkflowAccessControl(v.([]any))
	}

	if v, ok := d.GetOk("logic_app_integration_account_id"); ok {
		properties.Properties.IntegrationAccount = &workflows.ResourceReference{
			Id: pointer.To(v.(string)),
		}
	}

	if iseID, ok := d.GetOk("integration_service_environment_id"); ok {
		properties.Properties.IntegrationServiceEnvironment = &workflows.ResourceReference{
			Id: pointer.To(iseID.(string)),
		}
	}

	if _, err = client.CreateOrUpdate(ctx, *id, properties); err != nil {
		return fmt.Errorf("updating Logic App Workflow %s: %+v", id, err)
	}

	return resourceLogicAppWorkflowRead(d, meta)
}
